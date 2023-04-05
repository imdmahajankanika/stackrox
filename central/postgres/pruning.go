package postgres

import (
	"context"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
)

const (
	pruneActiveComponentsStmt = `DELETE FROM active_components child WHERE NOT EXISTS
		(SELECT 1 from deployments parent WHERE child.deploymentid = parent.id)`

	pruneClusterHealthStatusesStmt = `DELETE FROM cluster_health_statuses child WHERE NOT EXISTS
		(SELECT 1 FROM clusters parent WHERE
		child.Id = parent.Id)`

	getAllOrphanedAlerts = `SELECT id from alerts WHERE lifecyclestage = 0 and state = 0 and time < NOW() - INTERVAL '30 MINUTES' and NOT EXISTS
		"(SELECT 1 FROM deployments WHERE alerts.deployment_id = deployments.Id)`
)

var (
	log = logging.LoggerForModule()
)

// PruneActiveComponents - prunes active components
// TODO (ROX-12710):  This will no longer be necessary when the foreign keys are added back
func PruneActiveComponents(ctx context.Context, pool *postgres.DB) {
	if _, err := pool.Exec(ctx, pruneActiveComponentsStmt); err != nil {
		log.Errorf("failed to prune active components: %v", err)
	}
}

// PruneClusterHealthStatuses - prunes cluster health statuses
// TODO (ROX-12711):  This will no longer be necessary when the foreign keys are added back
func PruneClusterHealthStatuses(ctx context.Context, pool *postgres.DB) {
	if _, err := pool.Exec(ctx, pruneClusterHealthStatusesStmt); err != nil {
		log.Errorf("failed to prune cluster health statuses: %v", err)
	}
}

func getOrphanedAlertIDs(ctx context.Context, pool *postgres.DB) ([]string, error) {
	var ids []string
	rows, err := pool.Query(ctx, getAllOrphanedAlerts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get orphaned alerts")
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, errors.Wrap(err, "getting ids from orphaned alerts query")
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// GetOrphanedAlertIDs returns the alert IDs for alerts that are orphaned so they can be resolved
func GetOrphanedAlertIDs(ctx context.Context, pool *postgres.DB) ([]string, error) {
	return pgutils.Retry2(func() ([]string, error) {
		return getOrphanedAlertIDs(ctx, pool)
	})
}
