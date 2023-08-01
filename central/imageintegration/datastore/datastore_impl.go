package datastore

import (
	"context"

	"github.com/stackrox/rox/central/imageintegration/search"
	"github.com/stackrox/rox/central/imageintegration/store"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/sac/resources"
	searchPkg "github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/uuid"
)

var (
	integrationSAC = sac.ForResource(resources.Integration)
)

type datastoreImpl struct {
	storage           store.Store
	formattedSearcher search.Searcher
}

func (ds *datastoreImpl) Count(ctx context.Context, q *v1.Query) (int, error) {
	return ds.formattedSearcher.Count(ctx, q)
}

func (ds *datastoreImpl) Search(ctx context.Context, q *v1.Query) ([]searchPkg.Result, error) {
	return ds.formattedSearcher.Search(ctx, q)
}

// GetImageIntegration is pass-through to the underlying store.
func (ds *datastoreImpl) GetImageIntegration(ctx context.Context, id string) (*storage.ImageIntegration, bool, error) {
	if ok, err := integrationSAC.ReadAllowed(ctx); err != nil {
		return nil, false, err
	} else if !ok {
		return nil, false, nil
	}

	return ds.storage.Get(ctx, id)
}

// GetImageIntegrations provides an in memory layer on top of the underlying DB based storage.
func (ds *datastoreImpl) GetImageIntegrations(ctx context.Context, request *v1.GetImageIntegrationsRequest) ([]*storage.ImageIntegration, error) {
	if ok, err := integrationSAC.ReadAllowed(ctx); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	integrations, err := ds.storage.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	integrationSlice := integrations[:0]
	for _, integration := range integrations {
		if request.GetCluster() != "" {
			continue
		}
		if request.GetName() != "" && request.GetName() != integration.GetName() {
			continue
		}
		integrationSlice = append(integrationSlice, integration)
	}
	return integrationSlice, nil
}

// AddImageIntegration is pass-through to the underlying store.
func (ds *datastoreImpl) AddImageIntegration(ctx context.Context, integration *storage.ImageIntegration) (string, error) {
	if ok, err := integrationSAC.WriteAllowed(ctx); err != nil {
		return "", err
	} else if !ok {
		return "", sac.ErrResourceAccessDenied
	}

	if !features.SourcedAutogeneratedIntegrations.Enabled() || integration.GetId() == "" {
		integration.Id = uuid.NewV4().String()
	}
	err := ds.storage.Upsert(ctx, integration)
	if err != nil {
		return "", err
	}
	return integration.Id, nil
}

// UpdateImageIntegration is pass-through to the underlying store.
func (ds *datastoreImpl) UpdateImageIntegration(ctx context.Context, integration *storage.ImageIntegration) error {
	if ok, err := integrationSAC.WriteAllowed(ctx); err != nil {
		return err
	} else if !ok {
		return sac.ErrResourceAccessDenied
	}

	return ds.storage.Upsert(ctx, integration)
}

// RemoveImageIntegration is pass-through to the underlying store.
func (ds *datastoreImpl) RemoveImageIntegration(ctx context.Context, id string) error {
	if ok, err := integrationSAC.WriteAllowed(ctx); err != nil {
		return err
	} else if !ok {
		return sac.ErrResourceAccessDenied
	}
	return ds.storage.Delete(ctx, id)
}

// SearchImageIntegrations
func (ds *datastoreImpl) SearchImageIntegrations(ctx context.Context, q *v1.Query) ([]*v1.SearchResult, error) {
	return ds.formattedSearcher.SearchImageIntegrations(ctx, q)
}
