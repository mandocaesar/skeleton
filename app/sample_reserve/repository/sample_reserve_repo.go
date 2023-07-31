package repository

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/machtwatch/catalyst-go-skeleton/domain/common/repository"
	samplereserve "github.com/machtwatch/catalyst-go-skeleton/domain/sample_reserve"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/cache"
	"github.com/machtwatch/catalyst-go-skeleton/utils"
	"gorm.io/gorm"
)

// sampleReserveRepo represent repository type of sample that used as
// a collection oh the sample repo.
type sampleReserveRepo struct {
	repository.BaseRepo[samplereserve.SampleReserveEntity]
	http      *resty.Client
	cache     cache.Cache
	utils     utils.Utils
	writeConn *gorm.DB
	readConn  *gorm.DB
}

// NewSampleReserveRepo instantiate sample repository.
// Accepts base repo and some clients provider
func NewSampleReserveRepo(dbMaster *gorm.DB, dbSlave *gorm.DB, cache cache.Cache, httpClient *resty.Client, utils utils.Utils) samplereserve.SampleReserveRepo {
	baseRepo := repository.NewBaseRepo[samplereserve.SampleReserveEntity](dbMaster, dbSlave)
	return &sampleReserveRepo{
		BaseRepo:  *baseRepo,
		http:      httpClient,
		writeConn: dbMaster,
		readConn:  dbSlave,
		cache:     cache,
		utils:     utils,
	}
}

func (r *sampleReserveRepo) GetByOrderID(ctx context.Context, orderID int64, filter interface{}) (sampleReserveEntity []samplereserve.SampleReserveEntity, err error) {

	return sampleReserveEntity, err
}
