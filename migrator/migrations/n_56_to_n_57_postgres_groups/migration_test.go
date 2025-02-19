// Code originally generated by pg-bindings generator.
// Extended to cover additional edge cases in tests.

//go:build sql_integration

package n56ton57

import (
	"context"
	"testing"

	"github.com/stackrox/rox/generated/storage"
	legacy "github.com/stackrox/rox/migrator/migrations/n_56_to_n_57_postgres_groups/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_56_to_n_57_postgres_groups/postgres"
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
	s.postgresDB = pghelper.ForT(s.T(), false)
}

func (s *postgresMigrationSuite) TearDownTest() {
	testutils.TearDownDB(s.legacyDB)
	s.postgresDB.Teardown(s.T())
}

func (s *postgresMigrationSuite) TestGroupMigration() {
	newStore := pgStore.New(s.postgresDB.DB)
	legacyStore := legacy.New(s.legacyDB)

	// Prepare data and write to legacy DB
	var groups []*storage.Group
	for i := 0; i < 200; i++ {
		group := &storage.Group{}
		s.NoError(testutils.FullInit(group, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		groups = append(groups, group)
		s.NoError(legacyStore.Upsert(s.ctx, group))
	}

	// Special cases: as observed, there are sometimes groups still retained using the old format (i.e. stored by
	// composite key instead of UUID) as well as empty groups without an ID being set. These will lead to failures
	// when migrating. Below are some cases added, which should all be skipped during migration.

	// One completely empty group without ID and properties set stored in the old format.
	s.NoError(legacyStore.UpsertOldFormat(s.ctx, &storage.Group{}))
	// One completely initialized group stored in the old format.
	initializedGroupOldFormat := &storage.Group{}
	s.NoError(testutils.FullInit(initializedGroupOldFormat, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
	s.NoError(legacyStore.UpsertOldFormat(s.ctx, initializedGroupOldFormat))

	// Move
	s.NoError(move(s.ctx, s.postgresDB.GetGormDB(), s.postgresDB.DB, legacyStore))

	// Verify
	count, err := newStore.Count(s.ctx)
	s.NoError(err)
	s.Equal(len(groups), count)
	for _, group := range groups {
		fetched, exists, err := newStore.Get(s.ctx, group.GetProps().GetId())
		s.NoError(err)
		s.True(exists)
		s.Equal(group, fetched)
	}
}
