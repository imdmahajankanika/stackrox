package index

import (
	"time"

	"github.com/blevesearch/bleve"
	"github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/central/policy/index/mappings"
	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/blevesearch"
)

// AlertIndex provides storage functionality for alerts.
type indexerImpl struct {
	index bleve.Index
}

type policyWrapper struct {
	*storage.Policy `json:"policy"`
	Type            string `json:"type"`
}

// AddPolicy adds the policy to the index
func (b *indexerImpl) AddPolicy(policy *storage.Policy) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Add, "Policy")
	return b.index.Index(policy.GetId(), &policyWrapper{Type: v1.SearchCategory_POLICIES.String(), Policy: policy})
}

// AddPolicies adds the policies to the indexer
func (b *indexerImpl) AddPolicies(policies []*storage.Policy) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.AddMany, "Policy")
	batch := b.index.NewBatch()
	for _, policy := range policies {
		batch.Index(policy.GetId(), &policyWrapper{Type: v1.SearchCategory_POLICIES.String(), Policy: policy})
	}
	return b.index.Batch(batch)
}

// DeletePolicy deletes the policy from the index
func (b *indexerImpl) DeletePolicy(id string) error {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Remove, "Policy")
	return b.index.Delete(id)
}

// SearchPolicies takes a SearchRequest and finds any matches
func (b *indexerImpl) SearchPolicies(q *v1.Query) ([]search.Result, error) {
	defer metrics.SetIndexOperationDurationTime(time.Now(), ops.Search, "Policy")
	return blevesearch.RunSearchRequest(v1.SearchCategory_POLICIES, q, b.index, mappings.OptionsMap)
}
