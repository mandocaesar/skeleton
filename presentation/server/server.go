package server

import (
	"context"

	sampleRepo "github.com/machtwatch/catalyst-go-skeleton/app/sample/repository"
	sampleUC "github.com/machtwatch/catalyst-go-skeleton/app/sample/usecase"
	sampleReserveRepo "github.com/machtwatch/catalyst-go-skeleton/app/sample_reserve/repository"
	"github.com/machtwatch/catalyst-go-skeleton/domain/common/event"
	cache "github.com/machtwatch/catalyst-go-skeleton/infrastructure/cache"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config/secret"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/database"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/eventstore"
	flag "github.com/machtwatch/catalyst-go-skeleton/infrastructure/flag"
	http "github.com/machtwatch/catalyst-go-skeleton/infrastructure/http"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/integration"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging/rabbitmq/topics_factory"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/scheduler"

	"github.com/machtwatch/catalyst-go-skeleton/domain"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/logger"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/tracer"
	"github.com/machtwatch/catalyst-go-skeleton/presentation/api"
	"github.com/machtwatch/catalyst-go-skeleton/presentation/messaging"
	"github.com/machtwatch/catalyst-go-skeleton/presentation/middlewares"
	"github.com/machtwatch/catalyst-go-skeleton/utils"

	"github.com/go-chi/chi/v5"
)

// Start the application instance.
//
// This will build all the required system stacks for catalyst-go-skeleton,
// and then run the application with graceful shutdown.
func Start() {
	ctx := context.Background()

	config.SetConfig(".", ".env")

	logger.SetStandardLog(config.GetLogConfig())

	traceProvider := tracer.SetStandardTrace(config.GetTracerConfig())

	dbConnections := database.CreateDBConnection(ctx, database.GetConfigs())
	sampleDBRepo := dbConnections.Postgres[database.SamplePostgresDBName]
	sampleMongoRepo := dbConnections.Mongo[database.SampleMongoDBName]

	eventStore := eventstore.NewPgEventStore(dbConnections.Postgres[database.SamplePostgresDBName].Master)
	eventProducer := topics_factory.NewTopicsProducer(event.DOMAIN_EVENT_EXCHANGE, config.GetRmqConfig())

	httpClient := http.DefaultHTTPClient()

	redisConfig := cache.RedisConfig{
		Host: config.REDIS_HOST,
		Port: config.REDIS_PORT,
		DB:   config.REDIS_DB,
	}

	redis := cache.CreateRedisConnection(redisConfig)
	cache := cache.NewCache(redis)

	schedulerConfig := &scheduler.SchedulerConfig{
		RundeckUrl:      secret.RUNDECK_URL,
		RundeckAPIToken: secret.RUNDECK_API_TOKEN,
		RundeckProject:  secret.RUNDECK_PROJECT,
	}

	scheduler := scheduler.NewRundeckConnection(schedulerConfig)

	xmsSrv := integration.NewXMSService(config.XMS_API_URI)

	utls := utils.NewUtils()

	flag := flag.NewFlag(config.GetFlagConfig())

	repositoryCollection := domain.RepositoryCollection{
		SampleRepo:        sampleRepo.NewSampleRepo(sampleDBRepo.Master, sampleDBRepo.Slave, sampleMongoRepo.DB, cache, httpClient, utls),
		SampleReserveRepo: sampleReserveRepo.NewSampleReserveRepo(sampleDBRepo.Master, sampleDBRepo.Slave, cache, httpClient, utls),
	}

	usecaseCollection := domain.UsecaseCollection{
		SampleUC: sampleUC.NewSampleUC(repositoryCollection, eventProducer, eventStore, scheduler, xmsSrv, flag),
	}

	router := chi.NewRouter()

	middlewares.Set(router)
	apiMiddleware := middlewares.NewApiMiddleware(cache)
	api.Set(router, usecaseCollection, apiMiddleware)

	messaging.Set(usecaseCollection)

	addHealthCheck(router, dbConnections, cache)
	startServerWithGracefulShutdown(router, traceProvider)
}
