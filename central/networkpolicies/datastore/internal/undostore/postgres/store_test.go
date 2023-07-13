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

type NetworkpolicyapplicationundorecordsStoreSuite struct {
	suite.Suite
	store  Store
	testDB *pgtest.TestPostgres
}

func TestNetworkpolicyapplicationundorecordsStore(t *testing.T) {
	suite.Run(t, new(NetworkpolicyapplicationundorecordsStoreSuite))
}

func (s *NetworkpolicyapplicationundorecordsStoreSuite) SetupSuite() {

	s.testDB = pgtest.ForT(s.T())
	s.store = New(s.testDB.DB)
}

func (s *NetworkpolicyapplicationundorecordsStoreSuite) SetupTest() {
	ctx := sac.WithAllAccess(context.Background())
	tag, err := s.testDB.Exec(ctx, "TRUNCATE networkpolicyapplicationundorecords CASCADE")
	s.T().Log("networkpolicyapplicationundorecords", tag)
	s.NoError(err)
}

func (s *NetworkpolicyapplicationundorecordsStoreSuite) TearDownSuite() {
	s.testDB.Teardown(s.T())
}

func (s *NetworkpolicyapplicationundorecordsStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	networkPolicyApplicationUndoRecord := &storage.NetworkPolicyApplicationUndoRecord{}
	s.NoError(testutils.FullInit(networkPolicyApplicationUndoRecord, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundNetworkPolicyApplicationUndoRecord, exists, err := store.Get(ctx, networkPolicyApplicationUndoRecord.GetClusterId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundNetworkPolicyApplicationUndoRecord)

	withNoAccessCtx := sac.WithNoAccess(ctx)

	s.NoError(store.Upsert(ctx, networkPolicyApplicationUndoRecord))
	foundNetworkPolicyApplicationUndoRecord, exists, err = store.Get(ctx, networkPolicyApplicationUndoRecord.GetClusterId())
	s.NoError(err)
	s.True(exists)
	s.Equal(networkPolicyApplicationUndoRecord, foundNetworkPolicyApplicationUndoRecord)

	networkPolicyApplicationUndoRecordCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, networkPolicyApplicationUndoRecordCount)
	networkPolicyApplicationUndoRecordCount, err = store.Count(withNoAccessCtx)
	s.NoError(err)
	s.Zero(networkPolicyApplicationUndoRecordCount)

	networkPolicyApplicationUndoRecordExists, err := store.Exists(ctx, networkPolicyApplicationUndoRecord.GetClusterId())
	s.NoError(err)
	s.True(networkPolicyApplicationUndoRecordExists)
	s.NoError(store.Upsert(ctx, networkPolicyApplicationUndoRecord))
	s.ErrorIs(store.Upsert(withNoAccessCtx, networkPolicyApplicationUndoRecord), sac.ErrResourceAccessDenied)

	foundNetworkPolicyApplicationUndoRecord, exists, err = store.Get(ctx, networkPolicyApplicationUndoRecord.GetClusterId())
	s.NoError(err)
	s.True(exists)
	s.Equal(networkPolicyApplicationUndoRecord, foundNetworkPolicyApplicationUndoRecord)

	s.NoError(store.Delete(ctx, networkPolicyApplicationUndoRecord.GetClusterId()))
	foundNetworkPolicyApplicationUndoRecord, exists, err = store.Get(ctx, networkPolicyApplicationUndoRecord.GetClusterId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundNetworkPolicyApplicationUndoRecord)
	s.NoError(store.Delete(withNoAccessCtx, networkPolicyApplicationUndoRecord.GetClusterId()))

	var networkPolicyApplicationUndoRecords []*storage.NetworkPolicyApplicationUndoRecord
	var networkPolicyApplicationUndoRecordIDs []string
	for i := 0; i < 200; i++ {
		networkPolicyApplicationUndoRecord := &storage.NetworkPolicyApplicationUndoRecord{}
		s.NoError(testutils.FullInit(networkPolicyApplicationUndoRecord, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		networkPolicyApplicationUndoRecords = append(networkPolicyApplicationUndoRecords, networkPolicyApplicationUndoRecord)
		networkPolicyApplicationUndoRecordIDs = append(networkPolicyApplicationUndoRecordIDs, networkPolicyApplicationUndoRecord.GetClusterId())
	}

	s.NoError(store.UpsertMany(ctx, networkPolicyApplicationUndoRecords))

	networkPolicyApplicationUndoRecordCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, networkPolicyApplicationUndoRecordCount)

	s.NoError(store.DeleteMany(ctx, networkPolicyApplicationUndoRecordIDs))

	networkPolicyApplicationUndoRecordCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(0, networkPolicyApplicationUndoRecordCount)
}
