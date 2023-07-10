// Code generated by make bootstrap_migration generator. DO NOT EDIT.
package runner

import (
	// Postgres -> Postgres migrations
	_ "github.com/stackrox/rox/migrator/migrations/m_168_to_m_169_postgres_remove_clustercve_permission"
	_ "github.com/stackrox/rox/migrator/migrations/m_169_to_m_170_collections_sac_resource_migration"
	_ "github.com/stackrox/rox/migrator/migrations/m_170_to_m_171_create_policy_categories_and_edges"
	_ "github.com/stackrox/rox/migrator/migrations/m_171_to_m_172_move_scope_to_collection_in_report_configurations"
	_ "github.com/stackrox/rox/migrator/migrations/m_172_to_m_173_network_flows_partition"
	_ "github.com/stackrox/rox/migrator/migrations/m_173_to_m_174_group_unique_constraint"
	_ "github.com/stackrox/rox/migrator/migrations/m_174_to_m_175_enable_search_on_api_tokens"
	_ "github.com/stackrox/rox/migrator/migrations/m_175_to_m_176_create_notification_schedule_table"
	_ "github.com/stackrox/rox/migrator/migrations/m_176_to_m_177_network_baselines_cidr"
	_ "github.com/stackrox/rox/migrator/migrations/m_177_to_m_178_group_permissions"
	_ "github.com/stackrox/rox/migrator/migrations/m_178_to_m_179_embedded_collections_search_label"
	_ "github.com/stackrox/rox/migrator/migrations/m_179_to_m_180_openshift_policy_exclusions"
	_ "github.com/stackrox/rox/migrator/migrations/m_180_to_m_181_move_to_blobstore"
	_ "github.com/stackrox/rox/migrator/migrations/m_181_to_m_182_group_role_permission_with_access_one"
	_ "github.com/stackrox/rox/migrator/migrations/m_182_to_m_183_remove_default_scope_manager_role"
	_ "github.com/stackrox/rox/migrator/migrations/m_183_to_m_184_move_declarative_config_health"
	_ "github.com/stackrox/rox/migrator/migrations/m_184_to_m_185_remove_policy_vulnerability_report_resources"
	_ "github.com/stackrox/rox/migrator/migrations/m_185_to_m_186_more_policy_migrations"
)
