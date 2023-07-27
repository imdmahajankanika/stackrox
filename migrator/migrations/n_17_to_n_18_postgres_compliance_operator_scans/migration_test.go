// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package n17ton18

import (
	"context"
	"testing"

	"github.com/stackrox/rox/generated/storage"
	legacy "github.com/stackrox/rox/migrator/migrations/n_17_to_n_18_postgres_compliance_operator_scans/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_17_to_n_18_postgres_compliance_operator_scans/postgres"
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
	s.postgresDB = pghelper.ForT(s.T(), false)
}

func (s *postgresMigrationSuite) TearDownTest() {
	rocksdbtest.TearDownRocksDB(s.legacyDB)
	s.postgresDB.Teardown(s.T())
}

func (s *postgresMigrationSuite) TestComplianceOperatorScanMigration() {
	newStore := pgStore.New(s.postgresDB.DB)
	legacyStore, err := legacy.New(s.legacyDB)
	s.NoError(err)

	// Prepare data and write to legacy DB
	var complianceOperatorScans []*storage.ComplianceOperatorScan
	for i := 0; i < 200; i++ {
		complianceOperatorScan := &storage.ComplianceOperatorScan{}
		s.NoError(testutils.FullInit(complianceOperatorScan, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		complianceOperatorScans = append(complianceOperatorScans, complianceOperatorScan)
	}

	s.NoError(legacyStore.UpsertMany(s.ctx, complianceOperatorScans))

	// Move
	s.NoError(move(s.postgresDB.GetGormDB(), s.postgresDB.DB, legacyStore))

	// Verify
	count, err := newStore.Count(s.ctx)
	s.NoError(err)
	s.Equal(len(complianceOperatorScans), count)
	for _, complianceOperatorScan := range complianceOperatorScans {
		fetched, exists, err := newStore.Get(s.ctx, complianceOperatorScan.GetId())
		s.NoError(err)
		s.True(exists)
		s.Equal(complianceOperatorScan, fetched)
	}
}
