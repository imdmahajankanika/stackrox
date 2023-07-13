// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package n43ton44

import (
	"context"
	"testing"

	"github.com/stackrox/rox/generated/storage"
	legacy "github.com/stackrox/rox/migrator/migrations/n_43_to_n_44_postgres_report_configurations/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_43_to_n_44_postgres_report_configurations/postgres"
	pghelper "github.com/stackrox/rox/migrator/migrations/postgreshelper"

	"github.com/stackrox/rox/pkg/rocksdb"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/testutils/rocksdbtest"
	"github.com/stretchr/testify/suite"
)

func TestMigration(t *testing.T) {
	suite.Run(t, new(postgresMigrationSuite))
}

type postgresMigrationSuite struct {
	suite.Suite
	ctx context.Context

	legacyDB   *rocksdb.RocksDB
	postgresDB *pghelper.TestPostgres
}

var _ suite.TearDownTestSuite = (*postgresMigrationSuite)(nil)

func (s *postgresMigrationSuite) SetupTest() {
	var err error
	s.legacyDB, err = rocksdb.NewTemp(s.T().Name())
	s.NoError(err)

	s.Require().NoError(err)

	s.ctx = sac.WithAllAccess(context.Background())
	s.postgresDB = pghelper.ForT(s.T(), true)
}

func (s *postgresMigrationSuite) TearDownTest() {
	rocksdbtest.TearDownRocksDB(s.legacyDB)
	s.postgresDB.Teardown(s.T())
}

func (s *postgresMigrationSuite) TestReportConfigurationMigration() {
	newStore := pgStore.New(s.postgresDB.DB)
	legacyStore, err := legacy.New(s.legacyDB)
	s.NoError(err)

	// Prepare data and write to legacy DB
	var reportConfigurations []*storage.ReportConfiguration
	for i := 0; i < 200; i++ {
		reportConfiguration := &storage.ReportConfiguration{}
		s.NoError(testutils.FullInit(reportConfiguration, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		reportConfigurations = append(reportConfigurations, reportConfiguration)
	}

	s.NoError(legacyStore.UpsertMany(s.ctx, reportConfigurations))

	// Move
	s.NoError(move(s.postgresDB.GetGormDB(), s.postgresDB.DB, legacyStore))

	// Verify
	count, err := newStore.Count(s.ctx)
	s.NoError(err)
	s.Equal(len(reportConfigurations), count)
	for _, reportConfiguration := range reportConfigurations {
		fetched, exists, err := newStore.Get(s.ctx, reportConfiguration.GetId())
		s.NoError(err)
		s.True(exists)
		s.Equal(reportConfiguration, fetched)
	}
}
