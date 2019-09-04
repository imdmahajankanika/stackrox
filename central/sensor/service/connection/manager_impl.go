package connection

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/central/sensor/service/connection/upgradecontroller"
	"github.com/stackrox/rox/central/sensor/service/pipeline"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/sync"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	clusterCheckinInterval = 30 * time.Second
)

var (
	managerCtx = sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(
			sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
			sac.ResourceScopeKeys(resources.Cluster)))

	clusterSAC = sac.ForResource(resources.Cluster)
)

type connectionAndUpgradeController struct {
	connection  *sensorConnection
	upgradeCtrl upgradecontroller.UpgradeController
}

type manager struct {
	connectionsByClusterID      map[string]connectionAndUpgradeController
	connectionsByClusterIDMutex sync.RWMutex

	clusters ClusterManager
}

func newManager() *manager {
	return &manager{
		connectionsByClusterID: make(map[string]connectionAndUpgradeController),
	}
}

func (m *manager) initializeUpgradeControllers() error {
	clusters, err := m.clusters.GetClusters(managerCtx)
	if err != nil {
		return err
	}
	m.connectionsByClusterIDMutex.Lock()
	defer m.connectionsByClusterIDMutex.Unlock()
	for _, cluster := range clusters {
		upgradeCtrl, err := upgradecontroller.New(cluster.GetId(), m.clusters)
		if err != nil {
			return err
		}
		m.connectionsByClusterID[cluster.GetId()] = connectionAndUpgradeController{
			upgradeCtrl: upgradeCtrl,
		}
	}
	return nil
}

func (m *manager) Start(clusterManager ClusterManager) error {
	m.clusters = clusterManager
	if features.SensorAutoUpgrade.Enabled() {
		err := m.initializeUpgradeControllers()
		if err != nil {
			return errors.Wrap(err, "failed to initialize upgrade controllers")
		}
	}

	go m.updateClusterContactTimesForever()
	return nil
}

func (m *manager) updateClusterContactTimesForever() {
	t := time.NewTicker(clusterCheckinInterval)
	defer t.Stop()

	for range t.C {
		connections := m.GetActiveConnections()
		clusterIDs := make([]string, 0, len(connections))
		for _, c := range connections {
			clusterIDs = append(clusterIDs, c.ClusterID())
		}
		if err := m.clusters.UpdateClusterContactTimes(managerCtx, time.Now(), clusterIDs...); err != nil {
			log.Errorf("error checking in clusters: %v", err)
		}
	}
}

func (m *manager) GetConnection(clusterID string) SensorConnection {
	m.connectionsByClusterIDMutex.RLock()
	defer m.connectionsByClusterIDMutex.RUnlock()

	conn := m.connectionsByClusterID[clusterID].connection
	if conn == nil {
		return nil
	}
	return conn
}

func (m *manager) replaceConnection(ctx context.Context, clusterID string, newConnection *sensorConnection) (oldConnection *sensorConnection, err error) {
	m.connectionsByClusterIDMutex.Lock()
	defer m.connectionsByClusterIDMutex.Unlock()

	connAndUpgradeCtrl := m.connectionsByClusterID[clusterID]
	oldConnection = connAndUpgradeCtrl.connection
	upgradeCtrl := connAndUpgradeCtrl.upgradeCtrl

	if features.SensorAutoUpgrade.Enabled() {
		if upgradeCtrl == nil {
			upgradeCtrl, err = upgradecontroller.New(clusterID, m.clusters)
			if err != nil {
				return nil, err
			}
			go newConnection.stopSig.SignalWhen(upgradeCtrl.ErrorSignal(), concurrency.Never())
		}
		upgradeCtrl.RegisterConnection(ctx, newConnection)
	}
	m.connectionsByClusterID[clusterID] = connectionAndUpgradeController{
		connection:  newConnection,
		upgradeCtrl: upgradeCtrl,
	}
	return oldConnection, nil
}

func (m *manager) HandleConnection(ctx context.Context, clusterID string, pf pipeline.Factory, server central.SensorService_CommunicateServer) error {
	conn, err := newConnection(ctx, clusterID, pf, m.clusters)
	if err != nil {
		return errors.Wrap(err, "creating sensor connection")
	}

	oldConnection, err := m.replaceConnection(ctx, clusterID, conn)
	if err != nil {
		return errors.Wrap(err, "replacing old connection")
	}

	if oldConnection != nil {
		oldConnection.Terminate(errors.New("replaced by new connection"))
	}

	err = conn.Run(ctx, server)
	log.Warnf("Connection to server in cluster %s terminated: %v", clusterID, err)

	concurrency.WithLock(&m.connectionsByClusterIDMutex, func() {
		connAndUpgradeCtrl := m.connectionsByClusterID[clusterID]
		if connAndUpgradeCtrl.connection == conn {
			connAndUpgradeCtrl.connection = nil
			m.connectionsByClusterID[clusterID] = connAndUpgradeCtrl
		}
	})

	return err
}

func (m *manager) getOrCreateUpgradeCtrl(clusterID string) (upgradecontroller.UpgradeController, error) {
	m.connectionsByClusterIDMutex.Lock()
	defer m.connectionsByClusterIDMutex.Unlock()

	connAndUpgradeCtrl := m.connectionsByClusterID[clusterID]
	if connAndUpgradeCtrl.upgradeCtrl == nil {
		var err error
		connAndUpgradeCtrl.upgradeCtrl, err = upgradecontroller.New(clusterID, m.clusters)
		if err != nil {
			return nil, err
		}
		m.connectionsByClusterID[clusterID] = connAndUpgradeCtrl
	}
	return connAndUpgradeCtrl.upgradeCtrl, nil
}

func (m *manager) RecordUpgradeProgress(clusterID, upgradeProcessID string, upgradeProgress *storage.UpgradeProgress) error {
	if !features.SensorAutoUpgrade.Enabled() {
		return errors.New("cannot record upgrade progress; auto-upgrade feature flag disabled")
	}
	upgradeCtrl, err := m.getOrCreateUpgradeCtrl(clusterID)
	if err != nil {
		return err
	}

	return upgradeCtrl.RecordUpgradeProgress(upgradeProcessID, upgradeProgress)
}

func (m *manager) TriggerUpgrade(ctx context.Context, clusterID string) error {
	if ok, err := clusterSAC.WriteAllowed(ctx, sac.ClusterScopeKey(clusterID)); err != nil {
		return status.Error(codes.Internal, err.Error())
	} else if !ok {
		return status.Error(codes.PermissionDenied, sac.ErrPermissionDenied.Error())
	}

	var upgradeCtrl upgradecontroller.UpgradeController
	concurrency.WithRLock(&m.connectionsByClusterIDMutex, func() {
		upgradeCtrl = m.connectionsByClusterID[clusterID].upgradeCtrl
	})
	if upgradeCtrl == nil {
		return errors.Errorf("no upgrade controller found for cluster ID %s; either the sensor has not checked in or the clusterID is invalid. Cannot trigger upgrade", clusterID)
	}
	return upgradeCtrl.Trigger(ctx)
}

func (m *manager) GetActiveConnections() []SensorConnection {
	m.connectionsByClusterIDMutex.RLock()
	defer m.connectionsByClusterIDMutex.RUnlock()

	result := make([]SensorConnection, 0, len(m.connectionsByClusterID))

	for _, connAndUpgradeCtrl := range m.connectionsByClusterID {
		if conn := connAndUpgradeCtrl.connection; conn != nil {
			result = append(result, conn)
		}
	}

	return result
}
