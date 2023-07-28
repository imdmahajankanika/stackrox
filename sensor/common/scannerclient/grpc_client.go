package scannerclient

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	scannerV4 "github.com/stackrox/rox/generated/internalapi/scanner/v4"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/clientconn"
	"github.com/stackrox/rox/pkg/images/utils"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/mtls"
	"github.com/stackrox/rox/pkg/registries/types"
	scannerV1 "github.com/stackrox/scanner/generated/scanner/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	log = logging.LoggerForModule()
)

// Client represents a gRPC client interface.
type Client interface {
	Dial(endpoint string) (Client, error)
	GetImageAnalysis(ctx context.Context, image *storage.Image, cfg *types.Config) (*scannerV1.GetImageComponentsResponse, *scannerV4.IndexReport, error)
	Close() error
}

// GrpcClient is a Scanner gRPC GrpcClient.
type GrpcClient struct {
	client scannerV1.ImageScanServiceClient
	conn   *grpc.ClientConn
}

// Dial Scanner and return a new GrpcClient.
// Dial is non-blocking and returns a non-nil error upon configuration error.
func (c *GrpcClient) Dial(endpoint string) (Client, error) {
	if endpoint == "" {
		return nil, errors.New("Invalid Scanner endpoint (empty)")
	}

	endpoint = strings.TrimPrefix(endpoint, "https://")
	if strings.Contains(endpoint, "://") {
		return nil, errors.Errorf("Scanner endpoint has unsupported scheme: %s", endpoint)
	}

	tlsConfig, err := clientconn.TLSConfig(mtls.ScannerSubject, clientconn.TLSConfigOptions{
		UseClientCert: clientconn.MustUseClientCert,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize Scanner TLS config")
	}

	// This is non-blocking. If we ever want this to block,
	// then add the grpc.WithBlock() DialOption.
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to Dial Scanner")
	}

	log.Infof("Dialing Scanner at %s", endpoint)

	return &GrpcClient{
		client: scannerV1.NewImageScanServiceClient(conn),
		conn:   conn,
	}, nil
}

// GetImageAnalysis retrieves the image analysis results for the given image.
func (c *GrpcClient) GetImageAnalysis(ctx context.Context, image *storage.Image, cfg *types.Config) (*scannerV1.GetImageComponentsResponse, *scannerV4.IndexReport, error) {
	name := image.GetName().GetFullName()

	// The WaitForReady option will cause invocations to block (until server ready or ctx done/expires)
	// This was added so that on fresh install of sensor when scanner is not ready yet, local scans will
	// not all fail and have to wait for next reprocess to succeed
	resp, err := c.client.GetImageComponents(ctx, &scannerV1.GetImageComponentsRequest{
		Image: utils.GetFullyQualifiedFullName(image),
		Registry: &scannerV1.RegistryData{
			Url:      cfg.URL,
			Username: cfg.Username,
			Password: cfg.Password,
			Insecure: cfg.Insecure,
		},
	}, grpc.WaitForReady(true))
	if err != nil {
		log.Debugf("Unable to get image components from local Scanner for image %s: %v", name, err)
		return nil, nil, errors.Wrap(err, "getting image components from scanner")
	}

	log.Debugf("Received image components from local Scanner for image %s", name)

	return resp, nil, nil
}

// Close closes the underlying grpc.ClientConn.
func (c *GrpcClient) Close() error {
	return c.conn.Close()
}
