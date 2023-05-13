// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package postgres

import (
	"context"
	"testing"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/suite"
)

type ContinuousIntegrationConfigsStoreSuite struct {
	suite.Suite
	store  Store
	testDB *pgtest.TestPostgres
}

func TestContinuousIntegrationConfigsStore(t *testing.T) {
	suite.Run(t, new(ContinuousIntegrationConfigsStoreSuite))
}

func (s *ContinuousIntegrationConfigsStoreSuite) SetupSuite() {
	s.T().Setenv(env.PostgresDatastoreEnabled.EnvVar(), "true")

	if !env.PostgresDatastoreEnabled.BooleanSetting() {
		s.T().Skip("Skip postgres store tests")
		s.T().SkipNow()
	}

	s.testDB = pgtest.ForT(s.T())
	s.store = New(s.testDB.DB)
}

func (s *ContinuousIntegrationConfigsStoreSuite) SetupTest() {
	ctx := sac.WithAllAccess(context.Background())
	tag, err := s.testDB.Exec(ctx, "TRUNCATE continuous_integration_configs CASCADE")
	s.T().Log("continuous_integration_configs", tag)
	s.NoError(err)
}

func (s *ContinuousIntegrationConfigsStoreSuite) TearDownSuite() {
	s.testDB.Teardown(s.T())
}

func (s *ContinuousIntegrationConfigsStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	continuousIntegrationConfig := &storage.ContinuousIntegrationConfig{}
	s.NoError(testutils.FullInit(continuousIntegrationConfig, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundContinuousIntegrationConfig, exists, err := store.Get(ctx, continuousIntegrationConfig.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundContinuousIntegrationConfig)

	withNoAccessCtx := sac.WithNoAccess(ctx)

	s.NoError(store.Upsert(ctx, continuousIntegrationConfig))
	foundContinuousIntegrationConfig, exists, err = store.Get(ctx, continuousIntegrationConfig.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(continuousIntegrationConfig, foundContinuousIntegrationConfig)

	continuousIntegrationConfigCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, continuousIntegrationConfigCount)
	continuousIntegrationConfigCount, err = store.Count(withNoAccessCtx)
	s.NoError(err)
	s.Zero(continuousIntegrationConfigCount)

	continuousIntegrationConfigExists, err := store.Exists(ctx, continuousIntegrationConfig.GetId())
	s.NoError(err)
	s.True(continuousIntegrationConfigExists)
	s.NoError(store.Upsert(ctx, continuousIntegrationConfig))
	s.ErrorIs(store.Upsert(withNoAccessCtx, continuousIntegrationConfig), sac.ErrResourceAccessDenied)

	foundContinuousIntegrationConfig, exists, err = store.Get(ctx, continuousIntegrationConfig.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(continuousIntegrationConfig, foundContinuousIntegrationConfig)

	s.NoError(store.Delete(ctx, continuousIntegrationConfig.GetId()))
	foundContinuousIntegrationConfig, exists, err = store.Get(ctx, continuousIntegrationConfig.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundContinuousIntegrationConfig)
	s.ErrorIs(store.Delete(withNoAccessCtx, continuousIntegrationConfig.GetId()), sac.ErrResourceAccessDenied)

	var continuousIntegrationConfigs []*storage.ContinuousIntegrationConfig
	var continuousIntegrationConfigIDs []string
	for i := 0; i < 200; i++ {
		continuousIntegrationConfig := &storage.ContinuousIntegrationConfig{}
		s.NoError(testutils.FullInit(continuousIntegrationConfig, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		continuousIntegrationConfigs = append(continuousIntegrationConfigs, continuousIntegrationConfig)
		continuousIntegrationConfigIDs = append(continuousIntegrationConfigIDs, continuousIntegrationConfig.GetId())
	}

	s.NoError(store.UpsertMany(ctx, continuousIntegrationConfigs))

	continuousIntegrationConfigCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, continuousIntegrationConfigCount)

	s.NoError(store.DeleteMany(ctx, continuousIntegrationConfigIDs))

	continuousIntegrationConfigCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(0, continuousIntegrationConfigCount)
}
