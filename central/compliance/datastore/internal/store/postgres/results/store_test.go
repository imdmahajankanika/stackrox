// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ComplianceRunResultsStoreSuite struct {
	suite.Suite
	store  Store
	testDB *pgtest.TestPostgres
}

func TestComplianceRunResultsStore(t *testing.T) {
	suite.Run(t, new(ComplianceRunResultsStoreSuite))
}

func (s *ComplianceRunResultsStoreSuite) SetupSuite() {

	s.testDB = pgtest.ForT(s.T())
	s.store = New(s.testDB.DB)
}

func (s *ComplianceRunResultsStoreSuite) SetupTest() {
	ctx := sac.WithAllAccess(context.Background())
	tag, err := s.testDB.Exec(ctx, "TRUNCATE compliance_run_results CASCADE")
	s.T().Log("compliance_run_results", tag)
	s.NoError(err)
}

func (s *ComplianceRunResultsStoreSuite) TearDownSuite() {
	s.testDB.Teardown(s.T())
}

func (s *ComplianceRunResultsStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	complianceRunResults := &storage.ComplianceRunResults{}
	s.NoError(testutils.FullInit(complianceRunResults, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundComplianceRunResults, exists, err := store.Get(ctx, complianceRunResults.GetRunMetadata().GetRunId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundComplianceRunResults)

	withNoAccessCtx := sac.WithNoAccess(ctx)

	s.NoError(store.Upsert(ctx, complianceRunResults))
	foundComplianceRunResults, exists, err = store.Get(ctx, complianceRunResults.GetRunMetadata().GetRunId())
	s.NoError(err)
	s.True(exists)
	s.Equal(complianceRunResults, foundComplianceRunResults)

	complianceRunResultsCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, complianceRunResultsCount)
	complianceRunResultsCount, err = store.Count(withNoAccessCtx)
	s.NoError(err)
	s.Zero(complianceRunResultsCount)

	complianceRunResultsExists, err := store.Exists(ctx, complianceRunResults.GetRunMetadata().GetRunId())
	s.NoError(err)
	s.True(complianceRunResultsExists)
	s.NoError(store.Upsert(ctx, complianceRunResults))
	s.ErrorIs(store.Upsert(withNoAccessCtx, complianceRunResults), sac.ErrResourceAccessDenied)

	foundComplianceRunResults, exists, err = store.Get(ctx, complianceRunResults.GetRunMetadata().GetRunId())
	s.NoError(err)
	s.True(exists)
	s.Equal(complianceRunResults, foundComplianceRunResults)

	s.NoError(store.Delete(ctx, complianceRunResults.GetRunMetadata().GetRunId()))
	foundComplianceRunResults, exists, err = store.Get(ctx, complianceRunResults.GetRunMetadata().GetRunId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundComplianceRunResults)
	s.NoError(store.Delete(withNoAccessCtx, complianceRunResults.GetRunMetadata().GetRunId()))

	var complianceRunResultss []*storage.ComplianceRunResults
	var complianceRunResultsIDs []string
	for i := 0; i < 200; i++ {
		complianceRunResults := &storage.ComplianceRunResults{}
		s.NoError(testutils.FullInit(complianceRunResults, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		complianceRunResultss = append(complianceRunResultss, complianceRunResults)
		complianceRunResultsIDs = append(complianceRunResultsIDs, complianceRunResults.GetRunMetadata().GetRunId())
	}

	s.NoError(store.UpsertMany(ctx, complianceRunResultss))

	complianceRunResultsCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, complianceRunResultsCount)

	s.NoError(store.DeleteMany(ctx, complianceRunResultsIDs))

	complianceRunResultsCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(0, complianceRunResultsCount)
}

const (
	withAllAccess                = "AllAccess"
	withNoAccess                 = "NoAccess"
	withAccess                   = "Access"
	withAccessToCluster          = "AccessToCluster"
	withNoAccessToCluster        = "NoAccessToCluster"
	withAccessToDifferentCluster = "AccessToDifferentCluster"
	withAccessToDifferentNs      = "AccessToDifferentNs"
)

var (
	withAllAccessCtx = sac.WithAllAccess(context.Background())
)

type testCase struct {
	context                context.Context
	expectedObjIDs         []string
	expectedIdentifiers    []string
	expectedMissingIndices []int
	expectedObjects        []*storage.ComplianceRunResults
	expectedWriteError     error
}

func (s *ComplianceRunResultsStoreSuite) getTestData(access storage.Access) (*storage.ComplianceRunResults, *storage.ComplianceRunResults, map[string]testCase) {
	objA := &storage.ComplianceRunResults{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	objB := &storage.ComplianceRunResults{}
	s.NoError(testutils.FullInit(objB, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	testCases := map[string]testCase{
		withAllAccess: {
			context:                sac.WithAllAccess(context.Background()),
			expectedObjIDs:         []string{objA.GetRunMetadata().GetRunId(), objB.GetRunMetadata().GetRunId()},
			expectedIdentifiers:    []string{objA.GetRunMetadata().GetRunId(), objB.GetRunMetadata().GetRunId()},
			expectedMissingIndices: []int{},
			expectedObjects:        []*storage.ComplianceRunResults{objA, objB},
			expectedWriteError:     nil,
		},
		withNoAccess: {
			context:                sac.WithNoAccess(context.Background()),
			expectedObjIDs:         []string{},
			expectedIdentifiers:    []string{},
			expectedMissingIndices: []int{0, 1},
			expectedObjects:        []*storage.ComplianceRunResults{},
			expectedWriteError:     sac.ErrResourceAccessDenied,
		},
		withNoAccessToCluster: {
			context: sac.WithGlobalAccessScopeChecker(context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(access),
					sac.ResourceScopeKeys(targetResource),
					sac.ClusterScopeKeys(uuid.Nil.String()),
				)),
			expectedObjIDs:         []string{},
			expectedIdentifiers:    []string{},
			expectedMissingIndices: []int{0, 1},
			expectedObjects:        []*storage.ComplianceRunResults{},
			expectedWriteError:     sac.ErrResourceAccessDenied,
		},
		withAccess: {
			context: sac.WithGlobalAccessScopeChecker(context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(access),
					sac.ResourceScopeKeys(targetResource),
					sac.ClusterScopeKeys(objA.GetRunMetadata().GetClusterId()),
				)),
			expectedObjIDs:         []string{objA.GetRunMetadata().GetRunId()},
			expectedIdentifiers:    []string{objA.GetRunMetadata().GetRunId()},
			expectedMissingIndices: []int{1},
			expectedObjects:        []*storage.ComplianceRunResults{objA},
			expectedWriteError:     nil,
		},
		withAccessToCluster: {
			context: sac.WithGlobalAccessScopeChecker(context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(access),
					sac.ResourceScopeKeys(targetResource),
					sac.ClusterScopeKeys(objA.GetRunMetadata().GetClusterId()),
				)),
			expectedObjIDs:         []string{objA.GetRunMetadata().GetRunId()},
			expectedIdentifiers:    []string{objA.GetRunMetadata().GetRunId()},
			expectedMissingIndices: []int{1},
			expectedObjects:        []*storage.ComplianceRunResults{objA},
			expectedWriteError:     nil,
		},
		withAccessToDifferentCluster: {
			context: sac.WithGlobalAccessScopeChecker(context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(access),
					sac.ResourceScopeKeys(targetResource),
					sac.ClusterScopeKeys("caaaaaaa-bbbb-4011-0000-111111111111"),
				)),
			expectedObjIDs:         []string{},
			expectedIdentifiers:    []string{},
			expectedMissingIndices: []int{0, 1},
			expectedObjects:        []*storage.ComplianceRunResults{},
			expectedWriteError:     sac.ErrResourceAccessDenied,
		},
		withAccessToDifferentNs: {
			context: sac.WithGlobalAccessScopeChecker(context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(access),
					sac.ResourceScopeKeys(targetResource),
					sac.ClusterScopeKeys(objA.GetRunMetadata().GetClusterId()),
					sac.NamespaceScopeKeys("unknown ns"),
				)),
			expectedObjIDs:         []string{objA.GetRunMetadata().GetRunId()},
			expectedIdentifiers:    []string{objA.GetRunMetadata().GetRunId()},
			expectedMissingIndices: []int{1},
			expectedObjects:        []*storage.ComplianceRunResults{objA},
			expectedWriteError:     nil,
		},
	}

	return objA, objB, testCases
}

func (s *ComplianceRunResultsStoreSuite) TestSACUpsert() {
	obj, _, testCases := s.getTestData(storage.Access_READ_WRITE_ACCESS)
	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			assert.ErrorIs(t, s.store.Upsert(testCase.context, obj), testCase.expectedWriteError)
		})
	}
}

func (s *ComplianceRunResultsStoreSuite) TestSACUpsertMany() {
	obj, _, testCases := s.getTestData(storage.Access_READ_WRITE_ACCESS)
	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			assert.ErrorIs(t, s.store.UpsertMany(testCase.context, []*storage.ComplianceRunResults{obj}), testCase.expectedWriteError)
		})
	}
}

func (s *ComplianceRunResultsStoreSuite) TestSACCount() {
	objA, objB, testCases := s.getTestData(storage.Access_READ_ACCESS)
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objA))
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objB))

	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			expectedCount := len(testCase.expectedObjects)
			count, err := s.store.Count(testCase.context)
			assert.NoError(t, err)
			assert.Equal(t, expectedCount, count)
		})
	}
}

func (s *ComplianceRunResultsStoreSuite) TestSACWalk() {
	objA, objB, testCases := s.getTestData(storage.Access_READ_ACCESS)
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objA))
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objB))

	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			identifiers := []string{}
			getIDs := func(obj *storage.ComplianceRunResults) error {
				identifiers = append(identifiers, obj.GetRunMetadata().GetRunId())
				return nil
			}
			err := s.store.Walk(testCase.context, getIDs)
			assert.NoError(t, err)
			assert.ElementsMatch(t, testCase.expectedIdentifiers, identifiers)
		})
	}
}

func (s *ComplianceRunResultsStoreSuite) TestSACGetIDs() {
	objA, objB, testCases := s.getTestData(storage.Access_READ_ACCESS)
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objA))
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objB))

	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			identifiers, err := s.store.GetIDs(testCase.context)
			assert.NoError(t, err)
			assert.EqualValues(t, testCase.expectedObjIDs, identifiers)
		})
	}
}

func (s *ComplianceRunResultsStoreSuite) TestSACExists() {
	objA, _, testCases := s.getTestData(storage.Access_READ_ACCESS)
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objA))

	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			exists, err := s.store.Exists(testCase.context, objA.GetRunMetadata().GetRunId())
			assert.NoError(t, err)

			// Assumption from the test case structure: objA is always in the visible list
			// in the first position.
			expectedFound := len(testCase.expectedObjects) > 0
			assert.Equal(t, expectedFound, exists)
		})
	}
}

func (s *ComplianceRunResultsStoreSuite) TestSACGet() {
	objA, _, testCases := s.getTestData(storage.Access_READ_ACCESS)
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objA))

	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			actual, exists, err := s.store.Get(testCase.context, objA.GetRunMetadata().GetRunId())
			assert.NoError(t, err)

			// Assumption from the test case structure: objA is always in the visible list
			// in the first position.
			expectedFound := len(testCase.expectedObjects) > 0
			assert.Equal(t, expectedFound, exists)
			if expectedFound {
				assert.Equal(t, objA, actual)
			} else {
				assert.Nil(t, actual)
			}
		})
	}
}

func (s *ComplianceRunResultsStoreSuite) TestSACDelete() {
	objA, objB, testCases := s.getTestData(storage.Access_READ_WRITE_ACCESS)

	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			s.SetupTest()

			s.NoError(s.store.Upsert(withAllAccessCtx, objA))
			s.NoError(s.store.Upsert(withAllAccessCtx, objB))

			assert.NoError(t, s.store.Delete(testCase.context, objA.GetRunMetadata().GetRunId()))
			assert.NoError(t, s.store.Delete(testCase.context, objB.GetRunMetadata().GetRunId()))

			count, err := s.store.Count(withAllAccessCtx)
			assert.NoError(t, err)
			assert.Equal(t, 2-len(testCase.expectedObjects), count)

			// Ensure objects allowed by test scope were actually deleted
			for _, obj := range testCase.expectedObjects {
				found, err := s.store.Exists(withAllAccessCtx, obj.GetRunMetadata().GetRunId())
				assert.NoError(t, err)
				assert.False(t, found)
			}
		})
	}
}

func (s *ComplianceRunResultsStoreSuite) TestSACDeleteMany() {
	objA, objB, testCases := s.getTestData(storage.Access_READ_WRITE_ACCESS)
	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			s.SetupTest()

			s.NoError(s.store.Upsert(withAllAccessCtx, objA))
			s.NoError(s.store.Upsert(withAllAccessCtx, objB))

			assert.NoError(t, s.store.DeleteMany(testCase.context, []string{
				objA.GetRunMetadata().GetRunId(),
				objB.GetRunMetadata().GetRunId(),
			}))

			count, err := s.store.Count(withAllAccessCtx)
			assert.NoError(t, err)
			assert.Equal(t, 2-len(testCase.expectedObjects), count)

			// Ensure objects allowed by test scope were actually deleted
			for _, obj := range testCase.expectedObjects {
				found, err := s.store.Exists(withAllAccessCtx, obj.GetRunMetadata().GetRunId())
				assert.NoError(t, err)
				assert.False(t, found)
			}
		})
	}
}

func (s *ComplianceRunResultsStoreSuite) TestSACGetMany() {
	objA, objB, testCases := s.getTestData(storage.Access_READ_ACCESS)
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objA))
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objB))

	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			actual, missingIndices, err := s.store.GetMany(testCase.context, []string{objA.GetRunMetadata().GetRunId(), objB.GetRunMetadata().GetRunId()})
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedObjects, actual)
			assert.Equal(t, testCase.expectedMissingIndices, missingIndices)
		})
	}

	s.T().Run("with no identifiers", func(t *testing.T) {
		actual, missingIndices, err := s.store.GetMany(withAllAccessCtx, []string{})
		assert.Nil(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, missingIndices)
	})
}
<<<<<<< HEAD
=======

const (
	withAllAccess                = "AllAccess"
	withNoAccess                 = "NoAccess"
	withAccessToDifferentNs      = "AccessToDifferentNs"
	withAccessToDifferentCluster = "AccessToDifferentCluster"
	withAccess                   = "Access"
	withAccessToCluster          = "AccessToCluster"
	withNoAccessToCluster        = "NoAccessToCluster"
)

func getSACContexts(obj *storage.ComplianceRunResults, access storage.Access) map[string]context.Context {
	return map[string]context.Context{
		withAllAccess: sac.WithAllAccess(context.Background()),
		withNoAccess:  sac.WithNoAccess(context.Background()),
		withAccessToDifferentCluster: sac.WithGlobalAccessScopeChecker(context.Background(),
			sac.AllowFixedClusterLevelScopes(
				sac.AccessModeScopeKeys(access),
				sac.ResourceScopeKeys(targetResource),
				sac.ClusterScopeKeys("caaaaaaa-bbbb-4011-0000-111111111111"),
			)),
		withAccessToDifferentNs: sac.WithGlobalAccessScopeChecker(context.Background(),
			sac.AllowFixedNamespaceLevelScopes(
				sac.AccessModeScopeKeys(access),
				sac.ResourceScopeKeys(targetResource),
				sac.ClusterScopeKeys(obj.GetRunMetadata().GetClusterId()),
				sac.NamespaceScopeKeys("unknown ns"),
			)),
		withAccess: sac.WithGlobalAccessScopeChecker(context.Background(),
			sac.AllowFixedClusterLevelScopes(
				sac.AccessModeScopeKeys(access),
				sac.ResourceScopeKeys(targetResource),
				sac.ClusterScopeKeys(obj.GetRunMetadata().GetClusterId()),
			)),
		withAccessToCluster: sac.WithGlobalAccessScopeChecker(context.Background(),
			sac.AllowFixedClusterLevelScopes(
				sac.AccessModeScopeKeys(access),
				sac.ResourceScopeKeys(targetResource),
				sac.ClusterScopeKeys(obj.GetRunMetadata().GetClusterId()),
			)),
		withNoAccessToCluster: sac.WithGlobalAccessScopeChecker(context.Background(),
			sac.AllowFixedClusterLevelScopes(
				sac.AccessModeScopeKeys(access),
				sac.ResourceScopeKeys(targetResource),
				sac.ClusterScopeKeys(uuid.Nil.String()),
			)),
	}
}
>>>>>>> d6486707d8 (Re-adjust function renames)
