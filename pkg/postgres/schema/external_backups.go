// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
)

var (
	// CreateTableExternalBackupsStmt holds the create statement for table `external_backups`.
	CreateTableExternalBackupsStmt = &postgres.CreateStmts{
		GormModel: (*ExternalBackups)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// ExternalBackupsSchema is the go schema for table `external_backups`.
	ExternalBackupsSchema = func() *walker.Schema {
		schema := GetSchemaForTable("external_backups")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.ExternalBackup)(nil)), "external_backups")
		RegisterTable(schema, CreateTableExternalBackupsStmt)
		return schema
	}()
)

const (
	// ExternalBackupsTableName specifies the name of the table in postgres.
	ExternalBackupsTableName = "external_backups"
)

// ExternalBackups holds the Gorm model for Postgres table `external_backups`.
type ExternalBackups struct {
	ID         string `gorm:"column:id;type:varchar;primaryKey"`
	Serialized []byte `gorm:"column:serialized;type:bytea"`
}
