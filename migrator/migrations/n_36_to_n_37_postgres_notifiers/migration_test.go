// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package n36ton37

import (
	"context"
	"testing"

	"github.com/stackrox/rox/generated/storage"
	legacy "github.com/stackrox/rox/migrator/migrations/n_36_to_n_37_postgres_notifiers/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_36_to_n_37_postgres_notifiers/postgres"
	pghelper "github.com/stackrox/rox/migrator/migrations/postgreshelper"

	"github.com/stackrox/rox/pkg/bolthelper"
	"github.com/stackrox/rox/pkg/sac"

	"github.com/stackrox/rox/pkg/env"

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
	s.T().Setenv(env.PostgresDatastoreEnabled.EnvVar(), "true")
	if !env.PostgresDatastoreEnabled.BooleanSetting() {
		s.T().Skip("Skip postgres store tests")
		s.T().SkipNow()
	}

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

func (s *postgresMigrationSuite) TestNotifierMigration() {
	newStore := pgStore.New(s.postgresDB.Pool)
	legacyStore := legacy.New(s.legacyDB)

	// Prepare data and write to legacy DB
	var notifiers []*storage.Notifier
	for i := 0; i < 200; i++ {
		notifier := &storage.Notifier{}
		s.NoError(testutils.FullInit(notifier, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		notifiers = append(notifiers, notifier)
		s.NoError(legacyStore.Upsert(s.ctx, notifier))
	}

	// Move
	s.NoError(move(s.postgresDB.GetGormDB(), s.postgresDB.Pool, legacyStore))

	// Verify
	count, err := newStore.Count(s.ctx)
	s.NoError(err)
	s.Equal(len(notifiers), count)
	for _, notifier := range notifiers {
		fetched, exists, err := newStore.Get(s.ctx, notifier.GetId())
		s.NoError(err)
		s.True(exists)
		s.Equal(notifier, fetched)
	}
}
