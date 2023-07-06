package store

import (
	"context"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
)

type Store interface {
	Upsert(ctx context.Context, obj *storage.ClusterRegistryMirrorSet) error
}

// Store provides an interface to the underlying data layer
type UnderlyingStore interface {
	Upsert(ctx context.Context, obj *storage.ClusterRegistryMirrorSet) error
	UpsertMany(ctx context.Context, objs []*storage.ClusterRegistryMirrorSet) error
	Delete(ctx context.Context, id string) error
	DeleteByQuery(ctx context.Context, q *v1.Query) error
	DeleteMany(ctx context.Context, identifiers []string) error

	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, id string) (bool, error)

	Get(ctx context.Context, id string) (*storage.ClusterRegistryMirrorSet, bool, error)
	GetMany(ctx context.Context, identifiers []string) ([]*storage.ClusterRegistryMirrorSet, []int, error)
	GetIDs(ctx context.Context) ([]string, error)

	Walk(ctx context.Context, fn func(obj *storage.ClusterRegistryMirrorSet) error) error
}

type storeImpl struct {
	store UnderlyingStore
}

// NewStore returns a wrapper store for cluster registry mirrors
func NewStore(store UnderlyingStore) Store {
	return &storeImpl{store: store}
}

func (s *storeImpl) Upsert(ctx context.Context, obj *storage.ClusterRegistryMirrorSet) error {
	return s.store.Upsert(ctx, obj)
}
