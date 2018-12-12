package predicate

import (
	"github.com/stackrox/rox/generated/storage"
)

func init() {
	compilers = append(compilers, newDisabledPredicate)
}

func newDisabledPredicate(policy *storage.Policy) (Predicate, error) {
	if !policy.GetDisabled() {
		return nil, nil
	}
	return shouldProcess, nil
}

// If the policy is disabled, we create a predicate that always returns false.
func shouldProcess(*storage.Deployment) bool {
	return false
}
