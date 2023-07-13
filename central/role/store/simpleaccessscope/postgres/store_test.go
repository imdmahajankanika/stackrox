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

type SimpleAccessScopesStoreSuite struct {
	suite.Suite
	store  Store
	testDB *pgtest.TestPostgres
}

func TestSimpleAccessScopesStore(t *testing.T) {
	suite.Run(t, new(SimpleAccessScopesStoreSuite))
}

func (s *SimpleAccessScopesStoreSuite) SetupSuite() {

	s.testDB = pgtest.ForT(s.T())
	s.store = New(s.testDB.DB)
}

func (s *SimpleAccessScopesStoreSuite) SetupTest() {
	ctx := sac.WithAllAccess(context.Background())
	tag, err := s.testDB.Exec(ctx, "TRUNCATE simple_access_scopes CASCADE")
	s.T().Log("simple_access_scopes", tag)
	s.NoError(err)
}

func (s *SimpleAccessScopesStoreSuite) TearDownSuite() {
	s.testDB.Teardown(s.T())
}

func (s *SimpleAccessScopesStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	simpleAccessScope := &storage.SimpleAccessScope{}
	s.NoError(testutils.FullInit(simpleAccessScope, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundSimpleAccessScope, exists, err := store.Get(ctx, simpleAccessScope.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundSimpleAccessScope)

	withNoAccessCtx := sac.WithNoAccess(ctx)

	s.NoError(store.Upsert(ctx, simpleAccessScope))
	foundSimpleAccessScope, exists, err = store.Get(ctx, simpleAccessScope.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(simpleAccessScope, foundSimpleAccessScope)

	simpleAccessScopeCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, simpleAccessScopeCount)
	simpleAccessScopeCount, err = store.Count(withNoAccessCtx)
	s.NoError(err)
	s.Zero(simpleAccessScopeCount)

	simpleAccessScopeExists, err := store.Exists(ctx, simpleAccessScope.GetId())
	s.NoError(err)
	s.True(simpleAccessScopeExists)
	s.NoError(store.Upsert(ctx, simpleAccessScope))
	s.ErrorIs(store.Upsert(withNoAccessCtx, simpleAccessScope), sac.ErrResourceAccessDenied)

	foundSimpleAccessScope, exists, err = store.Get(ctx, simpleAccessScope.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(simpleAccessScope, foundSimpleAccessScope)

	s.NoError(store.Delete(ctx, simpleAccessScope.GetId()))
	foundSimpleAccessScope, exists, err = store.Get(ctx, simpleAccessScope.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundSimpleAccessScope)
	s.ErrorIs(store.Delete(withNoAccessCtx, simpleAccessScope.GetId()), sac.ErrResourceAccessDenied)

	var simpleAccessScopes []*storage.SimpleAccessScope
	var simpleAccessScopeIDs []string
	for i := 0; i < 200; i++ {
		simpleAccessScope := &storage.SimpleAccessScope{}
		s.NoError(testutils.FullInit(simpleAccessScope, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		simpleAccessScopes = append(simpleAccessScopes, simpleAccessScope)
		simpleAccessScopeIDs = append(simpleAccessScopeIDs, simpleAccessScope.GetId())
	}

	s.NoError(store.UpsertMany(ctx, simpleAccessScopes))

	simpleAccessScopeCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, simpleAccessScopeCount)

	s.NoError(store.DeleteMany(ctx, simpleAccessScopeIDs))

	simpleAccessScopeCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(0, simpleAccessScopeCount)
}
