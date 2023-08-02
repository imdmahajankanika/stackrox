// Code originally generated by pg-bindings generator.

package n16ton17

import (
	"context"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/migrations"
	frozenSchema "github.com/stackrox/rox/migrator/migrations/frozenschema/v73"
	"github.com/stackrox/rox/migrator/migrations/loghelper"
	legacy "github.com/stackrox/rox/migrator/migrations/n_16_to_n_17_postgres_compliance_operator_scan_setting_bindings/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_16_to_n_17_postgres_compliance_operator_scan_setting_bindings/postgres"
	"github.com/stackrox/rox/migrator/types"
	pkgMigrations "github.com/stackrox/rox/pkg/migrations"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	"gorm.io/gorm"
)

var (
	startingSeqNum = pkgMigrations.BasePostgresDBVersionSeqNum() + 16 // 127

	migration = types.Migration{
		StartingSeqNum: startingSeqNum,
		VersionAfter:   &storage.Version{SeqNum: int32(startingSeqNum + 1)}, // 128
		Run: func(databases *types.Databases) error {
			legacyStore, err := legacy.New(databases.PkgRocksDB)
			if err != nil {
				return err
			}
			if err := move(databases.DBCtx, databases.GormDB, databases.PostgresDB, legacyStore); err != nil {
				return errors.Wrap(err,
					"moving compliance_operator_scan_setting_bindings from rocksdb to postgres")
			}
			return nil
		},
	}
	batchSize = 10000
	schema    = frozenSchema.ComplianceOperatorScanSettingBindingsSchema
	log       = loghelper.LogWrapper{}
)

func move(ctx context.Context, gormDB *gorm.DB, postgresDB postgres.DB, legacyStore legacy.Store) error {
	store := pgStore.New(postgresDB)
	pgutils.CreateTableFromModel(context.Background(), gormDB, frozenSchema.CreateTableComplianceOperatorScanSettingBindingsStmt)

	var complianceOperatorScanSettingBindings []*storage.ComplianceOperatorScanSettingBinding
	err := walk(ctx, legacyStore, func(obj *storage.ComplianceOperatorScanSettingBinding) error {
		complianceOperatorScanSettingBindings = append(complianceOperatorScanSettingBindings, obj)
		if len(complianceOperatorScanSettingBindings) == batchSize {
			if err := store.UpsertMany(ctx, complianceOperatorScanSettingBindings); err != nil {
				log.WriteToStderrf("failed to persist compliance_operator_scan_setting_bindings to store %v", err)
				return err
			}
			complianceOperatorScanSettingBindings = complianceOperatorScanSettingBindings[:0]
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(complianceOperatorScanSettingBindings) > 0 {
		if err = store.UpsertMany(ctx, complianceOperatorScanSettingBindings); err != nil {
			log.WriteToStderrf("failed to persist compliance_operator_scan_setting_bindings to store %v", err)
			return err
		}
	}
	return nil
}

func walk(ctx context.Context, s legacy.Store, fn func(obj *storage.ComplianceOperatorScanSettingBinding) error) error {
	return s.Walk(ctx, fn)
}

func init() {
	migrations.MustRegisterMigration(migration)
}
