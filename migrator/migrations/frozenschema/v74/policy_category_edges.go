package schema

import (
	"fmt"
	"reflect"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	v73Schema "github.com/stackrox/rox/migrator/migrations/frozenschema/v73"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/search"
)

var (
	// CreateTablePolicyCategoryEdgesStmt holds the create statement for table `policy_category_edges`.
	CreateTablePolicyCategoryEdgesStmt = &postgres.CreateStmts{
		GormModel: (*PolicyCategoryEdges)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// PolicyCategoryEdgesSchema is the go schema for table `policy_category_edges`.
	PolicyCategoryEdgesSchema = func() *walker.Schema {
		schema := walker.Walk(reflect.TypeOf((*storage.PolicyCategoryEdge)(nil)), "policy_category_edges")
		referencedSchemas := map[string]*walker.Schema{
			"storage.Policy":         v73Schema.PoliciesSchema,
			"storage.PolicyCategory": v73Schema.PolicyCategoriesSchema,
		}

		schema.ResolveReferences(func(messageTypeName string) *walker.Schema {
			return referencedSchemas[fmt.Sprintf("storage.%s", messageTypeName)]
		})
		schema.SetOptionsMap(search.Walk(v1.SearchCategory_POLICY_CATEGORY_EDGE, "policycategoryedge", (*storage.PolicyCategoryEdge)(nil)))
		schema.SetSearchScope([]v1.SearchCategory{
			v1.SearchCategory_POLICY_CATEGORY_EDGE,
			v1.SearchCategory_POLICY_CATEGORIES,
		}...)
		return schema
	}()
)

const (
	// PolicyCategoryEdgesTableName is the name of the table in which items are stored
	PolicyCategoryEdgesTableName = "policy_category_edges"
)

// PolicyCategoryEdges holds the Gorm model for Postgres table `policy_category_edges`.
type PolicyCategoryEdges struct {
	ID                  string                     `gorm:"column:id;type:varchar;primaryKey"`
	PolicyID            string                     `gorm:"column:policyid;type:varchar"`
	CategoryID          string                     `gorm:"column:categoryid;type:varchar"`
	Serialized          []byte                     `gorm:"column:serialized;type:bytea"`
	PoliciesRef         v73Schema.Policies         `gorm:"foreignKey:policyid;references:id;belongsTo;constraint:OnDelete:CASCADE"`
	PolicyCategoriesRef v73Schema.PolicyCategories `gorm:"foreignKey:categoryid;references:id;belongsTo;constraint:OnDelete:CASCADE"`
}
