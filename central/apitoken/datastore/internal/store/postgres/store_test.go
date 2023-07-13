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

type APITokensStoreSuite struct {
	suite.Suite
	store  Store
	testDB *pgtest.TestPostgres
}

func TestAPITokensStore(t *testing.T) {
	suite.Run(t, new(APITokensStoreSuite))
}

func (s *APITokensStoreSuite) SetupSuite() {

	s.testDB = pgtest.ForT(s.T())
	s.store = New(s.testDB.DB)
}

func (s *APITokensStoreSuite) SetupTest() {
	ctx := sac.WithAllAccess(context.Background())
	tag, err := s.testDB.Exec(ctx, "TRUNCATE api_tokens CASCADE")
	s.T().Log("api_tokens", tag)
	s.NoError(err)
}

func (s *APITokensStoreSuite) TearDownSuite() {
	s.testDB.Teardown(s.T())
}

func (s *APITokensStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	tokenMetadata := &storage.TokenMetadata{}
	s.NoError(testutils.FullInit(tokenMetadata, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundTokenMetadata, exists, err := store.Get(ctx, tokenMetadata.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundTokenMetadata)

	withNoAccessCtx := sac.WithNoAccess(ctx)

	s.NoError(store.Upsert(ctx, tokenMetadata))
	foundTokenMetadata, exists, err = store.Get(ctx, tokenMetadata.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(tokenMetadata, foundTokenMetadata)

	tokenMetadataCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, tokenMetadataCount)
	tokenMetadataCount, err = store.Count(withNoAccessCtx)
	s.NoError(err)
	s.Zero(tokenMetadataCount)

	tokenMetadataExists, err := store.Exists(ctx, tokenMetadata.GetId())
	s.NoError(err)
	s.True(tokenMetadataExists)
	s.NoError(store.Upsert(ctx, tokenMetadata))
	s.ErrorIs(store.Upsert(withNoAccessCtx, tokenMetadata), sac.ErrResourceAccessDenied)

	foundTokenMetadata, exists, err = store.Get(ctx, tokenMetadata.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(tokenMetadata, foundTokenMetadata)

	s.NoError(store.Delete(ctx, tokenMetadata.GetId()))
	foundTokenMetadata, exists, err = store.Get(ctx, tokenMetadata.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundTokenMetadata)
	s.ErrorIs(store.Delete(withNoAccessCtx, tokenMetadata.GetId()), sac.ErrResourceAccessDenied)

	var tokenMetadatas []*storage.TokenMetadata
	var tokenMetadataIDs []string
	for i := 0; i < 200; i++ {
		tokenMetadata := &storage.TokenMetadata{}
		s.NoError(testutils.FullInit(tokenMetadata, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		tokenMetadatas = append(tokenMetadatas, tokenMetadata)
		tokenMetadataIDs = append(tokenMetadataIDs, tokenMetadata.GetId())
	}

	s.NoError(store.UpsertMany(ctx, tokenMetadatas))

	tokenMetadataCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, tokenMetadataCount)

	s.NoError(store.DeleteMany(ctx, tokenMetadataIDs))

	tokenMetadataCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(0, tokenMetadataCount)
}
