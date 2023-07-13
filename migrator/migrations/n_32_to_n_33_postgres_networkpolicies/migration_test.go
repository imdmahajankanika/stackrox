// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package n32ton33

import (
	"context"
	"testing"

	"github.com/stackrox/rox/generated/storage"
	legacy "github.com/stackrox/rox/migrator/migrations/n_32_to_n_33_postgres_networkpolicies/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_32_to_n_33_postgres_networkpolicies/postgres"
	pghelper "github.com/stackrox/rox/migrator/migrations/postgreshelper"

	"github.com/stackrox/rox/pkg/bolthelper"
	"github.com/stackrox/rox/pkg/sac"

	"github.com/stackrox/rox/pkg/testutils"

	"github.com/stretchr/testify/suite"

	bolt "go.etcd.io/bbolt"
)

func TestMigration(t *testing.T) {
	suite.Run(t, new(postgresMigrationSuite))
}

type postgresMigrationSuite struct {
	suite.Suite
	ctx context.Context

	legacyDB   *bolt.DB
	postgresDB *pghelper.TestPostgres
}

var _ suite.TearDownTestSuite = (*postgresMigrationSuite)(nil)

func (s *postgresMigrationSuite) SetupTest() {
	var err error
	s.legacyDB, err = bolthelper.NewTemp(s.T().Name() + ".db")
	s.NoError(err)

	s.Require().NoError(err)

	s.ctx = sac.WithAllAccess(context.Background())
	s.postgresDB = pghelper.ForT(s.T(), true)
}

func (s *postgresMigrationSuite) TearDownTest() {
	testutils.TearDownDB(s.legacyDB)
	s.postgresDB.Teardown(s.T())
}

func (s *postgresMigrationSuite) TestNetworkPolicyMigration() {
	newStore := pgStore.New(s.postgresDB.DB)
	legacyStore := legacy.New(s.legacyDB)

	// Prepare data and write to legacy DB
	var networkPolicys []*storage.NetworkPolicy
	for i := 0; i < 200; i++ {
		networkPolicy := &storage.NetworkPolicy{}
		s.NoError(testutils.FullInit(networkPolicy, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		networkPolicys = append(networkPolicys, networkPolicy)
		s.NoError(legacyStore.Upsert(s.ctx, networkPolicy))
	}

	// Move
	s.NoError(move(s.postgresDB.GetGormDB(), s.postgresDB.DB, legacyStore))

	// Verify
	count, err := newStore.Count(s.ctx)
	s.NoError(err)
	s.Equal(len(networkPolicys), count)
	for _, networkPolicy := range networkPolicys {
		fetched, exists, err := newStore.Get(s.ctx, networkPolicy.GetId())
		s.NoError(err)
		s.True(exists)
		s.Equal(networkPolicy, fetched)
	}
}
