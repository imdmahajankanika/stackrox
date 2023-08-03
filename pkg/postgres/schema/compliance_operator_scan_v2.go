// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"fmt"
	"reflect"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/postgres/mapping"
)

var (
	// CreateTableComplianceOperatorScanV2Stmt holds the create statement for table `compliance_operator_scan_v2`.
	CreateTableComplianceOperatorScanV2Stmt = &postgres.CreateStmts{
		GormModel: (*ComplianceOperatorScanV2)(nil),
		Children: []*postgres.CreateStmts{
			&postgres.CreateStmts{
				GormModel: (*ComplianceOperatorScanV2Profiles)(nil),
				Children:  []*postgres.CreateStmts{},
			},
		},
	}

	// ComplianceOperatorScanV2Schema is the go schema for table `compliance_operator_scan_v2`.
	ComplianceOperatorScanV2Schema = func() *walker.Schema {
		schema := GetSchemaForTable("compliance_operator_scan_v2")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.ComplianceOperatorScanV2)(nil)), "compliance_operator_scan_v2")
		referencedSchemas := map[string]*walker.Schema{
			"storage.Cluster":                         ClustersSchema,
			"storage.ComplianceOperatorProfileV2":     ComplianceOperatorProfileV2Schema,
			"storage.ComplianceOperatorScanSettingV2": ComplianceOperatorScanSettingV2Schema,
		}

		schema.ResolveReferences(func(messageTypeName string) *walker.Schema {
			return referencedSchemas[fmt.Sprintf("storage.%s", messageTypeName)]
		})
		schema.SetOptionsMap(search.Walk(v1.SearchCategory_COMPLIANCE_SCAN, "complianceoperatorscanv2", (*storage.ComplianceOperatorScanV2)(nil)))
		RegisterTable(schema, CreateTableComplianceOperatorScanV2Stmt, features.ComplianceEnhancements.Enabled)
		mapping.RegisterCategoryToTable(v1.SearchCategory_COMPLIANCE_SCAN, schema)
		return schema
	}()
)

const (
	// ComplianceOperatorScanV2TableName specifies the name of the table in postgres.
	ComplianceOperatorScanV2TableName = "compliance_operator_scan_v2"
	// ComplianceOperatorScanV2ProfilesTableName specifies the name of the table in postgres.
	ComplianceOperatorScanV2ProfilesTableName = "compliance_operator_scan_v2_profiles"
)

// ComplianceOperatorScanV2 holds the Gorm model for Postgres table `compliance_operator_scan_v2`.
type ComplianceOperatorScanV2 struct {
	ID         string `gorm:"column:id;type:varchar;primaryKey"`
	ScanName   string `gorm:"column:scanname;type:varchar;uniqueIndex:scan_unique_indicator"`
	ClusterID  string `gorm:"column:clusterid;type:uuid;uniqueIndex:scan_unique_indicator;index:complianceoperatorscanv2_sac_filter,type:btree"`
	Serialized []byte `gorm:"column:serialized;type:bytea"`
}

// ComplianceOperatorScanV2Profiles holds the Gorm model for Postgres table `compliance_operator_scan_v2_profiles`.
type ComplianceOperatorScanV2Profiles struct {
	ComplianceOperatorScanV2ID  string                   `gorm:"column:compliance_operator_scan_v2_id;type:varchar;primaryKey"`
	Idx                         int                      `gorm:"column:idx;type:integer;primaryKey;index:complianceoperatorscanv2profiles_idx,type:btree"`
	ProfileID                   string                   `gorm:"column:profileid;type:varchar"`
	ComplianceOperatorScanV2Ref ComplianceOperatorScanV2 `gorm:"foreignKey:compliance_operator_scan_v2_id;references:id;belongsTo;constraint:OnDelete:CASCADE"`
}
