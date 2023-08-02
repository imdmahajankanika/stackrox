package scannerclient

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	scannerV4 "github.com/stackrox/rox/generated/internalapi/scanner/v4"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/clientconn"
	"github.com/stackrox/rox/pkg/mtls"
	"github.com/stackrox/rox/pkg/registries/types"
	scannerV1 "github.com/stackrox/scanner/generated/scanner/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var _ Client = (*V4GRPCClient)(nil)
var languageComponents = map[string]bool{
	"Go":       true,
	"Maven":    true,
	"PiPI":     true,
	"RubyGems": true,
}

// V4GRPCClient represents a client implementation using the v4 gRPC protocol.
type V4GRPCClient struct {
	indexerClient scannerV4.IndexerClient
	matcherClient scannerV4.MatcherClient
	conn          *grpc.ClientConn
}

func (v V4GRPCClient) Dial(endpoint string) (Client, error) {
	if endpoint == "" {
		return nil, errors.New("Invalid Scanner endpoint (empty)")
	}

	endpoint = strings.TrimPrefix(endpoint, "https://")
	if strings.Contains(endpoint, "://") {
		return nil, errors.Errorf("ScannerV4 endpoint has unsupported scheme: %s", endpoint)
	}

	tlsConfig, err := clientconn.TLSConfig(mtls.ScannerSubject, clientconn.TLSConfigOptions{
		UseClientCert: clientconn.MustUseClientCert,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize ScannerV4 TLS config")
	}

	// This is non-blocking. If we ever want this to block,
	// then add the grpc.WithBlock() DialOption.
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial Scanner")
	}

	log.Infof("Dialing ScannerV4 at %s", endpoint)

	return &V4GRPCClient{
		indexerClient: scannerV4.NewIndexerClient(conn),
		matcherClient: scannerV4.NewMatcherClient(conn),
		conn:          conn,
	}, nil
}

func (v V4GRPCClient) GetImageAnalysis(ctx context.Context, image *storage.Image, cfg *types.Config) (*scannerV1.GetImageComponentsResponse, *scannerV4.IndexReport, error) {
	name := image.GetName().GetFullName()
	var hid string
	if hid = image.Metadata.V2.Digest; len(hid) == 0 {
		hid = image.Metadata.V1.Digest
	}
	indexReport, err := v.indexerClient.CreateIndexReport(ctx, &scannerV4.CreateIndexReportRequest{
		HashId:               hid,
		ResourceLocator:      nil,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}, grpc.WaitForReady(true))

	if err != nil {
		log.Debugf("Unable to get image components from local Scanner for image %s: %v", name, err)
		return nil, nil, errors.Wrap(err, "getting image components from scanner")
	}

	log.Debugf("Received image indexing report from local Scanner for image %s", name)

	// Convert indexReport to scannerV1.GetImageComponentsResponse
	res := &scannerV1.GetImageComponentsResponse{}
	if indexReport.Success {
		res.Status = scannerV1.ScanStatus_SUCCEEDED
		res.ScannerVersion = image.GetScan().ScannerVersion
	} else {
		// Handle conversion failure
		err = errors.New("indexer report conversion failed")
		log.Debugf("Failed to convert indexer report to image components from local Scanner for image %s", name)
		return nil, nil, errors.Wrap(err, "converting indexer report from scanner")
	}
	//return resp or return indexReport directly?
	if err != nil {
		log.Debugf("Failed to convert indexer report to image components from local Scanner for image %s: %v", name, err)
		return nil, nil, errors.Wrap(err, "converting indexer report from scanner")
	}

	return res, indexReport, nil
}

func getNamespace(report scannerV4.IndexReport) (string, bool) {
	if len(report.Distributions) == 1 {
		dist := report.Distributions[0]
		os := dist.Name + ":" + dist.Version
		return os, dist.Did == "rhel"
	}
	return "unknown", false
}

func (v V4GRPCClient) Close() error {
	return v.conn.Close()
}
