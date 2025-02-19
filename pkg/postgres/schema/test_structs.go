// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"
	"time"

	"github.com/lib/pq"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/postgres/mapping"
)

var (
	// CreateTableTestStructsStmt holds the create statement for table `test_structs`.
	CreateTableTestStructsStmt = &postgres.CreateStmts{
		GormModel: (*TestStructs)(nil),
		Children: []*postgres.CreateStmts{
			&postgres.CreateStmts{
				GormModel: (*TestStructsNesteds)(nil),
				Children:  []*postgres.CreateStmts{},
			},
		},
	}

	// TestStructsSchema is the go schema for table `test_structs`.
	TestStructsSchema = func() *walker.Schema {
		schema := GetSchemaForTable("test_structs")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.TestStruct)(nil)), "test_structs")
		schema.SetOptionsMap(search.Walk(v1.SearchCategory(101), "teststruct", (*storage.TestStruct)(nil)))
		RegisterTable(schema, CreateTableTestStructsStmt)
		mapping.RegisterCategoryToTable(v1.SearchCategory(101), schema)
		return schema
	}()
)

const (
	// TestStructsTableName specifies the name of the table in postgres.
	TestStructsTableName = "test_structs"
	// TestStructsNestedsTableName specifies the name of the table in postgres.
	TestStructsNestedsTableName = "test_structs_nesteds"
)

// TestStructs holds the Gorm model for Postgres table `test_structs`.
type TestStructs struct {
	Key1              string                  `gorm:"column:key1;type:varchar;primaryKey"`
	Key2              string                  `gorm:"column:key2;type:varchar"`
	StringSlice       *pq.StringArray         `gorm:"column:stringslice;type:text[]"`
	Bool              bool                    `gorm:"column:bool;type:bool"`
	Uint64            uint64                  `gorm:"column:uint64;type:bigint"`
	Int64             int64                   `gorm:"column:int64;type:bigint"`
	Float             float32                 `gorm:"column:float;type:numeric"`
	Labels            map[string]string       `gorm:"column:labels;type:jsonb"`
	Timestamp         *time.Time              `gorm:"column:timestamp;type:timestamp"`
	Enum              storage.TestStruct_Enum `gorm:"column:enum;type:integer"`
	Enums             *pq.Int32Array          `gorm:"column:enums;type:int[]"`
	String            string                  `gorm:"column:string_;type:varchar"`
	Int32Slice        *pq.Int32Array          `gorm:"column:int32slice;type:int[]"`
	OneofnestedNested string                  `gorm:"column:oneofnested_nested;type:varchar"`
	Serialized        []byte                  `gorm:"column:serialized;type:bytea"`
}

// TestStructsNesteds holds the Gorm model for Postgres table `test_structs_nesteds`.
type TestStructsNesteds struct {
	TestStructsKey1 string      `gorm:"column:test_structs_key1;type:varchar;primaryKey"`
	Idx             int         `gorm:"column:idx;type:integer;primaryKey;index:teststructsnesteds_idx,type:btree"`
	Nested          string      `gorm:"column:nested;type:varchar"`
	IsNested        bool        `gorm:"column:isnested;type:bool"`
	Int64           int64       `gorm:"column:int64;type:bigint"`
	Nested2Nested2  string      `gorm:"column:nested2_nested2;type:varchar"`
	Nested2IsNested bool        `gorm:"column:nested2_isnested;type:bool"`
	Nested2Int64    int64       `gorm:"column:nested2_int64;type:bigint"`
	TestStructsRef  TestStructs `gorm:"foreignKey:test_structs_key1;references:key1;belongsTo;constraint:OnDelete:CASCADE"`
}
