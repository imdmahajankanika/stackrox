// Code originally generated by pg-bindings generator.

package n13ton14

import (
	"context"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/migrations"
	frozenSchema "github.com/stackrox/rox/migrator/migrations/frozenschema/v73"
	"github.com/stackrox/rox/migrator/migrations/loghelper"
	legacy "github.com/stackrox/rox/migrator/migrations/n_13_to_n_14_postgres_compliance_operator_check_results/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_13_to_n_14_postgres_compliance_operator_check_results/postgres"
	"github.com/stackrox/rox/migrator/types"
	pkgMigrations "github.com/stackrox/rox/pkg/migrations"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	"gorm.io/gorm"
)

var (
	startingSeqNum = pkgMigrations.BasePostgresDBVersionSeqNum() + 13 // 124

	migration = types.Migration{
		StartingSeqNum: startingSeqNum,
		VersionAfter:   &storage.Version{SeqNum: int32(startingSeqNum + 1)}, // 125
		Run: func(databases *types.Databases) error {
			legacyStore, err := legacy.New(databases.PkgRocksDB)
			if err != nil {
				return err
			}
			if err := move(databases.DBCtx, databases.GormDB, databases.PostgresDB, legacyStore); err != nil {
				return errors.Wrap(err,
					"moving compliance_operator_check_results from rocksdb to postgres")
			}
			return nil
		},
	}
	batchSize = 10000
	schema    = frozenSchema.ComplianceOperatorCheckResultsSchema
	log       = loghelper.LogWrapper{}
)

func move(ctx context.Context, gormDB *gorm.DB, postgresDB postgres.DB, legacyStore legacy.Store) error {
	store := pgStore.New(postgresDB)
	pgutils.CreateTableFromModel(context.Background(), gormDB, frozenSchema.CreateTableComplianceOperatorCheckResultsStmt)

	var complianceOperatorCheckResults []*storage.ComplianceOperatorCheckResult
	err := walk(ctx, legacyStore, func(obj *storage.ComplianceOperatorCheckResult) error {
		complianceOperatorCheckResults = append(complianceOperatorCheckResults, obj)
		if len(complianceOperatorCheckResults) == batchSize {
			if err := store.UpsertMany(ctx, complianceOperatorCheckResults); err != nil {
				log.WriteToStderrf("failed to persist compliance_operator_check_results to store %v", err)
				return err
			}
			complianceOperatorCheckResults = complianceOperatorCheckResults[:0]
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(complianceOperatorCheckResults) > 0 {
		if err = store.UpsertMany(ctx, complianceOperatorCheckResults); err != nil {
			log.WriteToStderrf("failed to persist compliance_operator_check_results to store %v", err)
			return err
		}
	}
	return nil
}

func walk(ctx context.Context, s legacy.Store, fn func(obj *storage.ComplianceOperatorCheckResult) error) error {
	return s.Walk(ctx, fn)
}

func init() {
	migrations.MustRegisterMigration(migration)
}
