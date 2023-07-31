package tracer

import (
	"context"

	"github.com/machtwatch/catalystdk/go/log"
	"github.com/machtwatch/catalystdk/go/trace"
	config "github.com/machtwatch/catalystdk/go/trace/config"
)

// const segments
const (
	SegmentHandler           = "Handler "
	SegmentUsecase           = "Usecase "
	SegmentRepoAPI           = "Repository API "
	SegmentRepoGRPC          = "Repository GRPC "
	SegmentRepoDB            = "Repository DB "
	SegmentRepoCache         = "Repository Cache "
	SegmentRepoAppCache      = "Repository App Cache "
	SegmentRepoHTTP          = "Repository HTTP "
	SegmentRepoRabbitMQ      = "Repository RabbitMQ "
	SegmentRepoElasticSearch = "Repository ElasticSearch "
	SegmentRundeck           = "Rundeck Service "
)

// SetStandardTrace init catalystdk log globally for standard logging with fields
func SetStandardTrace(config *config.Config) *trace.TracerSet {
	tc, err := trace.SetStdTrace(config)
	if err != nil {
		log.StdError(context.Background(), config, err, "got error on init tracer  - SetStandardTrace()")
	}
	return tc
}
