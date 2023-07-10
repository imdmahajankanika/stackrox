// Code generated by pg-bindings generator. DO NOT EDIT.

package postgres

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v4"
	"github.com/stackrox/rox/central/metrics"
	"github.com/stackrox/rox/central/role/resources"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	ops "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	pkgSchema "github.com/stackrox/rox/pkg/postgres/schema"
	"github.com/stackrox/rox/pkg/sac"
	pgSearch "github.com/stackrox/rox/pkg/search/postgres"
	"github.com/stackrox/rox/pkg/sync"
	"gorm.io/gorm"
)

const (
	baseTable = "report_metadata"
	storeName = "ReportMetadata"

	batchAfter = 100

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

// Store is the interface to interact with the storage for storage.ReportMetadata
type Store interface {
	Upsert(ctx context.Context, obj *storage.ReportMetadata) error
	UpsertMany(ctx context.Context, objs []*storage.ReportMetadata) error
	Delete(ctx context.Context, reportID string) error
	DeleteByQuery(ctx context.Context, q *v1.Query) error
	DeleteMany(ctx context.Context, identifiers []string) error

	Count(ctx context.Context) (int, error)
	Exists(ctx context.Context, reportID string) (bool, error)

	Get(ctx context.Context, reportID string) (*storage.ReportMetadata, bool, error)
	GetByQuery(ctx context.Context, query *v1.Query) ([]*storage.ReportMetadata, error)
	GetMany(ctx context.Context, identifiers []string) ([]*storage.ReportMetadata, []int, error)
	GetIDs(ctx context.Context) ([]string, error)

	Walk(ctx context.Context, fn func(obj *storage.ReportMetadata) error) error
}

type storeImpl struct {
	*pgSearch.GenericStore[storage.ReportMetadata, *storage.ReportMetadata]
	db    postgres.DB
	mutex sync.RWMutex
}

// New returns a new Store instance using the provided sql instance.
func New(db postgres.DB) Store {
	return &storeImpl{
		db: db,
		GenericStore: pgSearch.NewGenericStore[storage.ReportMetadata, *storage.ReportMetadata](
			db,
			schema,
			pkGetter,
			metricsSetAcquireDBConnDuration,
			metricsSetPostgresOperationDurationTime,
			targetResource,
		),
	}
}

// region Helper functions

func pkGetter(obj *storage.ReportMetadata) string {
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

func (s *storeImpl) copyFromReportMetadata(ctx context.Context, tx *postgres.Tx, objs ...*storage.ReportMetadata) error {

	inputRows := [][]interface{}{}

	var err error

	// This is a copy so first we must delete the rows and re-add them
	// Which is essentially the desired behaviour of an upsert.
	var deletes []string

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
			deletes = nil

			_, err = tx.CopyFrom(ctx, pgx.Identifier{"report_metadata"}, copyCols, pgx.CopyFromRows(inputRows))

			if err != nil {
				return err
			}

			// clear the input rows for the next batch
			inputRows = inputRows[:0]
		}
	}

	return err
}

func (s *storeImpl) copyFrom(ctx context.Context, objs ...*storage.ReportMetadata) error {
	conn, err := s.AcquireConn(ctx, ops.Get)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	if err := s.copyFromReportMetadata(ctx, tx, objs...); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (s *storeImpl) upsert(ctx context.Context, objs ...*storage.ReportMetadata) error {
	conn, err := s.AcquireConn(ctx, ops.Get)
	if err != nil {
		return err
	}
	defer conn.Release()

	for _, obj := range objs {
		batch := &pgx.Batch{}
		if err := insertIntoReportMetadata(ctx, batch, obj); err != nil {
			return err
		}
		batchResults := conn.SendBatch(ctx, batch)
		var result *multierror.Error
		for i := 0; i < batch.Len(); i++ {
			_, err := batchResults.Exec()
			result = multierror.Append(result, err)
		}
		if err := batchResults.Close(); err != nil {
			return err
		}
		if err := result.ErrorOrNil(); err != nil {
			return err
		}
	}
	return nil
}

// endregion Helper functions

//// Interface functions

// Upsert saves the current state of an object in storage.
func (s *storeImpl) Upsert(ctx context.Context, obj *storage.ReportMetadata) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.Upsert, "ReportMetadata")

	scopeChecker := sac.GlobalAccessScopeChecker(ctx).AccessMode(storage.Access_READ_WRITE_ACCESS).Resource(targetResource)
	if !scopeChecker.IsAllowed() {
		return sac.ErrResourceAccessDenied
	}

	return pgutils.Retry(func() error {
		return s.upsert(ctx, obj)
	})
}

// UpsertMany saves the state of multiple objects in the storage.
func (s *storeImpl) UpsertMany(ctx context.Context, objs []*storage.ReportMetadata) error {
	defer metrics.SetPostgresOperationDurationTime(time.Now(), ops.UpdateMany, "ReportMetadata")

	scopeChecker := sac.GlobalAccessScopeChecker(ctx).AccessMode(storage.Access_READ_WRITE_ACCESS).Resource(targetResource)
	if !scopeChecker.IsAllowed() {
		return sac.ErrResourceAccessDenied
	}

	return pgutils.Retry(func() error {
		// Lock since copyFrom requires a delete first before being executed.  If multiple processes are updating
		// same subset of rows, both deletes could occur before the copyFrom resulting in unique constraint
		// violations
		if len(objs) < batchAfter {
			s.mutex.RLock()
			defer s.mutex.RUnlock()

			return s.upsert(ctx, objs...)
		}
		s.mutex.Lock()
		defer s.mutex.Unlock()

		return s.copyFrom(ctx, objs...)
	})
}

//// Interface functions - END

//// Used for testing

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

//// Used for testing - END
