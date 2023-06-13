package index

import (
	"context"

	v1 "github.com/stackrox/rox/generated/api/v1"
	storage "github.com/stackrox/rox/generated/storage"
	search "github.com/stackrox/rox/pkg/search"
	blevesearch "github.com/stackrox/rox/pkg/search/blevesearch"
)

// Indexer is the image-component edge indexer.
//
//go:generate mockgen-wrapper
type Indexer interface {
	AddImageComponentEdge(imagecomponentedge *storage.ImageComponentEdge) error
	AddImageComponentEdges(imagecomponentedges []*storage.ImageComponentEdge) error
	Count(ctx context.Context, q *v1.Query, opts ...blevesearch.SearchOption) (int, error)
	DeleteImageComponentEdge(id string) error
	DeleteImageComponentEdges(ids []string) error
	Search(ctx context.Context, q *v1.Query, opts ...blevesearch.SearchOption) ([]search.Result, error)
}
