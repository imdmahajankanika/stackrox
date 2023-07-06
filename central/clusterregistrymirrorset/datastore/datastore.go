package datastore

import (
	"context"

	"github.com/stackrox/rox/central/clusterregistrymirrorset/store"
	"github.com/stackrox/rox/central/role/resources"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/sac"
)

var (
	clusterSAC = sac.ForResource(resources.Cluster)
)

// DataStore is the entry point for modifying cluster registry mirror sets.
type DataStore interface {
	UpsertMirror(context.Context, *storage.ClusterRegistryMirrorSet) error
}

// New returns an instance of DataStore.
func New(store store.Store) DataStore {
	return &datastoreImpl{
		store: store,
	}
}

type datastoreImpl struct {
	store store.Store
}

// UpsertMirror upserts mirrors into the datastore
func (d *datastoreImpl) UpsertMirror(ctx context.Context, mirrorSet *storage.ClusterRegistryMirrorSet) error {
	// TODO: Confirm if sac still necessary/required?
	if ok, err := clusterSAC.WriteAllowed(ctx); err != nil {
		return err
	} else if !ok {
		return sac.ErrResourceAccessDenied
	}

	return d.store.Upsert(ctx, mirrorSet)
}
