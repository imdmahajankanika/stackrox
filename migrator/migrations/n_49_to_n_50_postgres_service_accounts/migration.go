// Code originally generated by pg-bindings generator.

package n49ton50

import (
	"context"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/migrator/migrations"
	frozenSchema "github.com/stackrox/rox/migrator/migrations/frozenschema/v73"
	"github.com/stackrox/rox/migrator/migrations/loghelper"
	legacy "github.com/stackrox/rox/migrator/migrations/n_49_to_n_50_postgres_service_accounts/legacy"
	pgStore "github.com/stackrox/rox/migrator/migrations/n_49_to_n_50_postgres_service_accounts/postgres"
	"github.com/stackrox/rox/migrator/types"
	pkgMigrations "github.com/stackrox/rox/pkg/migrations"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	"gorm.io/gorm"
)

var (
	startingSeqNum = pkgMigrations.BasePostgresDBVersionSeqNum() + 49 // 160

	migration = types.Migration{
		StartingSeqNum: startingSeqNum,
		VersionAfter:   &storage.Version{SeqNum: int32(startingSeqNum + 1)}, // 161
		Run: func(databases *types.Databases) error {
			legacyStore, err := legacy.New(databases.PkgRocksDB)
			if err != nil {
				return err
			}
			if err := move(databases.DBCtx, databases.GormDB, databases.PostgresDB, legacyStore); err != nil {
				return errors.Wrap(err,
					"moving service_accounts from rocksdb to postgres")
			}
			return nil
		},
	}
	batchSize = 10000
	schema    = frozenSchema.ServiceAccountsSchema
	log       = loghelper.LogWrapper{}
)

func move(ctx context.Context, gormDB *gorm.DB, postgresDB postgres.DB, legacyStore legacy.Store) error {
	store := pgStore.New(postgresDB)
	pgutils.CreateTableFromModel(context.Background(), gormDB, frozenSchema.CreateTableServiceAccountsStmt)

	var serviceAccounts []*storage.ServiceAccount
	err := walk(ctx, legacyStore, func(obj *storage.ServiceAccount) error {
		serviceAccounts = append(serviceAccounts, obj)
		if len(serviceAccounts) == batchSize {
			if err := store.UpsertMany(ctx, serviceAccounts); err != nil {
				log.WriteToStderrf("failed to persist service_accounts to store %v", err)
				return err
			}
			serviceAccounts = serviceAccounts[:0]
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(serviceAccounts) > 0 {
		if err = store.UpsertMany(ctx, serviceAccounts); err != nil {
			log.WriteToStderrf("failed to persist service_accounts to store %v", err)
			return err
		}
	}
	return nil
}

func walk(ctx context.Context, s legacy.Store, fn func(obj *storage.ServiceAccount) error) error {
	return s.Walk(ctx, fn)
}

func init() {
	migrations.MustRegisterMigration(migration)
}
