package service

import (
	"context"

	"github.com/stackrox/rox/central/continuousintegration/datastore"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/grpc"
)

// Service provides the interface to the gRPC service for users.
type Service interface {
	grpc.APIService

	AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error)

	v1.ContinuousIntegrationServiceServer
}

func New(dataStore datastore.DataStore) Service {
	return &serviceImpl{
		dataStore: dataStore,
	}
}
