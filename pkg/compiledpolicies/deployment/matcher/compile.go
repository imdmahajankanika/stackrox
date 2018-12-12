package matcher

import (
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/compiledpolicies/deployment/predicate"
)

// Compiler is a function that turns a policy into a matcher.
type compiler func(*storage.Policy) (Matcher, error)

// compilers are all of the different Matcher creation functions that are registered.
var compilers []compiler

// Compile creates a new deployment policy matcher for build and deploy policies.
func Compile(policy *storage.Policy) (Matcher, error) {
	var matcher Matcher
	for _, compiler := range compilers {
		matcherFunction, err := compiler(policy)
		if err != nil {
			return nil, err
		}
		matcher = matcher.MustAlsoViolate(matcherFunction)
	}
	if matcher == nil {
		return nil, nil
	}

	pred, err := predicate.Compile(policy)
	if err != nil {
		return nil, err
	}
	return matcher.ProcessIf(pred), nil
}
