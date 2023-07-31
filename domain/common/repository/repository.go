package repository

import (
	"context"
	"database/sql"
	"math"

	sq "github.com/Masterminds/squirrel"
	"github.com/machtwatch/catalyst-go-skeleton/domain/common/model"
	request "github.com/machtwatch/catalyst-go-skeleton/domain/common/request/api"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/tracer"
	"github.com/machtwatch/catalystdk/go/log"
	"github.com/machtwatch/catalystdk/go/trace"
	"gorm.io/gorm"
)

type IBaseRepo[M model.EntityModel] interface {
	GetAll(ctx context.Context) ([]M, error)
	GetWithPagination(ctx context.Context, opt request.PaginationRequest, filters []interface{}) ([]M, PaginationResponse, error)
	GetByID(ctx context.Context, ID int64) (M, error)
	GetByIDs(ctx context.Context, IDs []int64, filters ...interface{}) ([]M, error)
	Create(ctx context.Context, model M) (M, error)
	CreateBulk(ctx context.Context, models []M) error
	Update(ctx context.Context, ID int64, model M) (M, error)
	UpdateBulk(ctx context.Context, IDs []int64, payload map[string]interface{}) error
	UpdateWithMap(ctx context.Context, ID int64, payload map[string]interface{}) (M, error)
	Delete(ctx context.Context, model M, ID int64) error
	DeleteBulk(ctx context.Context, IDs []int64) error
	CreateWithTx(ctx context.Context, model M, trx *gorm.DB) (M, error)
	CreateBulkWithTx(ctx context.Context, models []M, trx *gorm.DB) error
	UpdateWithTx(ctx context.Context, ID int64, model M, trx *gorm.DB) (M, error)
	UpdateBulkWithTx(ctx context.Context, IDs []int64, payload map[string]interface{}, trx *gorm.DB) error
	UpdateWithMapTx(ctx context.Context, ID int64, payload map[string]interface{}, trx *gorm.DB) (M, error)
	DeleteWithTx(ctx context.Context, model M, ID int64, trx *gorm.DB) error
	DeleteBulkWithTx(ctx context.Context, IDs []int64, trx *gorm.DB) error
	BeginTransaction(ctx context.Context) *gorm.DB
	Rollback(trx *gorm.DB) *gorm.DB
	Commit(trx *gorm.DB) *gorm.DB
}

const (
	segmentRepoDBCreate            = tracer.SegmentRepoDB + "BaseRepo.Create"
	segmentRepoDBUpdate            = tracer.SegmentRepoDB + "BaseRepo.Update"
	segmentRepoDBDelete            = tracer.SegmentRepoDB + "BaseRepo.Delete"
	segmentRepoDBCreateBulk        = tracer.SegmentRepoDB + "BaseRepo.CreateBulk"
	segmentRepoDBUpdateBulk        = tracer.SegmentRepoDB + "BaseRepo.UpdateBulk"
	segmentRepoDBDeleteBulk        = tracer.SegmentRepoDB + "BaseRepo.DeleteBulk"
	segmentRepoDBDeleteBulkWithTx  = tracer.SegmentRepoDB + "BaseRepo.DeleteBulkWithTx"
	segmentRepoDBGetPaginationMeta = tracer.SegmentRepoDB + "BaseRepo.GetPaginationMeta"
	segmentRepoDBGetAll            = tracer.SegmentRepoDB + "BaseRepo.GetAll"
	segmentRepoDBGetWithPagination = tracer.SegmentRepoDB + "BaseRepo.GetWithPagination"
	segmentRepoDBGetByID           = tracer.SegmentRepoDB + "BaseRepo.GetByID"
	segmentRepoDBGetByIDs          = tracer.SegmentRepoDB + "BaseRepo.GetByIDs"
	segmentRepoDBUpdateBulkWithTx  = tracer.SegmentRepoDB + "BaseRepo.UpdateBulkWithTx"
	segmentRepoDBUpdateWithTx      = tracer.SegmentRepoDB + "BaseRepo.UpdateWithTx"
	segmentRepoDBUpdateWithMap     = tracer.SegmentRepoDB + "BaseRepo.UpdateWithMap"
	segmentRepoDBUpdateWithMapTx   = tracer.SegmentRepoDB + "BaseRepo.UpdateWithMapTx"
	segmentRepoDBDeleteWithTx      = tracer.SegmentRepoDB + "BaseRepo.DeleteWithTx"
)

type (
	TracingContextKey struct{}
)

type BaseRepo[M model.EntityModel] struct {
	writeConn *gorm.DB
	readConn  *gorm.DB
}

func NewBaseRepo[M model.EntityModel](dbMaster *gorm.DB, dbSlave *gorm.DB) *BaseRepo[M] {
	return &BaseRepo[M]{
		writeConn: dbMaster,
		readConn:  dbSlave,
	}
}

func (r *BaseRepo[M]) GetAll(ctx context.Context) ([]M, error) {
	ctx, span := trace.StartSpanFromContext(ctx, segmentRepoDBGetAll)
	defer span.End()
	var models []M
	err := r.readConn.WithContext(ctx).Find(&models).Where("deleted_at IS NULL").Error
	if err != nil {
		trace.SetErrorSpan(span, err, "r.readConn.WithContext() got error - repository.GetAll")
		log.StdDebug(ctx, nil, err, "r.readConn.WithContext() got error - repository.GetAll")
		return nil, err
	}

	return models, nil
}

func (r *BaseRepo[M]) GetWithPagination(ctx context.Context, opt request.PaginationRequest, filters []interface{}) (models []M, pagination PaginationResponse, err error) {
	ctx, span := trace.StartSpanFromContext(ctx, segmentRepoDBGetWithPagination)
	defer span.End()
	var model M

	pagination, err = r.GetPaginationMeta(ctx, model, opt.Size, filters)
	if err != nil {
		trace.SetErrorSpan(span, err, "r.GetPaginationMeta() got error - repository.GetWithPagination")
		log.StdDebug(ctx, nil, err, "r.GetPaginationMeta() got error - repository.GetWithPagination")
		return models, pagination, err
	}

	builderFilter := sq.Select("*").From(model.TableName())
	builderFilter = SetFilter(builderFilter, filters)
	builderFilter = SetPagination(builderFilter, opt)

	qry, args, err := builderFilter.ToSql()
	if err != nil {
		trace.SetErrorSpan(span, err, "builderFilter.ToSql() got error - repository.GetWithPagination")
		log.StdDebug(ctx, map[string]interface{}{"opt": opt}, err, "builderFilter.ToSql() got error - repository.GetWithPagination")
		return models, pagination, err
	}

	err = r.readConn.WithContext(ctx).Raw(qry, args...).Scan(&models).Error
	if err != nil && err == sql.ErrNoRows {
		return models, pagination, nil
	}

	if err != nil {
		trace.SetErrorSpan(span, err, "r.readConn.WithContext() got error - repository.GetWithPagination")
		log.StdDebug(ctx, map[string]interface{}{"query": qry, "args": args}, err, "r.readConn.WithContext() got error - repository.GetWithPagination")
		return models, pagination, err
	}

	return models, pagination, nil
}

func (r *BaseRepo[M]) GetByID(ctx context.Context, ID int64) (M, error) {
	ctx, span := trace.StartSpanFromContext(ctx, segmentRepoDBGetByID)
	defer span.End()
	var model M

	builder := sq.Select("*").From(model.TableName()).Where(sq.Eq{"id": ID}).Where("deleted_at IS NULL")
	qry, args, err := builder.ToSql()
	if err != nil {
		trace.SetErrorSpan(span, err, "builder.ToSql() got error - repository.GetByID")
		log.StdDebug(ctx, map[string]interface{}{"id": ID, "args": args}, err, "builder.ToSql() got error - repository.GetByID")
		return model, err
	}

	err = r.readConn.WithContext(ctx).Raw(qry, args...).Scan(&model).Error
	if err != nil && err == sql.ErrNoRows {
		return model, nil
	}

	if err != nil {
		trace.SetErrorSpan(span, err, "builder.ToSql() got error - repository.GetByID")
		log.StdDebug(ctx, map[string]interface{}{"query": qry, "args": args}, err, "builder.ToSql() got error - repository.GetByID")
		return model, err
	}

	return model, nil
}

func (r *BaseRepo[M]) GetByIDs(ctx context.Context, IDs []int64, filters ...interface{}) (models []M, err error) {
	ctx, span := trace.StartSpanFromContext(ctx, segmentRepoDBGetByIDs)
	defer span.End()

	var ids []int64
	ids = append(ids, IDs...)

	var model M

	builder := sq.Select("*").From(model.TableName()).Where(sq.Eq{"id": ids})
	builder = SetFilter(builder, filters)
	qry, args, err := builder.ToSql()
	if err != nil {
		trace.SetErrorSpan(span, err, "builder.ToSql() got error - repository.GetByIDs")
		log.StdDebug(ctx, map[string]interface{}{"query": qry, "args": args}, err, "builder.ToSql() got error - repository.GetByIDs")
		return models, err
	}

	err = r.readConn.WithContext(ctx).Raw(qry, args...).Scan(&models).Error
	if err != nil && err == sql.ErrNoRows {
		// no need to logs, error no rows returned
		return models, nil
	}

	if err != nil {
		trace.SetErrorSpan(span, err, "r.readConn.WithContext() got error - repository.GetByIDs")
		log.StdDebug(ctx, map[string]interface{}{"query": qry}, err, "r.readConn.WithContext() got error - repository.GetByIDs")
		return models, err
	}

	return models, nil
}

// Create execute a single insert without specified transaction
func (r *BaseRepo[M]) Create(ctx context.Context, model M) (M, error) {
	ctx, span := trace.StartSpanFromContext(ctx, segmentRepoDBCreate)
	defer span.End()

	err := r.writeConn.WithContext(ctx).Create(&model).Error
	if err != nil {
		trace.SetErrorSpan(span, err, "r.writeConn.WithContext() got error - repository.Create")
		log.StdDebug(ctx, map[string]interface{}{"model": model}, err, "r.writeConn.WithContext() got error - repository.Create")
		return model, err
	}

	return model, nil
}

// CreateWithTx execute a single insert with specified transaction
func (r *BaseRepo[M]) CreateWithTx(ctx context.Context, model M, trx *gorm.DB) (M, error) {
	ctx, span := trace.StartSpanFromContext(ctx, segmentRepoDBCreate)
	defer span.End()

	err := trx.Create(&model).Error
	if err != nil {
		trace.SetErrorSpan(span, err, "trx.Create() got error - repository.CreateWithTx")
		log.StdDebug(ctx, map[string]interface{}{"model": model}, err, "trx.Create() got error - repository.CreateWithTx")
		return model, err
	}

	return model, nil
}

// CreateBulk execute a bulk insert without specified transaction
func (r *BaseRepo[M]) CreateBulk(ctx context.Context, models []M) error {
	ctx, span := trace.StartSpanFromContext(ctx, segmentRepoDBCreateBulk)
	defer span.End()

	err := r.writeConn.WithContext(ctx).CreateInBatches(models, InsertBatchSize).Error
	if err != nil {
		trace.SetErrorSpan(span, err, "r.writeConn.WithContext() got error - repository.CreateBulk")
		log.StdDebug(ctx, map[string]interface{}{"models": models}, err, "r.writeConn.WithContext() got error - repository.CreateBulk")
	}

	return err
}

// CreateBulkWithTx execute a bulk insert without specified transaction
func (r *BaseRepo[M]) CreateBulkWithTx(ctx context.Context, models []M, trx *gorm.DB) error {
	ctx, span := trace.StartSpanFromContext(ctx, segmentRepoDBCreateBulk)
	defer span.End()

	err := trx.WithContext(ctx).CreateInBatches(models, InsertBatchSize).Error
	if err != nil {
		trace.SetErrorSpan(span, err, "trx.CreateInBatches() got error - repository.CreateBulkWithTx")
		log.StdDebug(ctx, map[string]interface{}{"models": models, "InsertBatchSize": InsertBatchSize}, err, "trx.CreateInBatches() got error - repository.CreateBulkWithTx")
	}

	return err
}

// Update execute bulk update without specified transaction
func (r *BaseRepo[M]) Update(ctx context.Context, ID int64, model M) (M, error) {
	ctx, span := trace.StartSpanFromContext(ctx, segmentRepoDBUpdate)
	defer span.End()

	err := r.writeConn.WithContext(ctx).Model(&model).Where("id=? AND deleted_at is NULL", ID).Updates(model).Scan(&model).Error
	if err != nil {
		trace.SetErrorSpan(span, err, "r.writeConn.WithContext() got error - repository.Update")
		log.StdDebug(ctx, map[string]interface{}{"model": model}, err, "r.writeConn.WithContext() got error - repository.Update")
		return model, err
	}

	return model, nil
}

// UpdateWithTx execute a single update with specified transaction
func (r *BaseRepo[M]) UpdateWithTx(ctx context.Context, ID int64, model M, trx *gorm.DB) (M, error) {
	ctx, span := trace.StartSpanFromContext(ctx, segmentRepoDBUpdateWithTx)
	defer span.End()

	err := trx.WithContext(ctx).Model(&model).Where("id=? AND deleted_at is NULL", ID).Updates(model).Scan(&model).Error
	if err != nil {
		trace.SetErrorSpan(span, err, "trx.WithContext() got error - repository.UpdateWithTx")
		log.StdDebug(ctx, map[string]interface{}{"model": model}, err, "trx.WithContext() got error - repository.UpdateWithTx")
		return model, err
	}

	return model, nil
}

// UpdateBulk execute a bulk update without specified transaction
func (r *BaseRepo[M]) UpdateBulk(ctx context.Context, IDs []int64, payload map[string]interface{}) error {
	var model M

	ctx, span := trace.StartSpanFromContext(ctx, segmentRepoDBUpdateBulk)
	defer span.End()

	err := r.writeConn.WithContext(ctx).Model(&model).Where("id IN ? AND deleted_at is NULL", IDs).Updates(payload).Error
	if err != nil {
		trace.SetErrorSpan(span, err, "r.writeConn.WithContext() got error - repository.UpdateBulk")
		log.StdDebug(ctx, map[string]interface{}{"model": model}, err, "r.writeConn.WithContext() got error - repository.UpdateBulk")
		return err
	}

	return nil
}

// UpdateBulkWithTx execute a bulk update without specified transaction
func (r *BaseRepo[M]) UpdateBulkWithTx(ctx context.Context, IDs []int64, payload map[string]interface{}, trx *gorm.DB) error {
	var model M

	ctx, span := trace.StartSpanFromContext(ctx, segmentRepoDBUpdateBulkWithTx)
	defer span.End()

	err := trx.WithContext(ctx).Model(&model).Where("id IN ? AND deleted_at is NULL", IDs).Updates(payload).Error
	if err != nil {
		trace.SetErrorSpan(span, err, "trx.WithContext() got error - repository.UpdateBulk")
		log.StdDebug(ctx, map[string]interface{}{"model": model}, err, "trx.WithContext() got error - repository.UpdateBulk")
		return err
	}
	return nil

}

// UpdateWithMap execute a single update with Map without specified transaction
func (r *BaseRepo[M]) UpdateWithMap(ctx context.Context, ID int64, payload map[string]interface{}) (M, error) {
	var model M

	ctx, span := trace.StartSpanFromContext(ctx, segmentRepoDBUpdateWithMap)
	defer span.End()

	err := r.writeConn.WithContext(ctx).Model(&model).Where("id=? AND deleted_at is NULL", ID).Updates(payload).Scan(&model).Error
	if err != nil {
		trace.SetErrorSpan(span, err, "r.writeConn.WithContext() got error - repository.UpdateWithMap")
		log.StdDebug(ctx, map[string]interface{}{"model": model}, err, "r.writeConn.WithContext() got error - repository.UpdateWithMap")
		return model, err
	}

	return model, nil
}

// UpdateWithMapTx execute a single update with Map with specified transaction
func (r *BaseRepo[M]) UpdateWithMapTx(ctx context.Context, ID int64, payload map[string]interface{}, trx *gorm.DB) (M, error) {
	var model M

	ctx, span := trace.StartSpanFromContext(ctx, segmentRepoDBUpdateWithMapTx)
	defer span.End()

	err := trx.WithContext(ctx).Model(&model).Where("id=? AND deleted_at is NULL", ID).Updates(payload).Scan(&model).Error
	if err != nil {
		trace.SetErrorSpan(span, err, "trx.WithContext() got error - repository.segmentRepoDBUpdateWithMapTx")
		log.StdDebug(ctx, map[string]interface{}{"model": model}, err, "trx.WithContext() got error - repository.segmentRepoDBUpdateWithMapTx")
		return model, err
	}

	return model, nil
}

// Delete execute a single delete without specified transaction
func (r *BaseRepo[M]) Delete(ctx context.Context, model M, ID int64) error {
	ctx, span := trace.StartSpanFromContext(ctx, segmentRepoDBDelete)
	defer span.End()

	err := r.writeConn.WithContext(ctx).Delete(&model, ID).Error
	if err != nil {
		trace.SetErrorSpan(span, err, "r.writeConn.WithContext() got error - repository.Delete")
		log.StdDebug(ctx, map[string]interface{}{"model": model}, err, "r.writeConn.WithContext() got error - repository.Delete")
		return err
	}

	return nil
}

// Delete execute a single delete with specified transaction
func (r *BaseRepo[M]) DeleteWithTx(ctx context.Context, model M, ID int64, trx *gorm.DB) error {
	ctx, span := trace.StartSpanFromContext(ctx, segmentRepoDBDeleteWithTx)
	defer span.End()

	err := trx.WithContext(ctx).Delete(&model, ID).Error
	if err != nil {
		trace.SetErrorSpan(span, err, "r.writeConn.WithContext() got error - repository.DeleteWithTx")
		log.StdDebug(ctx, map[string]interface{}{"model": model}, err, "r.writeConn.WithContext() got error - repository.DeleteWithTx")
		return err
	}

	return nil
}

// DeleteBulk execute a `soft` bulk delete without specified transaction
func (r *BaseRepo[M]) DeleteBulk(ctx context.Context, IDs []int64) error {
	ctx, span := trace.StartSpanFromContext(ctx, segmentRepoDBDeleteBulk)
	defer span.End()

	var model M
	var err error

	builder := sq.Update(model.TableName()).Set("deleted_at", "now()").Where(sq.Eq{"id": IDs})

	qry, args, err := builder.ToSql()
	if err != nil {
		trace.SetErrorSpan(span, err, "builder.ToSql() got error - DeleteBulk")
		log.StdDebug(ctx, map[string]interface{}{"model": model}, err, "builder.ToSql() got error - DeleteBulk")
		return err
	}

	err = r.writeConn.WithContext(ctx).Exec(qry, args...).Error
	if err != nil {
		trace.SetErrorSpan(span, err, "r.writeConn.WithContext() got error - DeleteBulk")
		log.StdDebug(ctx, map[string]interface{}{"model": model}, err, "r.writeConn.WithContext() got error - DeleteBulk")
		return err
	}

	return nil
}

// DeleteBulk execute a bulk `soft` delete with specified transaction
func (r *BaseRepo[M]) DeleteBulkWithTx(ctx context.Context, IDs []int64, trx *gorm.DB) error {
	ctx, span := trace.StartSpanFromContext(ctx, segmentRepoDBDeleteBulkWithTx)
	defer span.End()

	var model M
	var err error

	builder := sq.Update(model.TableName()).Set("deleted_at", "now()").Where(sq.Eq{"id": IDs})

	qry, args, err := builder.ToSql()
	if err != nil {
		trace.SetErrorSpan(span, err, "builder.ToSql() got error - DeleteBulkWithTx")
		log.StdDebug(ctx, map[string]interface{}{"model": model}, err, "builder.ToSql() got error - DeleteBulkWithTx")
		return err
	}

	err = trx.WithContext(ctx).Exec(qry, args...).Error
	if err != nil {
		trace.SetErrorSpan(span, err, "trx.WithContext() got error - DeleteBulkWithTx")
		log.StdDebug(ctx, map[string]interface{}{"model": model}, err, "trx.WithContext() got error - DeleteBulkWithTx")
		return err
	}

	return nil
}

func (r *BaseRepo[M]) GetPaginationMeta(ctx context.Context, model M, pageSize int, filters []interface{}) (pagination PaginationResponse, err error) {
	ctx, span := trace.StartSpanFromContext(ctx, segmentRepoDBGetPaginationMeta)
	defer span.End()

	builderTotal := sq.Select("COUNT(id) AS totalRows").From(model.TableName())
	builderTotal = SetFilter(builderTotal, filters)

	qry, args, err := builderTotal.ToSql()
	if err != nil {
		trace.SetErrorSpan(span, err, "builderTotal.ToSql() got error - GetPaginationMeta")
		log.StdDebug(ctx, map[string]interface{}{"model": model}, err, " builderTotal.ToSql() got error - GetPaginationMeta")
		return pagination, err
	}

	var totalRows int64
	err = r.readConn.WithContext(ctx).Model(model).Raw(qry, args...).Scan(&totalRows).Error
	if err != nil {
		trace.SetErrorSpan(span, err, "r.readConn.WithContext() got error - GetPaginationMeta")
		log.StdDebug(ctx, map[string]interface{}{"model": model}, err, " r.readConn.WithContext() got error - GetPaginationMeta")
		return pagination, err
	}

	pagination.TotalData = int(totalRows)

	totalPages := int(math.Ceil(float64(totalRows) / float64(pageSize)))
	pagination.TotalPage = totalPages

	return pagination, nil
}

func (r *BaseRepo[M]) BeginTransaction(ctx context.Context) *gorm.DB {
	return r.writeConn.WithContext(ctx).Begin()
}
func (r *BaseRepo[M]) Rollback(trx *gorm.DB) *gorm.DB {
	return trx.Rollback()
}
func (r *BaseRepo[M]) Commit(trx *gorm.DB) *gorm.DB {
	return trx.Commit()
}
