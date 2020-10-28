package service

import (
	"context"

	clusterDataStore "github.com/stackrox/rox/central/cluster/datastore"
	deploymentDataStore "github.com/stackrox/rox/central/deployment/datastore"
	nsDS "github.com/stackrox/rox/central/namespace/datastore"
	networkEntityDS "github.com/stackrox/rox/central/networkgraph/entity/datastore"
	nfDS "github.com/stackrox/rox/central/networkgraph/flow/datastore"
	npDS "github.com/stackrox/rox/central/networkpolicies/datastore"
	"github.com/stackrox/rox/central/networkpolicies/generator"
	"github.com/stackrox/rox/central/networkpolicies/graph"
	notifierDataStore "github.com/stackrox/rox/central/notifier/datastore"
	"github.com/stackrox/rox/central/sensor/service/connection"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/grpc"
	"github.com/stackrox/rox/pkg/logging"
)

var (
	log = logging.LoggerForModule()
)

// Service provides the interface to the microservice that serves alert data.
type Service interface {
	grpc.APIService

	AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error)

	v1.NetworkPolicyServiceServer
}

// New returns a new Service instance using the given DataStore.
func New(storage npDS.DataStore, deployments deploymentDataStore.DataStore, externalSrcs networkEntityDS.EntityDataStore, graphEvaluator graph.Evaluator, namespacesStore nsDS.DataStore, clusterStore clusterDataStore.DataStore, notifierStore notifierDataStore.DataStore, globalFlowDataStore nfDS.ClusterDataStore, sensorConnMgr connection.Manager) Service {
	return &serviceImpl{
		sensorConnMgr:   sensorConnMgr,
		deployments:     deployments,
		externalSrcs:    externalSrcs,
		networkPolicies: storage,
		notifierStore:   notifierStore,
		clusterStore:    clusterStore,
		graphEvaluator:  graphEvaluator,
		policyGenerator: generator.New(storage, deployments, namespacesStore, globalFlowDataStore, externalSrcs),
	}
}
