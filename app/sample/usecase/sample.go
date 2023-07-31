package usecase

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/machtwatch/catalyst-go-skeleton/domain"
	"github.com/machtwatch/catalyst-go-skeleton/domain/common/event"
	"github.com/machtwatch/catalyst-go-skeleton/domain/common/response"
	"github.com/machtwatch/catalyst-go-skeleton/domain/sample"
	samplereserve "github.com/machtwatch/catalyst-go-skeleton/domain/sample_reserve"
	scheduledjob "github.com/machtwatch/catalyst-go-skeleton/domain/scheduler"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config/secret"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/eventstore"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/integration"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/scheduler"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/tracer"
	flagprovider "github.com/machtwatch/catalystdk/go/flag/provider"
	"github.com/machtwatch/catalystdk/go/log"
	"github.com/machtwatch/catalystdk/go/trace"
)

const (
	segmentUsecaseCreateOrder               = tracer.SegmentUsecase + "SampleUseCase.CreateOrder"
	segmentUsecaseCreateSchedule            = tracer.SegmentUsecase + "SampleUseCase.CreateSchedule"
	segmentUsecaseReceiveMessageByScheduler = tracer.SegmentUsecase + "SampleUseCase.ReceiveMessageByScheduler"
	SAMPLE_FORMAT_DATE                      = "2006-01-02 15:04:05"
)

// Usecase type of sample usecase.
//
// It contains the repositories and messaging event producer.
type SampleUC struct {
	sampleRepo        sample.SampleRepo
	sampleReserveRepo samplereserve.SampleReserveRepo
	scheduler         scheduler.ISchedulerExecution
	eventProducer     messaging.Producer
	eventStore        *eventstore.PgEventStore
	flag              *flagprovider.Flag
}

// NewSampleUC instantiate sample usecase.
//
// Accept domain repo collections, event producer, and xms integration. It will return the instance of usecase.
func NewSampleUC(repoCollection domain.RepositoryCollection, eventProducer messaging.Producer, eventstore *eventstore.PgEventStore, scheduler scheduler.ISchedulerExecution, xms integration.IXMSService, flag *flagprovider.Flag) sample.SampleUC {
	return &SampleUC{
		sampleRepo:        repoCollection.SampleRepo,
		sampleReserveRepo: repoCollection.SampleReserveRepo,
		scheduler:         scheduler,
		eventProducer:     eventProducer,
		eventStore:        eventstore,
		flag:              flag,
	}
}

// CreateOrder create a sample order as an example on how to use this sample repo
//
// It returns success response as true if success
// Otherwise, it returns nil response with false flag.
func (uc *SampleUC) CreateOrder(ctx context.Context, req sample.CreateOrderRequest, token string) (res response.Response[bool]) {
	ctx, span := trace.StartSpanFromContext(ctx, segmentUsecaseCreateOrder)
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, config.CONTEXT_TIMEOUT)
	defer cancel()

	// check flag
	if uc.flag.Feature(ctx, "create-order-feature").Off {
		return res.NotFound(ctx, "Not Implemented")
	}

	// begin transaction
	trx := uc.sampleRepo.BeginTransaction(ctx)

	loc, err := time.LoadLocation(config.DB_TIMEZONE)
	sampleEntity := sample.SampleEntity{
		OrderDate: time.Now().In(loc),
		Note:      "test notes",
	}
	if err != nil {
		log.StdDebug(ctx, req, err, "time.LoadLocation() got error - uc.CreateOrder()")
		return res.InternalError(ctx, err, "Gagal membuat order, terjadi kesalahan saat menyimpan data order")
	}

	// insert order
	_, err = uc.sampleRepo.CreateWithTx(ctx, sampleEntity, trx)
	if err != nil {
		trx.Rollback()
		log.StdDebug(ctx, req, err, "uc.sampleRepo.Create() Failed to create order on insert - uc.CreateOrder()")
		return res.InternalError(ctx, err, "Gagal membuat order, terjadi kesalahan saat menyimpan data order")
	}

	sampleReserveEntity := samplereserve.SampleReserveEntity{
		OfficeID:  req.OfficeID,
		OrderDate: time.Now().In(loc),
	}

	// insert reserve order
	_, err = uc.sampleReserveRepo.CreateWithTx(ctx, sampleReserveEntity, trx)
	if err != nil {
		trx.Rollback()
		log.StdDebug(ctx, req, err, "uc.sampleReserveRepo.Create() Failed to create order on insert - uc.CreateOrder()")
		return res.InternalError(ctx, err, "Gagal membuat order reserve, terjadi kesalahan saat menyimpan data order")
	}

	if err := uc.PublishEvent(ctx, req); err != nil {
		log.StdDebug(ctx, req, err, "uc.PublishEvent() Error publishing event - uc.CreateOrder()")
		trx.Rollback()
		return res.InternalError(ctx, err, "Gagal membuat order, terjadi kesalahan saat mempublish order")
	}

	trx.Commit()

	return res.Success(ctx, true, "sukses membuat order")
}

// PublishEvent create a sample order as an example on how to use this sample repo
//
// It returns success response as true if success
// Otherwise, it returns nil response with false flag.
func (uc *SampleUC) PublishEvent(ctx context.Context, req sample.CreateOrderRequest) error {

	ctx, cancel := context.WithTimeout(ctx, config.CONTEXT_TIMEOUT)
	defer cancel()

	sampleEvent := &event.SampleUpdated{
		Id:      12345,
		OrderId: 56789,
		Status:  "success",
	}

	if err := uc.eventStore.SaveEvent(ctx, sampleEvent, event.SAMPLE_TOPIC); err != nil {
		log.StdDebug(ctx, sampleEvent, err, "uc.eventStore.SaveEvent() Error saving sample event - uc.PublishEvent()")
		return err
	}

	if err := uc.eventProducer.Publish(ctx, sampleEvent, event.SAMPLE_TOPIC); err != nil {
		log.StdDebug(ctx, sampleEvent, err, "uc.eventProducer.Publish() Error publish sample event - uc.PublishEvent()")

		if errSave := uc.eventStore.SaveFailedEvent(ctx, sampleEvent, event.SAMPLE_TOPIC, err.Error()); errSave != nil {
			log.StdDebug(ctx, sampleEvent, errSave, "uc.eventStore.SaveFailedEvent Error saving sample failed event - uc.PublishEvent()")
		}

		return err
	}

	log.StdInfo(ctx, nil, "sukses mempublish event")

	return nil
}

// ConsumeEvent consume a sample order as an example on how to use this sample repo
//
// It returns success response as true if success
// Otherwise, it returns nil response with false flag.
func (uc *SampleUC) ConsumeEvent(ctx context.Context, data *event.SampleUpdated) error {

	log.StdInfo(ctx, data, fmt.Sprintf("Event data: %+v", data))

	return nil
}

// CreateSchedule create a sample request schedule to rundeck as an example on how to use scheduler
//
// It returns success response as true if success
// Otherwise, it returns nil response with false flag.
func (uc *SampleUC) CreateSchedule(ctx context.Context, req sample.CreateScheduleRequest) (res response.Response[interface{}]) {

	ctx, span := trace.StartSpanFromContext(ctx, segmentUsecaseCreateSchedule)
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, config.CONTEXT_TIMEOUT)
	defer cancel()

	// in this sample app, we set message for rundeck request body, endpoint for url, and token for authorization header
	// for endpoint using app_url in environment, if empty get from os hostname
	url := config.APP_URL
	if url == "" {
		hostname, err := os.Hostname()
		if err != nil {
			return res.InternalError(ctx, err, "Failed getting hostname")
		}
		url = hostname
	}
	// if your app_url using domain that direct this service, you don't have to adding app_port to url
	// url = fmt.Sprintf("%s:%d", url, config.APP_PORT)
	payload := map[string]string{
		"message":  req.Message,
		"endpoint": url + "/scheduler/sample",
		"token":    secret.RUNDECK_JOB_BEARER_TOKEN,
	}

	// because rundeck is strict with timezone, we must set it when convert time.Time
	jakartaLoc, _ := time.LoadLocation("Asia/Jakarta")
	execJobDate, err := time.ParseInLocation(SAMPLE_FORMAT_DATE, req.RunAt, jakartaLoc)
	if err != nil {
		return res.InternalError(ctx, err, "Failed format timezone jakarta")
	}

	// you can get rundeck job id when creating new job
	job := scheduledjob.SchedulerJob{
		ID:    secret.RUNDECK_JOB_ID,
		Opt:   payload,
		RunAt: execJobDate,
	}

	execID, err := uc.scheduler.RunJob(ctx, job)
	if err != nil {
		return res.InternalError(ctx, err, "Failed request scheduler")
	}

	data := map[string]interface{}{}
	data["message"] = fmt.Sprintf("created job at %s with execID %v", req.RunAt, execID)

	return res.Success(ctx, data, "success")
}

// ReceiveMessageByScheduler receiving a sample request message from rundeck as an example on how to use scheduler
//
// It returns success response as true if success
// Otherwise, it returns nil response with false flag.
func (uc *SampleUC) ReceiveMessageByScheduler(ctx context.Context, payload map[string]interface{}) (res response.Response[interface{}]) {
	ctx, span := trace.StartSpanFromContext(ctx, segmentUsecaseReceiveMessageByScheduler)
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, config.CONTEXT_TIMEOUT)
	defer cancel()

	// this is sample receive message. in your case, put your logic here
	if payload["message"] != "test rundeck" {
		return res.InvalidPayload(ctx, "message not match")
	}

	return res.Success(ctx, payload, "success")
}
