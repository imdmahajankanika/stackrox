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

type SecretsStoreSuite struct {
	suite.Suite
	store  Store
	testDB *pgtest.TestPostgres
}

func TestSecretsStore(t *testing.T) {
	suite.Run(t, new(SecretsStoreSuite))
}

func (s *SecretsStoreSuite) SetupSuite() {

	s.testDB = pgtest.ForT(s.T())
	s.store = New(s.testDB.DB)
}

func (s *SecretsStoreSuite) SetupTest() {
	ctx := sac.WithAllAccess(context.Background())
	tag, err := s.testDB.Exec(ctx, "TRUNCATE secrets CASCADE")
	s.T().Log("secrets", tag)
	s.NoError(err)
}

func (s *SecretsStoreSuite) TearDownSuite() {
	s.testDB.Teardown(s.T())
}

func (s *SecretsStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	secret := &storage.Secret{}
	s.NoError(testutils.FullInit(secret, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundSecret, exists, err := store.Get(ctx, secret.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundSecret)

	withNoAccessCtx := sac.WithNoAccess(ctx)

	s.NoError(store.Upsert(ctx, secret))
	foundSecret, exists, err = store.Get(ctx, secret.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(secret, foundSecret)

	secretCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, secretCount)
	secretCount, err = store.Count(withNoAccessCtx)
	s.NoError(err)
	s.Zero(secretCount)

	secretExists, err := store.Exists(ctx, secret.GetId())
	s.NoError(err)
	s.True(secretExists)
	s.NoError(store.Upsert(ctx, secret))
	s.ErrorIs(store.Upsert(withNoAccessCtx, secret), sac.ErrResourceAccessDenied)

	foundSecret, exists, err = store.Get(ctx, secret.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(secret, foundSecret)

	s.NoError(store.Delete(ctx, secret.GetId()))
	foundSecret, exists, err = store.Get(ctx, secret.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundSecret)
	s.NoError(store.Delete(withNoAccessCtx, secret.GetId()))

	var secrets []*storage.Secret
	var secretIDs []string
	for i := 0; i < 200; i++ {
		secret := &storage.Secret{}
		s.NoError(testutils.FullInit(secret, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		secrets = append(secrets, secret)
		secretIDs = append(secretIDs, secret.GetId())
	}

	s.NoError(store.UpsertMany(ctx, secrets))

	secretCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, secretCount)

	s.NoError(store.DeleteMany(ctx, secretIDs))

	secretCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(0, secretCount)
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
	expectedObjects        []*storage.Secret
	expectedWriteError     error
}

func (s *SecretsStoreSuite) getTestData(access storage.Access) (*storage.Secret, *storage.Secret, map[string]testCase) {
	objA := &storage.Secret{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	objB := &storage.Secret{}
	s.NoError(testutils.FullInit(objB, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	testCases := map[string]testCase{
		withAllAccess: {
			context:                sac.WithAllAccess(context.Background()),
			expectedObjIDs:         []string{objA.GetId(), objB.GetId()},
			expectedIdentifiers:    []string{objA.GetId(), objB.GetId()},
			expectedMissingIndices: []int{},
			expectedObjects:        []*storage.Secret{objA, objB},
			expectedWriteError:     nil,
		},
		withNoAccess: {
			context:                sac.WithNoAccess(context.Background()),
			expectedObjIDs:         []string{},
			expectedIdentifiers:    []string{},
			expectedMissingIndices: []int{0, 1},
			expectedObjects:        []*storage.Secret{},
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
			expectedObjects:        []*storage.Secret{},
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
			expectedObjIDs:         []string{objA.GetId()},
			expectedIdentifiers:    []string{objA.GetId()},
			expectedMissingIndices: []int{1},
			expectedObjects:        []*storage.Secret{objA},
			expectedWriteError:     nil,
		},
		withAccessToCluster: {
			context: sac.WithGlobalAccessScopeChecker(context.Background(),
				sac.AllowFixedClusterLevelScopes(
					sac.AccessModeScopeKeyList(access),
					sac.ResourceScopeKeyList(targetResource),
					sac.ClusterScopeKeyList(objA.GetClusterId()),
				)),
			expectedObjIDs:         []string{objA.GetId()},
			expectedIdentifiers:    []string{objA.GetId()},
			expectedMissingIndices: []int{1},
			expectedObjects:        []*storage.Secret{objA},
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
			expectedObjects:        []*storage.Secret{},
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
			expectedObjects:        []*storage.Secret{},
			expectedWriteError:     sac.ErrResourceAccessDenied,
		},
	}

	return objA, objB, testCases
}

func (s *SecretsStoreSuite) TestSACUpsert() {
	obj, _, testCases := s.getTestData(storage.Access_READ_WRITE_ACCESS)
	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			assert.ErrorIs(t, s.store.Upsert(testCase.context, obj), testCase.expectedWriteError)
		})
	}
}

func (s *SecretsStoreSuite) TestSACUpsertMany() {
	obj, _, testCases := s.getTestData(storage.Access_READ_WRITE_ACCESS)
	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			assert.ErrorIs(t, s.store.UpsertMany(testCase.context, []*storage.Secret{obj}), testCase.expectedWriteError)
		})
	}
}

func (s *SecretsStoreSuite) TestSACCount() {
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

func (s *SecretsStoreSuite) TestSACWalk() {
	objA, objB, testCases := s.getTestData(storage.Access_READ_ACCESS)
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objA))
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objB))

	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			identifiers := []string{}
			getIDs := func(obj *storage.Secret) error {
				identifiers = append(identifiers, obj.GetId())
				return nil
			}
			err := s.store.Walk(testCase.context, getIDs)
			assert.NoError(t, err)
			assert.ElementsMatch(t, testCase.expectedIdentifiers, identifiers)
		})
	}
}

func (s *SecretsStoreSuite) TestSACGetIDs() {
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

func (s *SecretsStoreSuite) TestSACExists() {
	objA, _, testCases := s.getTestData(storage.Access_READ_ACCESS)
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objA))

	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			exists, err := s.store.Exists(testCase.context, objA.GetId())
			assert.NoError(t, err)

			// Assumption from the test case structure: objA is always in the visible list
			// in the first position.
			expectedFound := len(testCase.expectedObjects) > 0
			assert.Equal(t, expectedFound, exists)
		})
	}
}

func (s *SecretsStoreSuite) TestSACGet() {
	objA, _, testCases := s.getTestData(storage.Access_READ_ACCESS)
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objA))

	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			actual, exists, err := s.store.Get(testCase.context, objA.GetId())
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

func (s *SecretsStoreSuite) TestSACDelete() {
	objA, objB, testCases := s.getTestData(storage.Access_READ_WRITE_ACCESS)

	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			s.SetupTest()

			s.NoError(s.store.Upsert(withAllAccessCtx, objA))
			s.NoError(s.store.Upsert(withAllAccessCtx, objB))

			assert.NoError(t, s.store.Delete(testCase.context, objA.GetId()))
			assert.NoError(t, s.store.Delete(testCase.context, objB.GetId()))

			count, err := s.store.Count(withAllAccessCtx)
			assert.NoError(t, err)
			assert.Equal(t, 2-len(testCase.expectedObjects), count)

			// Ensure objects allowed by test scope were actually deleted
			for _, obj := range testCase.expectedObjects {
				found, err := s.store.Exists(withAllAccessCtx, obj.GetId())
				assert.NoError(t, err)
				assert.False(t, found)
			}
		})
	}
}

func (s *SecretsStoreSuite) TestSACDeleteMany() {
	objA, objB, testCases := s.getTestData(storage.Access_READ_WRITE_ACCESS)
	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			s.SetupTest()

			s.NoError(s.store.Upsert(withAllAccessCtx, objA))
			s.NoError(s.store.Upsert(withAllAccessCtx, objB))

			assert.NoError(t, s.store.DeleteMany(testCase.context, []string{
				objA.GetId(),
				objB.GetId(),
			}))

			count, err := s.store.Count(withAllAccessCtx)
			assert.NoError(t, err)
			assert.Equal(t, 2-len(testCase.expectedObjects), count)

			// Ensure objects allowed by test scope were actually deleted
			for _, obj := range testCase.expectedObjects {
				found, err := s.store.Exists(withAllAccessCtx, obj.GetId())
				assert.NoError(t, err)
				assert.False(t, found)
			}
		})
	}
}

func (s *SecretsStoreSuite) TestSACGetMany() {
	objA, objB, testCases := s.getTestData(storage.Access_READ_ACCESS)
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objA))
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objB))

	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			actual, missingIndices, err := s.store.GetMany(testCase.context, []string{objA.GetId(), objB.GetId()})
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
