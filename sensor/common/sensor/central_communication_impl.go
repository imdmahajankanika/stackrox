package sensor

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/pkg/centralsensor"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/sensor/common"
	"github.com/stackrox/rox/sensor/common/clusterid"
	"github.com/stackrox/rox/sensor/common/config"
	"github.com/stackrox/rox/sensor/common/detector"
	"github.com/stackrox/rox/sensor/common/sensor/helmconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// sensor implements the Sensor interface by sending inputs to central,
// and providing the output from central asynchronously.
type centralCommunicationImpl struct {
	receiver   CentralReceiver
	sender     CentralSender
	components []common.SensorComponent

	stopC    concurrency.ErrorSignal
	stoppedC concurrency.ErrorSignal
}

func (s *centralCommunicationImpl) Start(conn grpc.ClientConnInterface, centralReachable *concurrency.Flag, configHandler config.Handler, detector detector.Detector) {
	go s.sendEvents(central.NewSensorServiceClient(conn), centralReachable, configHandler, detector, s.receiver.Stop, s.sender.Stop)
}

func (s *centralCommunicationImpl) Stop(err error) {
	s.stopC.SignalWithError(err)
}

func (s *centralCommunicationImpl) Stopped() concurrency.ReadOnlyErrorSignal {
	return &s.stoppedC
}

func isUnimplemented(err error) bool {
	spb, ok := status.FromError(err)
	if spb == nil || !ok {
		return false
	}
	return spb.Code() == codes.Unimplemented
}

func communicateWithAutoSensedEncoding(ctx context.Context, client central.SensorServiceClient) (central.SensorService_CommunicateClient, error) {
	opts := []grpc.CallOption{grpc.UseCompressor(gzip.Name)}

	for {
		stream, err := client.Communicate(ctx, opts...)
		if err != nil {
			if isUnimplemented(err) && len(opts) > 0 {
				opts = nil
				continue
			}
			return nil, errors.Wrap(err, "opening stream")
		}

		_, err = stream.Header()
		if err != nil {
			if isUnimplemented(err) && len(opts) > 0 {
				opts = nil
				continue
			}
			return nil, errors.Wrap(err, "receiving initial metadata")
		}

		return stream, nil
	}
}

func (s *centralCommunicationImpl) sendEvents(client central.SensorServiceClient, centralReachable *concurrency.Flag, configHandler config.Handler, detector detector.Detector, onStops ...func(error)) {
	defer func() {
		s.stoppedC.SignalWithError(s.stopC.Err())
		runAll(s.stopC.Err(), onStops...)
	}()

	// Start the stream client.
	///////////////////////////
	ctx, err := centralsensor.AppendSensorVersionInfoToContext(context.Background())
	if err != nil {
		s.stopC.SignalWithError(err)
		return
	}

	capsSet := centralsensor.NewSensorCapabilitySet()
	for _, component := range s.components {
		capsSet.AddAll(component.Capabilities()...)
	}
	ctx = centralsensor.AppendCapsInfoToContext(ctx, capsSet)
	if configHandler.GetHelmManagedConfig() != nil {
		ctx = metadata.AppendToOutgoingContext(ctx, centralsensor.HelmManagedClusterMetadataKey, "true")
	}

	stream, err := communicateWithAutoSensedEncoding(ctx, client)
	if err != nil {
		s.stopC.SignalWithError(err)
		return
	}

	if err := s.initialSync(stream, capsSet, configHandler, detector); err != nil {
		s.stopC.SignalWithError(err)
		return
	}

	defer func() {
		if err := stream.CloseSend(); err != nil {
			log.Errorf("Failed to close stream cleanly: %v", err)
		}
	}()
	log.Info("Established connection to Central.")

	centralReachable.Set(true)
	defer centralReachable.Set(false)

	// Start receiving and sending with central.
	////////////////////////////////////////////
	s.receiver.Start(stream, s.Stop, s.sender.Stop)
	s.sender.Start(stream, s.Stop, s.receiver.Stop)
	log.Info("Communication with central started.")

	// Wait for stop.
	/////////////////
	_ = s.stopC.Wait()
	log.Info("Communication with central ended.")
}

func (s *centralCommunicationImpl) initialSync(stream central.SensorService_CommunicateClient, capabilities centralsensor.SensorCapabilitySet, configHandler config.Handler, detector detector.Detector) error {
	// DO NOT CHANGE THE ORDER. Please refer to `Run()` at `central/sensor/service/connection/connection_impl.go`
	if err := s.initialConfigSync(stream, configHandler); err != nil {
		return err
	}

	return s.initialPolicySync(stream, detector)
}

func (s *centralCommunicationImpl) initialConfigSync(stream central.SensorService_CommunicateClient, handler config.Handler) error {
	headerMD, err := stream.Header()
	if err != nil {
		return errors.Wrap(err, "receiving header metadata from central")
	}

	helmManagedCfg := handler.GetHelmManagedConfig()

	if metautils.NiceMD(headerMD).Get(centralsensor.HelmManagedClusterMetadataKey) == "true" {
		if helmManagedCfg == nil {
			return errors.New("central requested Helm-managed cluster config, but no Helm-managed config is available")
		}

		if helmManagedCfg.GetClusterId() == "" {
			cachedClusterID, err := helmconfig.LoadCachedClusterID()
			if err != nil {
				log.Warnf("Failed to load cached cluster ID: %s", err)
			} else if cachedClusterID != "" {
				helmManagedCfg = helmManagedCfg.Clone()
				helmManagedCfg.ClusterId = cachedClusterID
				log.Infof("Re-using cluster ID %s of previous run. If this is causing issues, re-apply a new Helm configuration via 'helm upgrade', or delete the sensor pod.", cachedClusterID)
			}
		}

		msg := &central.MsgFromSensor{
			Msg: &central.MsgFromSensor_HelmManagedConfigInit{
				HelmManagedConfigInit: helmManagedCfg,
			},
		}
		if err := stream.Send(msg); err != nil {
			return errors.Wrap(err, "could not send Helm-managed cluster config")
		}
	} else if helmManagedCfg != nil {
		log.Warn("Central instance does NOT support Helm-managed configuration. Dynamic cluster configuration MUST be changed via the UI in order to take effect. Please upgrade Central to a recent version to allow changing dynamic cluster configuration via 'helm upgrade'")
	}

	msg, err := stream.Recv()
	if err != nil {
		return errors.Wrap(err, "receiving initial cluster config")
	}
	if msg.GetClusterConfig() == nil {
		return errors.Errorf("initial message received from Sensor was not a cluster config: %T", msg.Msg)
	}
	if features.SensorInstallationExperience.Enabled() {
		clusterID := msg.GetClusterConfig().GetClusterId()
		clusterid.Set(clusterID)
		if err := helmconfig.StoreCachedClusterID(clusterID); err != nil {
			log.Warnf("Could not cache cluster ID: %v", err)
		}
	}
	// Send the initial cluster config to the config handler
	if err := handler.ProcessMessage(msg); err != nil {
		return errors.Wrap(err, "processing initial cluster config")
	}
	return nil
}

func (s *centralCommunicationImpl) initialPolicySync(stream central.SensorService_CommunicateClient, detector detector.Detector) error {
	// Policy sync
	msg, err := stream.Recv()
	if err != nil {
		return errors.Wrap(err, "receiving initial policies")
	}
	if msg.GetPolicySync() == nil {
		return errors.Errorf("second message received from Sensor was not a policy sync: %T", msg.Msg)
	}
	if err := detector.ProcessMessage(msg); err != nil {
		return errors.Wrap(err, "policy sync could not be successfully processed")
	}

	// Process baselines sync
	msg, err = stream.Recv()
	if err != nil {
		return errors.Wrap(err, "receiving initial baselines")
	}
	if err := detector.ProcessMessage(msg); err != nil {
		return errors.Wrap(err, "process baselines could not be successfully processed")
	}
	return nil
}

func runAll(err error, fs ...func(error)) {
	for _, f := range fs {
		f(err)
	}
}
