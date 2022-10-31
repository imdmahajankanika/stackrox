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

type NetworkEntitiesStoreSuite struct {
	suite.Suite
	store  Store
	testDB *pgtest.TestPostgres
}

func TestNetworkEntitiesStore(t *testing.T) {
	suite.Run(t, new(NetworkEntitiesStoreSuite))
}

func (s *NetworkEntitiesStoreSuite) SetupSuite() {
	s.T().Setenv(env.PostgresDatastoreEnabled.EnvVar(), "true")

	if !env.PostgresDatastoreEnabled.BooleanSetting() {
		s.T().Skip("Skip postgres store tests")
		s.T().SkipNow()
	}

	s.testDB = pgtest.ForT(s.T())
	s.store = New(s.testDB.Pool)
}

func (s *NetworkEntitiesStoreSuite) SetupTest() {
	ctx := sac.WithAllAccess(context.Background())
	tag, err := s.testDB.Exec(ctx, "TRUNCATE network_entities CASCADE")
	s.T().Log("network_entities", tag)
	s.NoError(err)
}

func (s *NetworkEntitiesStoreSuite) TearDownSuite() {
	s.testDB.Teardown(s.T())
}

func (s *NetworkEntitiesStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	networkEntity := &storage.NetworkEntity{}
	s.NoError(testutils.FullInit(networkEntity, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundNetworkEntity, exists, err := store.Get(ctx, networkEntity.GetInfo().GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundNetworkEntity)

	withNoAccessCtx := sac.WithNoAccess(ctx)

	s.NoError(store.Upsert(ctx, networkEntity))
	foundNetworkEntity, exists, err = store.Get(ctx, networkEntity.GetInfo().GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(networkEntity, foundNetworkEntity)

	networkEntityCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, networkEntityCount)
	networkEntityCount, err = store.Count(withNoAccessCtx)
	s.NoError(err)
	s.Zero(networkEntityCount)

	networkEntityExists, err := store.Exists(ctx, networkEntity.GetInfo().GetId())
	s.NoError(err)
	s.True(networkEntityExists)
	s.NoError(store.Upsert(ctx, networkEntity))
	s.ErrorIs(store.Upsert(withNoAccessCtx, networkEntity), sac.ErrResourceAccessDenied)

	foundNetworkEntity, exists, err = store.Get(ctx, networkEntity.GetInfo().GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(networkEntity, foundNetworkEntity)

	s.NoError(store.Delete(ctx, networkEntity.GetInfo().GetId()))
	foundNetworkEntity, exists, err = store.Get(ctx, networkEntity.GetInfo().GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundNetworkEntity)
	s.ErrorIs(store.Delete(withNoAccessCtx, networkEntity.GetInfo().GetId()), sac.ErrResourceAccessDenied)

	var networkEntitys []*storage.NetworkEntity
	var networkEntityIDs []string
	for i := 0; i < 200; i++ {
		networkEntity := &storage.NetworkEntity{}
		s.NoError(testutils.FullInit(networkEntity, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		networkEntitys = append(networkEntitys, networkEntity)
		networkEntityIDs = append(networkEntityIDs, networkEntity.GetInfo().GetId())
	}

	s.NoError(store.UpsertMany(ctx, networkEntitys))

	networkEntityCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, networkEntityCount)

	s.NoError(store.DeleteMany(ctx, networkEntityIDs))

	networkEntityCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(0, networkEntityCount)
}
