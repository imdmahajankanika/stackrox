// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
)

var (
	// CreateTableDeclarativeConfigHealthsStmt holds the create statement for table `declarative_config_healths`.
	CreateTableDeclarativeConfigHealthsStmt = &postgres.CreateStmts{
		GormModel: (*DeclarativeConfigHealths)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// DeclarativeConfigHealthsSchema is the go schema for table `declarative_config_healths`.
	DeclarativeConfigHealthsSchema = func() *walker.Schema {
		schema := GetSchemaForTable("declarative_config_healths")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.DeclarativeConfigHealth)(nil)), "declarative_config_healths")
		RegisterTable(schema, CreateTableDeclarativeConfigHealthsStmt)
		return schema
	}()
)

const (
	// DeclarativeConfigHealthsTableName specifies the name of the table in postgres.
	DeclarativeConfigHealthsTableName = "declarative_config_healths"
)

// DeclarativeConfigHealths holds the Gorm model for Postgres table `declarative_config_healths`.
type DeclarativeConfigHealths struct {
	ID         string `gorm:"column:id;type:uuid;primaryKey"`
	Serialized []byte `gorm:"column:serialized;type:bytea"`
}
