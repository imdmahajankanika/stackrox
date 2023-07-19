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

type ProcessBaselineResultsStoreSuite struct {
	suite.Suite
	store  Store
	testDB *pgtest.TestPostgres
}

func TestProcessBaselineResultsStore(t *testing.T) {
	suite.Run(t, new(ProcessBaselineResultsStoreSuite))
}

func (s *ProcessBaselineResultsStoreSuite) SetupSuite() {

	s.testDB = pgtest.ForT(s.T())
	s.store = New(s.testDB.DB)
}

func (s *ProcessBaselineResultsStoreSuite) SetupTest() {
	ctx := sac.WithAllAccess(context.Background())
	tag, err := s.testDB.Exec(ctx, "TRUNCATE process_baseline_results CASCADE")
	s.T().Log("process_baseline_results", tag)
	s.NoError(err)
}

func (s *ProcessBaselineResultsStoreSuite) TearDownSuite() {
	s.testDB.Teardown(s.T())
}

func (s *ProcessBaselineResultsStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	processBaselineResults := &storage.ProcessBaselineResults{}
	s.NoError(testutils.FullInit(processBaselineResults, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundProcessBaselineResults, exists, err := store.Get(ctx, processBaselineResults.GetDeploymentId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundProcessBaselineResults)

	withNoAccessCtx := sac.WithNoAccess(ctx)

	s.NoError(store.Upsert(ctx, processBaselineResults))
	foundProcessBaselineResults, exists, err = store.Get(ctx, processBaselineResults.GetDeploymentId())
	s.NoError(err)
	s.True(exists)
	s.Equal(processBaselineResults, foundProcessBaselineResults)

	processBaselineResultsCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, processBaselineResultsCount)
	processBaselineResultsCount, err = store.Count(withNoAccessCtx)
	s.NoError(err)
	s.Zero(processBaselineResultsCount)

	processBaselineResultsExists, err := store.Exists(ctx, processBaselineResults.GetDeploymentId())
	s.NoError(err)
	s.True(processBaselineResultsExists)
	s.NoError(store.Upsert(ctx, processBaselineResults))
	s.ErrorIs(store.Upsert(withNoAccessCtx, processBaselineResults), sac.ErrResourceAccessDenied)

	foundProcessBaselineResults, exists, err = store.Get(ctx, processBaselineResults.GetDeploymentId())
	s.NoError(err)
	s.True(exists)
	s.Equal(processBaselineResults, foundProcessBaselineResults)

	s.NoError(store.Delete(ctx, processBaselineResults.GetDeploymentId()))
	foundProcessBaselineResults, exists, err = store.Get(ctx, processBaselineResults.GetDeploymentId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundProcessBaselineResults)
	s.NoError(store.Delete(withNoAccessCtx, processBaselineResults.GetDeploymentId()))

	var processBaselineResultss []*storage.ProcessBaselineResults
	var processBaselineResultsIDs []string
	for i := 0; i < 200; i++ {
		processBaselineResults := &storage.ProcessBaselineResults{}
		s.NoError(testutils.FullInit(processBaselineResults, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		processBaselineResultss = append(processBaselineResultss, processBaselineResults)
		processBaselineResultsIDs = append(processBaselineResultsIDs, processBaselineResults.GetDeploymentId())
	}

	s.NoError(store.UpsertMany(ctx, processBaselineResultss))

	processBaselineResultsCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, processBaselineResultsCount)

	s.NoError(store.DeleteMany(ctx, processBaselineResultsIDs))

	processBaselineResultsCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(0, processBaselineResultsCount)
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
	expectedObjects        []*storage.ProcessBaselineResults
	expectedWriteError     error
}

func (s *ProcessBaselineResultsStoreSuite) getTestData(access storage.Access) (*storage.ProcessBaselineResults, *storage.ProcessBaselineResults, map[string]testCase) {
	objA := &storage.ProcessBaselineResults{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	objB := &storage.ProcessBaselineResults{}
	s.NoError(testutils.FullInit(objB, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	testCases := map[string]testCase{
		withAllAccess: {
			context:                sac.WithAllAccess(context.Background()),
			expectedObjIDs:         []string{objA.GetDeploymentId(), objB.GetDeploymentId()},
			expectedIdentifiers:    []string{objA.GetDeploymentId(), objB.GetDeploymentId()},
			expectedMissingIndices: []int{},
			expectedObjects:        []*storage.ProcessBaselineResults{objA, objB},
			expectedWriteError:     nil,
		},
		withNoAccess: {
			context:                sac.WithNoAccess(context.Background()),
			expectedObjIDs:         []string{},
			expectedIdentifiers:    []string{},
			expectedMissingIndices: []int{0, 1},
			expectedObjects:        []*storage.ProcessBaselineResults{},
			expectedWriteError:     sac.ErrResourceAccessDenied,
		},
		withNoAccessToCluster: {
			context: sac.WithGlobalAccessScopeChecker(context.Background(),
				sac.AllowFixedClusterLevelScopes(
					sac.AccessModeScopeKeyList(access),
					sac.ResourceScopeKeyList(targetResource),
					sac.ClusterScopeKeyList(uuid.Nil.String()),
				)),
			expectedObjIDs:         []string{},
			expectedIdentifiers:    []string{},
			expectedMissingIndices: []int{0, 1},
			expectedObjects:        []*storage.ProcessBaselineResults{},
			expectedWriteError:     sac.ErrResourceAccessDenied,
		},
		withAccess: {
			context: sac.WithGlobalAccessScopeChecker(context.Background(),
				sac.AllowFixedNamespaceLevelScopes(
					sac.AccessModeScopeKeyList(access),
					sac.ResourceScopeKeyList(targetResource),
					sac.ClusterScopeKeyList(objA.GetClusterId()),
					sac.NamespaceScopeKeyList(objA.GetNamespace()),
				)),
			expectedObjIDs:         []string{objA.GetDeploymentId()},
			expectedIdentifiers:    []string{objA.GetDeploymentId()},
			expectedMissingIndices: []int{1},
			expectedObjects:        []*storage.ProcessBaselineResults{objA},
			expectedWriteError:     nil,
		},
		withAccessToCluster: {
			context: sac.WithGlobalAccessScopeChecker(context.Background(),
				sac.AllowFixedClusterLevelScopes(
					sac.AccessModeScopeKeyList(access),
					sac.ResourceScopeKeyList(targetResource),
					sac.ClusterScopeKeyList(objA.GetClusterId()),
				)),
			expectedObjIDs:         []string{objA.GetDeploymentId()},
			expectedIdentifiers:    []string{objA.GetDeploymentId()},
			expectedMissingIndices: []int{1},
			expectedObjects:        []*storage.ProcessBaselineResults{objA},
			expectedWriteError:     nil,
		},
		withAccessToDifferentCluster: {
			context: sac.WithGlobalAccessScopeChecker(context.Background(),
				sac.AllowFixedClusterLevelScopes(
					sac.AccessModeScopeKeyList(access),
					sac.ResourceScopeKeyList(targetResource),
					sac.ClusterScopeKeyList("caaaaaaa-bbbb-4011-0000-111111111111"),
				)),
			expectedObjIDs:         []string{},
			expectedIdentifiers:    []string{},
			expectedMissingIndices: []int{0, 1},
			expectedObjects:        []*storage.ProcessBaselineResults{},
			expectedWriteError:     sac.ErrResourceAccessDenied,
		},
		withAccessToDifferentNs: {
			context: sac.WithGlobalAccessScopeChecker(context.Background(),
				sac.AllowFixedNamespaceLevelScopes(
					sac.AccessModeScopeKeyList(access),
					sac.ResourceScopeKeyList(targetResource),
					sac.ClusterScopeKeyList(objA.GetClusterId()),
					sac.NamespaceScopeKeyList("unknown ns"),
				)),
			expectedObjIDs:         []string{},
			expectedIdentifiers:    []string{},
			expectedMissingIndices: []int{0, 1},
			expectedObjects:        []*storage.ProcessBaselineResults{},
			expectedWriteError:     sac.ErrResourceAccessDenied,
		},
	}

	return objA, objB, testCases
}

func (s *ProcessBaselineResultsStoreSuite) TestSACUpsert() {
	obj, _, testCases := s.getTestData(storage.Access_READ_WRITE_ACCESS)
	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			assert.ErrorIs(t, s.store.Upsert(testCase.context, obj), testCase.expectedWriteError)
		})
	}
}

func (s *ProcessBaselineResultsStoreSuite) TestSACUpsertMany() {
	obj, _, testCases := s.getTestData(storage.Access_READ_WRITE_ACCESS)
	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			assert.ErrorIs(t, s.store.UpsertMany(testCase.context, []*storage.ProcessBaselineResults{obj}), testCase.expectedWriteError)
		})
	}
}

func (s *ProcessBaselineResultsStoreSuite) TestSACCount() {
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

func (s *ProcessBaselineResultsStoreSuite) TestSACWalk() {
	objA, objB, testCases := s.getTestData(storage.Access_READ_ACCESS)
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objA))
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objB))

	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			identifiers := []string{}
			getIDs := func(obj *storage.ProcessBaselineResults) error {
				identifiers = append(identifiers, obj.GetDeploymentId())
				return nil
			}
			err := s.store.Walk(testCase.context, getIDs)
			assert.NoError(t, err)
			assert.ElementsMatch(t, testCase.expectedIdentifiers, identifiers)
		})
	}
}

func (s *ProcessBaselineResultsStoreSuite) TestSACGetIDs() {
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

func (s *ProcessBaselineResultsStoreSuite) TestSACExists() {
	objA, _, testCases := s.getTestData(storage.Access_READ_ACCESS)
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objA))

	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			exists, err := s.store.Exists(testCase.context, objA.GetDeploymentId())
			assert.NoError(t, err)

			// Assumption from the test case structure: objA is always in the visible list
			// in the first position.
			expectedFound := len(testCase.expectedObjects) > 0
			assert.Equal(t, expectedFound, exists)
		})
	}
}

func (s *ProcessBaselineResultsStoreSuite) TestSACGet() {
	objA, _, testCases := s.getTestData(storage.Access_READ_ACCESS)
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objA))

	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			actual, exists, err := s.store.Get(testCase.context, objA.GetDeploymentId())
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

func (s *ProcessBaselineResultsStoreSuite) TestSACDelete() {
	objA, objB, testCases := s.getTestData(storage.Access_READ_WRITE_ACCESS)

	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			s.SetupTest()

			s.NoError(s.store.Upsert(withAllAccessCtx, objA))
			s.NoError(s.store.Upsert(withAllAccessCtx, objB))

			assert.NoError(t, s.store.Delete(testCase.context, objA.GetDeploymentId()))
			assert.NoError(t, s.store.Delete(testCase.context, objB.GetDeploymentId()))

			count, err := s.store.Count(withAllAccessCtx)
			assert.NoError(t, err)
			assert.Equal(t, 2-len(testCase.expectedObjects), count)

			// Ensure objects allowed by test scope were actually deleted
			for _, obj := range testCase.expectedObjects {
				found, err := s.store.Exists(withAllAccessCtx, obj.GetDeploymentId())
				assert.NoError(t, err)
				assert.False(t, found)
			}
		})
	}
}

func (s *ProcessBaselineResultsStoreSuite) TestSACDeleteMany() {
	objA, objB, testCases := s.getTestData(storage.Access_READ_WRITE_ACCESS)
	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			s.SetupTest()

			s.NoError(s.store.Upsert(withAllAccessCtx, objA))
			s.NoError(s.store.Upsert(withAllAccessCtx, objB))

			assert.NoError(t, s.store.DeleteMany(testCase.context, []string{
				objA.GetDeploymentId(),
				objB.GetDeploymentId(),
			}))

			count, err := s.store.Count(withAllAccessCtx)
			assert.NoError(t, err)
			assert.Equal(t, 2-len(testCase.expectedObjects), count)

			// Ensure objects allowed by test scope were actually deleted
			for _, obj := range testCase.expectedObjects {
				found, err := s.store.Exists(withAllAccessCtx, obj.GetDeploymentId())
				assert.NoError(t, err)
				assert.False(t, found)
			}
		})
	}
}

func (s *ProcessBaselineResultsStoreSuite) TestSACGetMany() {
	objA, objB, testCases := s.getTestData(storage.Access_READ_ACCESS)
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objA))
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objB))

	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			actual, missingIndices, err := s.store.GetMany(testCase.context, []string{objA.GetDeploymentId(), objB.GetDeploymentId()})
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
