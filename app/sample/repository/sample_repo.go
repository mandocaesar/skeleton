package repository

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/machtwatch/catalyst-go-skeleton/domain/common/repository"
	"github.com/machtwatch/catalyst-go-skeleton/domain/common/response"
	"github.com/machtwatch/catalyst-go-skeleton/domain/sample"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/cache"
	"github.com/machtwatch/catalyst-go-skeleton/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// sampleRepo represent repository type of sample that used as
// a collection oh the sample repo.
type sampleRepo struct {
	repository.BaseRepo[sample.SampleEntity]
	http       *resty.Client
	cache      cache.Cache
	utils      utils.Utils
	writeConn  *gorm.DB
	readConn   *gorm.DB
	monggoConn *mongo.Database
}

// NewSampleRepo instantiate sample repository.
// Accepts base repo and some clients provider
func NewSampleRepo(dbMaster *gorm.DB, dbSlave *gorm.DB, mongoDb *mongo.Database, cache cache.Cache, httpClient *resty.Client, utils utils.Utils) sample.SampleRepo {
	baseRepo := repository.NewBaseRepo[sample.SampleEntity](dbMaster, dbSlave)
	return &sampleRepo{
		BaseRepo:   *baseRepo,
		http:       httpClient,
		writeConn:  dbMaster,
		readConn:   dbSlave,
		monggoConn: mongoDb,
		cache:      cache,
		utils:      utils,
	}
}

func (r *sampleRepo) GetSampleRequest(ctx context.Context, req sample.SampleListRequest) (orders []sample.SampleJoinModel, pagination response.MetaPagination, err error) {

	return orders, pagination, err
}
