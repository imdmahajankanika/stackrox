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

type K8sRolesStoreSuite struct {
	suite.Suite
	store  Store
	testDB *pgtest.TestPostgres
}

func TestK8sRolesStore(t *testing.T) {
	suite.Run(t, new(K8sRolesStoreSuite))
}

func (s *K8sRolesStoreSuite) SetupSuite() {

	s.testDB = pgtest.ForT(s.T())
	s.store = New(s.testDB.DB)
}

func (s *K8sRolesStoreSuite) SetupTest() {
	ctx := sac.WithAllAccess(context.Background())
	tag, err := s.testDB.Exec(ctx, "TRUNCATE k8s_roles CASCADE")
	s.T().Log("k8s_roles", tag)
	s.NoError(err)
}

func (s *K8sRolesStoreSuite) TearDownSuite() {
	s.testDB.Teardown(s.T())
}

func (s *K8sRolesStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	k8SRole := &storage.K8SRole{}
	s.NoError(testutils.FullInit(k8SRole, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundK8SRole, exists, err := store.Get(ctx, k8SRole.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundK8SRole)

	withNoAccessCtx := sac.WithNoAccess(ctx)

	s.NoError(store.Upsert(ctx, k8SRole))
	foundK8SRole, exists, err = store.Get(ctx, k8SRole.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(k8SRole, foundK8SRole)

	k8SRoleCount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, k8SRoleCount)
	k8SRoleCount, err = store.Count(withNoAccessCtx)
	s.NoError(err)
	s.Zero(k8SRoleCount)

	k8SRoleExists, err := store.Exists(ctx, k8SRole.GetId())
	s.NoError(err)
	s.True(k8SRoleExists)
	s.NoError(store.Upsert(ctx, k8SRole))
	s.ErrorIs(store.Upsert(withNoAccessCtx, k8SRole), sac.ErrResourceAccessDenied)

	foundK8SRole, exists, err = store.Get(ctx, k8SRole.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(k8SRole, foundK8SRole)

	s.NoError(store.Delete(ctx, k8SRole.GetId()))
	foundK8SRole, exists, err = store.Get(ctx, k8SRole.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundK8SRole)
	s.NoError(store.Delete(withNoAccessCtx, k8SRole.GetId()))

	var k8SRoles []*storage.K8SRole
	var k8SRoleIDs []string
	for i := 0; i < 200; i++ {
		k8SRole := &storage.K8SRole{}
		s.NoError(testutils.FullInit(k8SRole, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		k8SRoles = append(k8SRoles, k8SRole)
		k8SRoleIDs = append(k8SRoleIDs, k8SRole.GetId())
	}

	s.NoError(store.UpsertMany(ctx, k8SRoles))

	k8SRoleCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, k8SRoleCount)

	s.NoError(store.DeleteMany(ctx, k8SRoleIDs))

	k8SRoleCount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(0, k8SRoleCount)
}

const (
	withAllAccess           = "AllAccess"
	withNoAccess            = "NoAccess"
	withAccessToDifferentNs = "AccessToDifferentNs"
	withAccess              = "Access"
	withAccessToCluster     = "AccessToCluster"
	withNoAccessToCluster   = "NoAccessToCluster"
)

var (
	withAllAccessCtx = sac.WithAllAccess(context.Background())
)

type testCase struct {
	context                context.Context
	expectedIDs            []string
	expectedIdentifiers    []string
	expectedMissingIndices []int
	expectedObjects        []*storage.K8SRole
	expectedWriteError     error
}

func (s *K8sRolesStoreSuite) getTestData(access storage.Access) (*storage.K8SRole, *storage.K8SRole, map[string]testCase) {
	objA := &storage.K8SRole{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	objB := &storage.K8SRole{}
	s.NoError(testutils.FullInit(objB, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	testCases := map[string]testCase{
		withAllAccess: {
			context:                sac.WithAllAccess(context.Background()),
			expectedMissingIndices: []int{},
			expectedObjects:        []*storage.K8SRole{objA, objB},
			expectedWriteError:     nil,
		},
		withNoAccess: {
			context:                sac.WithNoAccess(context.Background()),
			expectedMissingIndices: []int{0, 1},
			expectedObjects:        []*storage.K8SRole{},
			expectedWriteError:     sac.ErrResourceAccessDenied,
		},
		withNoAccessToCluster: {
			context: sac.WithGlobalAccessScopeChecker(context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(access),
					sac.ResourceScopeKeys(targetResource),
					sac.ClusterScopeKeys(uuid.Nil.String()),
				),
			),
			expectedMissingIndices: []int{0, 1},
			expectedObjects:        []*storage.K8SRole{},
			expectedWriteError:     sac.ErrResourceAccessDenied,
		},
		withAccessToDifferentNs: {
			context: sac.WithGlobalAccessScopeChecker(context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(access),
					sac.ResourceScopeKeys(targetResource),
					sac.ClusterScopeKeys(objA.GetClusterId()),
					sac.NamespaceScopeKeys("unknown ns"),
				),
			),
			expectedMissingIndices: []int{0, 1},
			expectedObjects:        []*storage.K8SRole{},
			expectedWriteError:     sac.ErrResourceAccessDenied,
		},
		withAccess: {
			context: sac.WithGlobalAccessScopeChecker(context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(access),
					sac.ResourceScopeKeys(targetResource),
					sac.ClusterScopeKeys(objA.GetClusterId()),
					sac.NamespaceScopeKeys(objA.GetNamespace()),
				),
			),
			expectedMissingIndices: []int{1},
			expectedObjects:        []*storage.K8SRole{objA},
			expectedWriteError:     nil,
		},
		withAccessToCluster: {
			context: sac.WithGlobalAccessScopeChecker(context.Background(),
				sac.AllowFixedScopes(
					sac.AccessModeScopeKeys(access),
					sac.ResourceScopeKeys(targetResource),
					sac.ClusterScopeKeys(objA.GetClusterId()),
				),
			),
			expectedMissingIndices: []int{1},
			expectedObjects:        []*storage.K8SRole{objA},
			expectedWriteError:     nil,
		},
	}

	return objA, objB, testCases
}

func (s *K8sRolesStoreSuite) TestSACUpsert() {
	obj, _, testCases := s.getTestData(storage.Access_READ_WRITE_ACCESS)
	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			assert.ErrorIs(t, s.store.Upsert(testCase.context, obj), testCase.expectedWriteError)
		})
	}
}

func (s *K8sRolesStoreSuite) TestSACUpsertMany() {
	obj, _, testCases := s.getTestData(storage.Access_READ_WRITE_ACCESS)
	for name, testCase := range testCases {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			assert.ErrorIs(t, s.store.UpsertMany(testCase.context, []*storage.K8SRole{obj}), testCase.expectedWriteError)
		})
	}
}

func (s *K8sRolesStoreSuite) TestSACCount() {
	objA := &storage.K8SRole{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	objB := &storage.K8SRole{}
	s.NoError(testutils.FullInit(objB, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	withAllAccessCtx := sac.WithAllAccess(context.Background())
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objA))
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objB))

	ctxs := getSACContexts(objA, storage.Access_READ_ACCESS)
	for name, expectedCount := range map[string]int{
		withAllAccess:           2,
		withNoAccess:            0,
		withNoAccessToCluster:   0,
		withAccessToDifferentNs: 0,
		withAccess:              1,
		withAccessToCluster:     1,
	} {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			count, err := s.store.Count(ctxs[name])
			assert.NoError(t, err)
			assert.Equal(t, expectedCount, count)
		})
	}
}

func (s *K8sRolesStoreSuite) TestSACWalk() {
	objA := &storage.K8SRole{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	objB := &storage.K8SRole{}
	s.NoError(testutils.FullInit(objB, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	withAllAccessCtx := sac.WithAllAccess(context.Background())
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objA))
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objB))

	ctxs := getSACContexts(objA, storage.Access_READ_ACCESS)
	for name, expectedIDs := range map[string][]string{
		withAllAccess:           []string{objA.GetId(), objB.GetId()},
		withNoAccess:            []string{},
		withNoAccessToCluster:   []string{},
		withAccessToDifferentNs: []string{},
		withAccess:              []string{objA.GetId()},
		withAccessToCluster:     []string{objA.GetId()},
	} {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			identifiers := []string{}
			getIDs := func(obj *storage.K8SRole) error {
				identifiers = append(identifiers, obj.GetId())
				return nil
			}
			err := s.store.Walk(ctxs[name], getIDs)
			assert.NoError(t, err)
			assert.ElementsMatch(t, expectedIDs, identifiers)
		})
	}
}

func (s *K8sRolesStoreSuite) TestSACGetIDs() {
	objA := &storage.K8SRole{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	objB := &storage.K8SRole{}
	s.NoError(testutils.FullInit(objB, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	withAllAccessCtx := sac.WithAllAccess(context.Background())
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objA))
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objB))

	ctxs := getSACContexts(objA, storage.Access_READ_ACCESS)
	for name, expectedIDs := range map[string][]string{
		withAllAccess:           []string{objA.GetId(), objB.GetId()},
		withNoAccess:            []string{},
		withNoAccessToCluster:   []string{},
		withAccessToDifferentNs: []string{},
		withAccess:              []string{objA.GetId()},
		withAccessToCluster:     []string{objA.GetId()},
	} {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			identifiers, err := s.store.GetIDs(ctxs[name])
			assert.NoError(t, err)
			assert.EqualValues(t, expectedIDs, identifiers)
		})
	}
}

func (s *K8sRolesStoreSuite) TestSACExists() {
	objA := &storage.K8SRole{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	withAllAccessCtx := sac.WithAllAccess(context.Background())
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objA))

	ctxs := getSACContexts(objA, storage.Access_READ_ACCESS)
	for name, expected := range map[string]bool{
		withAllAccess:           true,
		withNoAccess:            false,
		withNoAccessToCluster:   false,
		withAccessToDifferentNs: false,
		withAccess:              true,
		withAccessToCluster:     true,
	} {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			exists, err := s.store.Exists(ctxs[name], objA.GetId())
			assert.NoError(t, err)
			assert.Equal(t, expected, exists)
		})
	}
}

func (s *K8sRolesStoreSuite) TestSACGet() {
	objA := &storage.K8SRole{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	withAllAccessCtx := sac.WithAllAccess(context.Background())
	s.Require().NoError(s.store.Upsert(withAllAccessCtx, objA))

	ctxs := getSACContexts(objA, storage.Access_READ_ACCESS)
	for name, expected := range map[string]bool{
		withAllAccess:           true,
		withNoAccess:            false,
		withNoAccessToCluster:   false,
		withAccessToDifferentNs: false,
		withAccess:              true,
		withAccessToCluster:     true,
	} {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			actual, exists, err := s.store.Get(ctxs[name], objA.GetId())
			assert.NoError(t, err)
			assert.Equal(t, expected, exists)
			if expected == true {
				assert.Equal(t, objA, actual)
			} else {
				assert.Nil(t, actual)
			}
		})
	}
}

func (s *K8sRolesStoreSuite) TestSACDelete() {
	objA := &storage.K8SRole{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	objB := &storage.K8SRole{}
	s.NoError(testutils.FullInit(objB, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
	withAllAccessCtx := sac.WithAllAccess(context.Background())

	ctxs := getSACContexts(objA, storage.Access_READ_WRITE_ACCESS)
	for name, expectedCount := range map[string]int{
		withAllAccess:           0,
		withNoAccess:            2,
		withNoAccessToCluster:   2,
		withAccessToDifferentNs: 2,
		withAccess:              1,
		withAccessToCluster:     1,
	} {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			s.SetupTest()

			s.NoError(s.store.Upsert(withAllAccessCtx, objA))
			s.NoError(s.store.Upsert(withAllAccessCtx, objB))

			assert.NoError(t, s.store.Delete(ctxs[name], objA.GetId()))
			assert.NoError(t, s.store.Delete(ctxs[name], objB.GetId()))

			count, err := s.store.Count(withAllAccessCtx)
			assert.NoError(t, err)
			assert.Equal(t, expectedCount, count)
		})
	}
}

func (s *K8sRolesStoreSuite) TestSACDeleteMany() {
	objA := &storage.K8SRole{}
	s.NoError(testutils.FullInit(objA, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))

	objB := &storage.K8SRole{}
	s.NoError(testutils.FullInit(objB, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
	withAllAccessCtx := sac.WithAllAccess(context.Background())

	ctxs := getSACContexts(objA, storage.Access_READ_WRITE_ACCESS)
	for name, expectedCount := range map[string]int{
		withAllAccess:           0,
		withNoAccess:            2,
		withNoAccessToCluster:   2,
		withAccessToDifferentNs: 2,
		withAccess:              1,
		withAccessToCluster:     1,
	} {
		s.T().Run(fmt.Sprintf("with %s", name), func(t *testing.T) {
			s.SetupTest()

			s.NoError(s.store.Upsert(withAllAccessCtx, objA))
			s.NoError(s.store.Upsert(withAllAccessCtx, objB))

			assert.NoError(t, s.store.DeleteMany(ctxs[name], []string{
				objA.GetId(),
				objB.GetId(),
			}))

			count, err := s.store.Count(withAllAccessCtx)
			assert.NoError(t, err)
			assert.Equal(t, expectedCount, count)
		})
	}
}

func (s *K8sRolesStoreSuite) TestSACGetMany() {
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

func getSACContexts(obj *storage.K8SRole, access storage.Access) map[string]context.Context {
	return map[string]context.Context{
		withAllAccess: sac.WithAllAccess(context.Background()),
		withNoAccess:  sac.WithNoAccess(context.Background()),
		withAccessToDifferentNs: sac.WithGlobalAccessScopeChecker(context.Background(),
			sac.AllowFixedScopes(
				sac.AccessModeScopeKeys(access),
				sac.ResourceScopeKeys(targetResource),
				sac.ClusterScopeKeys(obj.GetClusterId()),
				sac.NamespaceScopeKeys("unknown ns"),
			)),
		withAccess: sac.WithGlobalAccessScopeChecker(context.Background(),
			sac.AllowFixedScopes(
				sac.AccessModeScopeKeys(access),
				sac.ResourceScopeKeys(targetResource),
				sac.ClusterScopeKeys(obj.GetClusterId()),
				sac.NamespaceScopeKeys(obj.GetNamespace()),
			)),
		withAccessToCluster: sac.WithGlobalAccessScopeChecker(context.Background(),
			sac.AllowFixedScopes(
				sac.AccessModeScopeKeys(access),
				sac.ResourceScopeKeys(targetResource),
				sac.ClusterScopeKeys(obj.GetClusterId()),
			)),
		withNoAccessToCluster: sac.WithGlobalAccessScopeChecker(context.Background(),
			sac.AllowFixedScopes(
				sac.AccessModeScopeKeys(access),
				sac.ResourceScopeKeys(targetResource),
				sac.ClusterScopeKeys(uuid.Nil.String()),
			)),
	}
}
