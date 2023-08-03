// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package postgres

import (
	"context"
	"testing"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stretchr/testify/suite"
)

type ComplianceOperatorCheckResultV2StoreSuite struct {
	suite.Suite
	store  Store
	testDB *pgtest.TestPostgres
}

func TestComplianceOperatorCheckResultV2Store(t *testing.T) {
	suite.Run(t, new(ComplianceOperatorCheckResultV2StoreSuite))
}

func (s *ComplianceOperatorCheckResultV2StoreSuite) SetupSuite() {

	s.T().Setenv(features.ComplianceEnhancements.EnvVar(), "true")
	if !features.ComplianceEnhancements.Enabled() {
		s.T().Skip("Skip postgres store tests because feature flag is off")
		s.T().SkipNow()
	}

	s.testDB = pgtest.ForT(s.T())
	s.store = New(s.testDB.DB)
}

func (s *ComplianceOperatorCheckResultV2StoreSuite) SetupTest() {
	ctx := sac.WithAllAccess(context.Background())
	tag, err := s.testDB.Exec(ctx, "TRUNCATE compliance_operator_check_result_v2 CASCADE")
	s.T().Log("compliance_operator_check_result_v2", tag)
	s.NoError(err)
}

func (s *ComplianceOperatorCheckResultV2StoreSuite) TearDownSuite() {
	s.testDB.Teardown(s.T())
}

func (s *ComplianceOperatorCheckResultV2StoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	complianceOperatorCheckResultV2 := &storage.ComplianceOperatorCheckResultV2{}
	s.NoError(testutils.FullInit(complianceOperatorCheckResultV2, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundComplianceOperatorCheckResultV2, exists, err := store.Get(ctx, complianceOperatorCheckResultV2.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundComplianceOperatorCheckResultV2)

	withNoAccessCtx := sac.WithNoAccess(ctx)

	s.NoError(store.Upsert(ctx, complianceOperatorCheckResultV2))
	foundComplianceOperatorCheckResultV2, exists, err = store.Get(ctx, complianceOperatorCheckResultV2.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(complianceOperatorCheckResultV2, foundComplianceOperatorCheckResultV2)

	complianceOperatorCheckResultV2Count, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, complianceOperatorCheckResultV2Count)
	complianceOperatorCheckResultV2Count, err = store.Count(withNoAccessCtx)
	s.NoError(err)
	s.Zero(complianceOperatorCheckResultV2Count)

	complianceOperatorCheckResultV2Exists, err := store.Exists(ctx, complianceOperatorCheckResultV2.GetId())
	s.NoError(err)
	s.True(complianceOperatorCheckResultV2Exists)
	s.NoError(store.Upsert(ctx, complianceOperatorCheckResultV2))
	s.ErrorIs(store.Upsert(withNoAccessCtx, complianceOperatorCheckResultV2), sac.ErrResourceAccessDenied)

	foundComplianceOperatorCheckResultV2, exists, err = store.Get(ctx, complianceOperatorCheckResultV2.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(complianceOperatorCheckResultV2, foundComplianceOperatorCheckResultV2)

	s.NoError(store.Delete(ctx, complianceOperatorCheckResultV2.GetId()))
	foundComplianceOperatorCheckResultV2, exists, err = store.Get(ctx, complianceOperatorCheckResultV2.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundComplianceOperatorCheckResultV2)
	s.ErrorIs(store.Delete(withNoAccessCtx, complianceOperatorCheckResultV2.GetId()), sac.ErrResourceAccessDenied)

	var complianceOperatorCheckResultV2s []*storage.ComplianceOperatorCheckResultV2
	var complianceOperatorCheckResultV2IDs []string
	for i := 0; i < 200; i++ {
		complianceOperatorCheckResultV2 := &storage.ComplianceOperatorCheckResultV2{}
		s.NoError(testutils.FullInit(complianceOperatorCheckResultV2, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		complianceOperatorCheckResultV2s = append(complianceOperatorCheckResultV2s, complianceOperatorCheckResultV2)
		complianceOperatorCheckResultV2IDs = append(complianceOperatorCheckResultV2IDs, complianceOperatorCheckResultV2.GetId())
	}

	s.NoError(store.UpsertMany(ctx, complianceOperatorCheckResultV2s))
	allComplianceOperatorCheckResultV2, err := store.GetAll(ctx)
	s.NoError(err)
	s.ElementsMatch(complianceOperatorCheckResultV2s, allComplianceOperatorCheckResultV2)

	complianceOperatorCheckResultV2Count, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, complianceOperatorCheckResultV2Count)

	s.NoError(store.DeleteMany(ctx, complianceOperatorCheckResultV2IDs))

	complianceOperatorCheckResultV2Count, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(0, complianceOperatorCheckResultV2Count)
}
