package check121

import (
	"context"
	"testing"

	"github.com/stackrox/rox/central/compliance/framework"
	"github.com/stackrox/rox/central/compliance/framework/mocks"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/networkgraph"
	"github.com/stackrox/rox/pkg/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestCheck(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(suiteImpl))
}

type suiteImpl struct {
	suite.Suite

	mockCtrl *gomock.Controller
}

func (s *suiteImpl) SetupSuite() {
	s.mockCtrl = gomock.NewController(s.T())
}

func (s *suiteImpl) TearDownSuite() {
	s.mockCtrl.Finish()
}

func (s *suiteImpl) TestHostNetwork() {
	check := s.verifyCheckRegistered()

	testCluster := s.cluster()

	testDeployments := []*storage.Deployment{
		{
			Id:          uuid.NewV4().String(),
			HostNetwork: true,
		},
	}

	testNodes := s.nodes()

	testPolicies := s.networkPolicies()

	testNetworkGraph := &v1.NetworkGraph{
		Nodes: []*v1.NetworkNode{
			{
				Entity:    networkgraph.EntityForDeployment(testDeployments[0].GetId()).ToProto(),
				PolicyIds: []string{testPolicies[0].GetId(), testPolicies[1].GetId()},
			},
		},
	}

	data := mocks.NewMockComplianceDataRepository(s.mockCtrl)
	data.EXPECT().NetworkPolicies().AnyTimes().Return(toMap(testPolicies))
	data.EXPECT().NetworkGraph().AnyTimes().Return(testNetworkGraph)

	run, err := framework.NewComplianceRun(check)
	s.NoError(err)

	domain := framework.NewComplianceDomain(testCluster, testNodes, testDeployments, nil, nil)
	err = run.Run(context.Background(), "standard", domain, data)
	s.NoError(err)

	results := run.GetAllResults()
	checkResults := results[checkID]
	s.NotNil(checkResults)

	for _, deployment := range domain.Deployments() {
		deploymentResults := checkResults.ForChild(deployment)
		s.NoError(deploymentResults.Error())
		s.Len(deploymentResults.Evidence(), 1)
		s.Equal(framework.FailStatus, deploymentResults.Evidence()[0].Status)
	}
}

func (s *suiteImpl) TestMissingIngress() {
	check := s.verifyCheckRegistered()

	testCluster := s.cluster()

	testDeployments := []*storage.Deployment{
		{
			Id:          uuid.NewV4().String(),
			HostNetwork: false,
		},
	}

	testNodes := s.nodes()

	testPolicies := s.networkPolicies()

	testNetworkGraph := &v1.NetworkGraph{
		Nodes: []*v1.NetworkNode{
			{
				Entity:    networkgraph.EntityForDeployment(testDeployments[0].GetId()).ToProto(),
				PolicyIds: []string{testPolicies[1].GetId()}, // Only egress
			},
		},
	}

	data := mocks.NewMockComplianceDataRepository(s.mockCtrl)
	data.EXPECT().NetworkPolicies().AnyTimes().Return(toMap(testPolicies))
	data.EXPECT().NetworkGraph().AnyTimes().Return(testNetworkGraph)

	run, err := framework.NewComplianceRun(check)
	s.NoError(err)

	domain := framework.NewComplianceDomain(testCluster, testNodes, testDeployments, nil, nil)
	err = run.Run(context.Background(), "standard", domain, data)
	s.NoError(err)

	results := run.GetAllResults()
	checkResults := results[checkID]
	s.NotNil(checkResults)

	for _, deployment := range domain.Deployments() {
		deploymentResults := checkResults.ForChild(deployment)
		s.NoError(deploymentResults.Error())
		s.Len(deploymentResults.Evidence(), 1)
		s.Equal(framework.FailStatus, deploymentResults.Evidence()[0].Status)
	}
}

func (s *suiteImpl) TestMissingEgress() {
	check := s.verifyCheckRegistered()

	testCluster := s.cluster()

	testDeployments := []*storage.Deployment{
		{
			Id:          uuid.NewV4().String(),
			HostNetwork: false,
		},
	}

	testNodes := s.nodes()

	testPolicies := s.networkPolicies()

	testNetworkGraph := &v1.NetworkGraph{
		Nodes: []*v1.NetworkNode{
			{
				Entity:    networkgraph.EntityForDeployment(testDeployments[0].GetId()).ToProto(),
				PolicyIds: []string{testPolicies[0].GetId()}, // Only ingress
			},
		},
	}

	data := mocks.NewMockComplianceDataRepository(s.mockCtrl)
	data.EXPECT().NetworkPolicies().AnyTimes().Return(toMap(testPolicies))
	data.EXPECT().NetworkGraph().AnyTimes().Return(testNetworkGraph)

	run, err := framework.NewComplianceRun(check)
	s.NoError(err)

	domain := framework.NewComplianceDomain(testCluster, testNodes, testDeployments, nil, nil)
	err = run.Run(context.Background(), "standard", domain, data)
	s.NoError(err)

	results := run.GetAllResults()
	checkResults := results[checkID]
	s.NotNil(checkResults)

	for _, deployment := range domain.Deployments() {
		deploymentResults := checkResults.ForChild(deployment)
		s.NoError(deploymentResults.Error())
		s.Len(deploymentResults.Evidence(), 1)
		s.Equal(framework.FailStatus, deploymentResults.Evidence()[0].Status)
	}
}

func (s *suiteImpl) TestSkipKubeSystem() {
	check := s.verifyCheckRegistered()

	testCluster := s.cluster()

	testDeployments := []*storage.Deployment{
		{
			Id:          uuid.NewV4().String(),
			HostNetwork: true,
			Namespace:   "kube-system",
		},
	}

	testNodes := s.nodes()
	testPolicies := s.networkPolicies()
	testNetworkGraph := &v1.NetworkGraph{}

	data := mocks.NewMockComplianceDataRepository(s.mockCtrl)
	data.EXPECT().NetworkPolicies().AnyTimes().Return(toMap(testPolicies))
	data.EXPECT().NetworkGraph().AnyTimes().Return(testNetworkGraph)

	run, err := framework.NewComplianceRun(check)
	s.NoError(err)

	domain := framework.NewComplianceDomain(testCluster, testNodes, testDeployments, nil, nil)
	err = run.Run(context.Background(), "standard", domain, data)
	s.NoError(err)

	results := run.GetAllResults()
	checkResults := results[checkID]
	s.NotNil(checkResults)

	for _, deployment := range domain.Deployments() {
		deploymentResults := checkResults.ForChild(deployment)
		s.NoError(deploymentResults.Error())
		if s.Len(deploymentResults.Evidence(), 1) {
			s.Equal(framework.SkipStatus, deploymentResults.Evidence()[0].Status)
		}
	}
}

func (s *suiteImpl) TestPass() {
	check := s.verifyCheckRegistered()

	testCluster := s.cluster()

	testDeployments := []*storage.Deployment{
		{
			Id:          uuid.NewV4().String(),
			HostNetwork: false,
		},
		{
			Id:          uuid.NewV4().String(),
			HostNetwork: false,
		},
	}

	testNodes := s.nodes()

	testPolicies := s.networkPolicies()

	testNetworkGraph := &v1.NetworkGraph{
		Nodes: []*v1.NetworkNode{
			{
				Entity:    networkgraph.EntityForDeployment(testDeployments[0].GetId()).ToProto(),
				PolicyIds: []string{testPolicies[0].GetId(), testPolicies[1].GetId()},
			},
			{
				Entity:    networkgraph.EntityForDeployment(testDeployments[1].GetId()).ToProto(),
				PolicyIds: []string{testPolicies[0].GetId(), testPolicies[1].GetId()},
			},
		},
	}

	data := mocks.NewMockComplianceDataRepository(s.mockCtrl)
	data.EXPECT().NetworkPolicies().AnyTimes().Return(toMap(testPolicies))
	data.EXPECT().NetworkGraph().AnyTimes().Return(testNetworkGraph)

	run, err := framework.NewComplianceRun(check)
	s.NoError(err)

	domain := framework.NewComplianceDomain(testCluster, testNodes, testDeployments, nil, nil)
	err = run.Run(context.Background(), "standard", domain, data)
	s.NoError(err)

	results := run.GetAllResults()
	checkResults := results[checkID]
	s.NotNil(checkResults)

	for _, deployment := range domain.Deployments() {
		deploymentResults := checkResults.ForChild(deployment)
		s.NoError(deploymentResults.Error())
		s.Len(deploymentResults.Evidence(), 1)
		s.Equal(framework.PassStatus, deploymentResults.Evidence()[0].Status)
	}
}

// Helper functions for test data.
//////////////////////////////////

func (s *suiteImpl) verifyCheckRegistered() framework.Check {
	registry := framework.RegistrySingleton()
	check := registry.Lookup(checkID)
	s.NotNil(check)
	return check
}

func (s *suiteImpl) cluster() *storage.Cluster {
	return &storage.Cluster{
		Id: uuid.NewV4().String(),
	}
}

func (s *suiteImpl) networkPolicies() []*storage.NetworkPolicy {
	return []*storage.NetworkPolicy{
		{
			Id: uuid.NewV4().String(),
			Spec: &storage.NetworkPolicySpec{
				PolicyTypes: []storage.NetworkPolicyType{
					storage.NetworkPolicyType_INGRESS_NETWORK_POLICY_TYPE,
				},
				Ingress: []*storage.NetworkPolicyIngressRule{
					{},
				},
			},
		},
		{
			Id: uuid.NewV4().String(),
			Spec: &storage.NetworkPolicySpec{
				PolicyTypes: []storage.NetworkPolicyType{
					storage.NetworkPolicyType_EGRESS_NETWORK_POLICY_TYPE,
				},
				Egress: []*storage.NetworkPolicyEgressRule{
					{},
				},
			},
		},
	}
}

func (s *suiteImpl) nodes() []*storage.Node {
	return []*storage.Node{
		{
			Id: uuid.NewV4().String(),
		},
		{
			Id: uuid.NewV4().String(),
		},
	}
}

func toMap(in []*storage.NetworkPolicy) map[string]*storage.NetworkPolicy {
	merp := make(map[string]*storage.NetworkPolicy, len(in))
	for _, np := range in {
		merp[np.GetId()] = np
	}
	return merp
}
