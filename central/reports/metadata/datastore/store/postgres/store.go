// Code generated by pg-bindings generator. DO NOT EDIT.

package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/stackrox/rox/central/metrics"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/sac/resources"
	pgSearch "github.com/stackrox/rox/pkg/search/postgres"
	"github.com/stackrox/rox/pkg/sync"
	"gorm.io/gorm"
)

const (
	baseTable = "report_metadata"
	storeName = "ReportMetadata"

	// using copyFrom, we may not even want to batch.  It would probably be simpler
	// to deal with failures if we just sent it all.  Something to think about as we
	// proceed and move into more e2e and larger performance testing
	batchSize = 10000
)

var (
	log            = logging.LoggerForModule()
	schema         = pkgSchema.ReportMetadataSchema
	targetResource = resources.WorkflowAdministration
)

type storeType = storage.ReportMetadata

// Store is the interface to interact with the storage for storage.ReportMetadata
type Store interface {
	Upsert(ctx context.Context, obj *storeType) error
	UpsertMany(ctx context.Context, objs []*storeType) error
	Delete(ctx context.Context, reportID string) error
	DeleteByQuery(ctx context.Context, q *v1.Query) error
	DeleteMany(ctx context.Context, identifiers []string) error

	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, reportID string) (bool, error)

	Get(ctx context.Context, reportID string) (*storeType, bool, error)
	GetByQuery(ctx context.Context, query *v1.Query) ([]*storeType, error)
	GetMany(ctx context.Context, identifiers []string) ([]*storeType, []int, error)
	GetIDs(ctx context.Context) ([]string, error)

	Walk(ctx context.Context, fn func(obj *storeType) error) error
}

type storeImpl struct {
	*pgSearch.GenericStore[storeType, *storeType]
	mutex sync.RWMutex
}

// New returns a new Store instance using the provided sql instance.
func New(db postgres.DB) Store {
	return &storeImpl{
		GenericStore: pgSearch.NewGenericStore[storeType, *storeType](
			db,
			schema,
			pkGetter,
			insertIntoReportMetadata,
			copyFromReportMetadata,
			metricsSetAcquireDBConnDuration,
			metricsSetPostgresOperationDurationTime,
			pgSearch.GloballyScopedUpsertChecker[storeType, *storeType](targetResource),
			targetResource,
		),
	}
}

// region Helper functions

func pkGetter(obj *storeType) string {
	return obj.GetReportId()
}

func metricsSetPostgresOperationDurationTime(start time.Time, op ops.Op) {
	metrics.SetPostgresOperationDurationTime(start, op, storeName)
}

func metricsSetAcquireDBConnDuration(start time.Time, op ops.Op) {
	metrics.SetAcquireDBConnDuration(start, op, storeName)
}

func insertIntoReportMetadata(_ context.Context, batch *pgx.Batch, obj *storage.ReportMetadata) error {

	serialized, marshalErr := obj.Marshal()
	if marshalErr != nil {
		return marshalErr
	}

	values := []interface{}{
		// parent primary keys start
		pgutils.NilOrUUID(obj.GetReportId()),
		obj.GetReportConfigId(),
		obj.GetRequester().GetName(),
		obj.GetReportStatus().GetRunState(),
		pgutils.NilOrTime(obj.GetReportStatus().GetQueuedAt()),
		pgutils.NilOrTime(obj.GetReportStatus().GetCompletedAt()),
		obj.GetReportStatus().GetReportRequestType(),
		obj.GetReportStatus().GetReportNotificationMethod(),
		serialized,
	}

	finalStr := "INSERT INTO report_metadata (ReportId, ReportConfigId, Requester_Name, ReportStatus_RunState, ReportStatus_QueuedAt, ReportStatus_CompletedAt, ReportStatus_ReportRequestType, ReportStatus_ReportNotificationMethod, serialized) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT(ReportId) DO UPDATE SET ReportId = EXCLUDED.ReportId, ReportConfigId = EXCLUDED.ReportConfigId, Requester_Name = EXCLUDED.Requester_Name, ReportStatus_RunState = EXCLUDED.ReportStatus_RunState, ReportStatus_QueuedAt = EXCLUDED.ReportStatus_QueuedAt, ReportStatus_CompletedAt = EXCLUDED.ReportStatus_CompletedAt, ReportStatus_ReportRequestType = EXCLUDED.ReportStatus_ReportRequestType, ReportStatus_ReportNotificationMethod = EXCLUDED.ReportStatus_ReportNotificationMethod, serialized = EXCLUDED.serialized"
	batch.Queue(finalStr, values...)

	return nil
}

func copyFromReportMetadata(ctx context.Context, s pgSearch.Deleter, tx *postgres.Tx, objs ...*storage.ReportMetadata) error {
	inputRows := make([][]interface{}, 0, batchSize)

	// This is a copy so first we must delete the rows and re-add them
	// Which is essentially the desired behaviour of an upsert.
	deletes := make([]string, 0, batchSize)

	copyCols := []string{
		"reportid",
		"reportconfigid",
		"requester_name",
		"reportstatus_runstate",
		"reportstatus_queuedat",
		"reportstatus_completedat",
		"reportstatus_reportrequesttype",
		"reportstatus_reportnotificationmethod",
		"serialized",
	}

	for idx, obj := range objs {
		// Todo: ROX-9499 Figure out how to more cleanly template around this issue.
		log.Debugf("This is here for now because there is an issue with pods_TerminatedInstances where the obj "+
			"in the loop is not used as it only consists of the parent ID and the index.  Putting this here as a stop gap "+
			"to simply use the object.  %s", obj)

		serialized, marshalErr := obj.Marshal()
		if marshalErr != nil {
			return marshalErr
		}

		inputRows = append(inputRows, []interface{}{
			pgutils.NilOrUUID(obj.GetReportId()),
			obj.GetReportConfigId(),
			obj.GetRequester().GetName(),
			obj.GetReportStatus().GetRunState(),
			pgutils.NilOrTime(obj.GetReportStatus().GetQueuedAt()),
			pgutils.NilOrTime(obj.GetReportStatus().GetCompletedAt()),
			obj.GetReportStatus().GetReportRequestType(),
			obj.GetReportStatus().GetReportNotificationMethod(),
			serialized,
		})

		// Add the ID to be deleted.
		deletes = append(deletes, obj.GetReportId())

		// if we hit our batch size we need to push the data
		if (idx+1)%batchSize == 0 || idx == len(objs)-1 {
			// copy does not upsert so have to delete first.  parent deletion cascades so only need to
			// delete for the top level parent

			if err := s.DeleteMany(ctx, deletes); err != nil {
				return err
			}
			// clear the inserts and vals for the next batch
			deletes = deletes[:0]

			if _, err := tx.CopyFrom(ctx, pgx.Identifier{"report_metadata"}, copyCols, pgx.CopyFromRows(inputRows)); err != nil {
				return err
			}
			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	return nil
}

// endregion Helper functions

// region Used for testing

// CreateTableAndNewStore returns a new Store instance for testing.
func CreateTableAndNewStore(ctx context.Context, db postgres.DB, gormDB *gorm.DB) Store {
	pkgSchema.ApplySchemaForTable(ctx, gormDB, baseTable)
	return New(db)
}

// Destroy drops the tables associated with the target object type.
func Destroy(ctx context.Context, db postgres.DB) {
	dropTableReportMetadata(ctx, db)
}

func dropTableReportMetadata(ctx context.Context, db postgres.DB) {
	_, _ = db.Exec(ctx, "DROP TABLE IF EXISTS report_metadata CASCADE")

}

// endregion Used for testing
