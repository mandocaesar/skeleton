package plugin

import (
	"fmt"

	"github.com/machtwatch/catalyst-go-skeleton/domain/common/repository"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config/secret"
	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/gorm"
)

const (
	gormTransactionOperation = "COMMIT/ROLLBACK"
	newrelicSegmentKey       = "newrelic_segment:%v"
)

// NewRelicTracer is a gorm plugin to instrument newrelic datastore segment for DB operations.
// It implements gorm Plugin interface which require Name() and Initialize() method.
type NewRelicTracer struct{}

// NewTracer returns new NewRelicTracer.
func NewTracer() NewRelicTracer {
	return NewRelicTracer{}
}

// Name returns the plugin name
func (nr NewRelicTracer) Name() string {
	return "newrelic-tracer-plugin"
}

// Initialize is called when adding an plugin to a gorm DB instance.
// Every gorm operation has a `Callback` which can be registered with custom functionality.
// This plugin registers `before` and `after` callback on each operation to start and end newrelic datastore segment.
func (nr NewRelicTracer) Initialize(db *gorm.DB) error {
	db.Callback().Create().Before("gorm:create").Register("newrelic:before_create", before)
	db.Callback().Create().After("gorm:create").Register("newrelic:after_create", after)
	db.Callback().Create().Before("gorm:begin_transaction").Register("newrelic:begin_transaction_create", beginTransaction)
	db.Callback().Create().After("gorm:commit_or_rollback_transaction").Register("newrelic:commit_or_rollback_transaction_create", commitOrRollback)

	db.Callback().Update().Before("gorm:update").Register("newrelic:before_update", before)
	db.Callback().Update().After("gorm:update").Register("newrelic:after_update", after)
	db.Callback().Update().Before("gorm:begin_transaction").Register("newrelic:begin_transaction_update", beginTransaction)
	db.Callback().Update().After("gorm:commit_or_rollback_transaction").Register("newrelic:commit_or_rollback_transaction_update", commitOrRollback)

	db.Callback().Delete().Before("gorm:delete").Register("newrelic:before_delete", before)
	db.Callback().Delete().After("gorm:delete").Register("newrelic:after_delete", after)
	db.Callback().Delete().Before("gorm:begin_transaction").Register("newrelic:begin_transaction_delete", beginTransaction)
	db.Callback().Delete().After("gorm:commit_or_rollback_transaction").Register("newrelic:commit_or_rollback_transaction_delete", commitOrRollback)

	db.Callback().Query().Before("gorm:query").Register("newrelic:before_query", before)
	db.Callback().Query().After("gorm:query").Register("newrelic:after_query", after)

	db.Callback().Row().Before("gorm:row").Register("newrelic:before_row", before)
	db.Callback().Row().After("gorm:row").Register("newrelic:after_row", after)

	db.Callback().Raw().Before("gorm:raw").Register("newrelic:before_raw", before)
	db.Callback().Raw().After("gorm:raw").Register("newrelic:after_raw", after)

	return nil
}

// before starts a newrelic transaction segment.
func before(db *gorm.DB) {
	ctx := db.Statement.Context
	ctxData := ctx.Value(repository.TracingContextKey{})

	segment, ok := ctxData.(*newrelic.DatastoreSegment)
	if !ok {
		return
	}

	tx := newrelic.FromContext(ctx)
	if tx == nil {
		return
	}

	segment.ParameterizedQuery = db.Statement.SQL.String()
	segment.StartTime = tx.StartSegmentNow()
}

// after ends a newrelic transaction segment.
func after(db *gorm.DB) {
	ctx := db.Statement.Context
	ctxData := ctx.Value(repository.TracingContextKey{})

	segment, ok := ctxData.(*newrelic.DatastoreSegment)
	if !ok {
		return
	}

	segment.End()
}

// beginTransaction starts a segment for DB transaction operation.
func beginTransaction(db *gorm.DB) {
	ctx := db.Statement.Context

	tx := newrelic.FromContext(ctx)
	if tx == nil {
		return
	}

	segment := newDatastoreSegment(gormTransactionOperation, db.Statement.Table)
	segment.StartTime = tx.StartSegmentNow()

	db.Statement.Set(segmentKey(gormTransactionOperation), segment)
}

// commitOrRollback ends a segment on DB transaction's commit or rollback.
func commitOrRollback(db *gorm.DB) {
	segment, ok := db.Statement.Get(segmentKey(gormTransactionOperation))
	if !ok {
		return
	}

	segment.(*newrelic.DatastoreSegment).ParameterizedQuery = db.Statement.SQL.String()
	segment.(*newrelic.DatastoreSegment).End()
}

// newDatastoreSegment returns a *newrelic.DatastoreSegment
func newDatastoreSegment(operation string, collection string) *newrelic.DatastoreSegment {
	return &newrelic.DatastoreSegment{
		Product:      newrelic.DatastorePostgres,
		Host:         secret.POSTGRES_HOST_MASTER,
		DatabaseName: secret.POSTGRES_DATABASE_MASTER,
		PortPathOrID: fmt.Sprint(secret.POSTGRES_PORT_MASTER),
		Operation:    operation,  // Operation is the relevant action, e.g. "SELECT", "CREATE" etc.
		Collection:   collection, // Collection is the table or group being operated upon in the datastore.
	}
}

func segmentKey(operation string) string {
	return fmt.Sprintf(newrelicSegmentKey, operation)
}
