package runtime

import (
	"github.com/stackrox/rox/central/deployment/datastore"
	"github.com/stackrox/rox/central/detection/deployment"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
)

var (
	log = logging.LoggerForModule()
)

// Detector provides an interface for performing runtime policy violation detection.
type Detector interface {
	AlertsForDeployments(deploymentIDs ...string) ([]*storage.Alert, error)
	AlertsForPolicy(policyID string) ([]*storage.Alert, error)
	DeploymentWhitelistedForPolicy(deploymentID, policyID string) bool
	UpsertPolicy(policy *storage.Policy) error
	RemovePolicy(policyID string) error
}

// NewDetector returns a new instance of a Detector.
func NewDetector(policySet deployment.PolicySet, deployments datastore.DataStore) Detector {
	return &detectorImpl{
		policySet:   policySet,
		deployments: deployments,
	}
}
