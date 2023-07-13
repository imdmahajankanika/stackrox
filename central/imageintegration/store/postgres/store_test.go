// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package postgres

import (
	"context"
	"testing"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/suite"
)

type ImageIntegrationsStoreSuite struct {
	suite.Suite
	store  Store
	testDB *pgtest.TestPostgres
}

func TestImageIntegrationsStore(t *testing.T) {
	suite.Run(t, new(ImageIntegrationsStoreSuite))
}

func (s *ImageIntegrationsStoreSuite) SetupSuite() {

	s.testDB = pgtest.ForT(s.T())
	s.store = New(s.testDB.DB)
}

func (s *ImageIntegrationsStoreSuite) SetupTest() {
	ctx := sac.WithAllAccess(context.Background())
	tag, err := s.testDB.Exec(ctx, "TRUNCATE image_integrations CASCADE")
	s.T().Log("image_integrations", tag)
	s.NoError(err)
}

func (s *ImageIntegrationsStoreSuite) TearDownSuite() {
	s.testDB.Teardown(s.T())
}

func (s *ImageIntegrationsStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	imageIntegration := &storage.ImageIntegration{}
	s.NoError(testutils.FullInit(imageIntegration, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundImageIntegration, exists, err := store.Get(ctx, imageIntegration.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundImageIntegration)

	withNoAccessCtx := sac.WithNoAccess(ctx)

	s.NoError(store.Upsert(ctx, imageIntegration))
	foundImageIntegration, exists, err = store.Get(ctx, imageIntegration.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(imageIntegration, foundImageIntegration)

	imageIntegrationCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, imageIntegrationCount)
	imageIntegrationCount, err = store.Count(withNoAccessCtx)
	s.NoError(err)
	s.Zero(imageIntegrationCount)

	imageIntegrationExists, err := store.Exists(ctx, imageIntegration.GetId())
	s.NoError(err)
	s.True(imageIntegrationExists)
	s.NoError(store.Upsert(ctx, imageIntegration))
	s.ErrorIs(store.Upsert(withNoAccessCtx, imageIntegration), sac.ErrResourceAccessDenied)

	foundImageIntegration, exists, err = store.Get(ctx, imageIntegration.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(imageIntegration, foundImageIntegration)

	s.NoError(store.Delete(ctx, imageIntegration.GetId()))
	foundImageIntegration, exists, err = store.Get(ctx, imageIntegration.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundImageIntegration)
	s.ErrorIs(store.Delete(withNoAccessCtx, imageIntegration.GetId()), sac.ErrResourceAccessDenied)

	var imageIntegrations []*storage.ImageIntegration
	var imageIntegrationIDs []string
	for i := 0; i < 200; i++ {
		imageIntegration := &storage.ImageIntegration{}
		s.NoError(testutils.FullInit(imageIntegration, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		imageIntegrations = append(imageIntegrations, imageIntegration)
		imageIntegrationIDs = append(imageIntegrationIDs, imageIntegration.GetId())
	}

	s.NoError(store.UpsertMany(ctx, imageIntegrations))
	allImageIntegration, err := store.GetAll(ctx)
	s.NoError(err)
	s.ElementsMatch(imageIntegrations, allImageIntegration)

	imageIntegrationCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, imageIntegrationCount)

	s.NoError(store.DeleteMany(ctx, imageIntegrationIDs))

	imageIntegrationCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(0, imageIntegrationCount)
}
