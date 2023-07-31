package messaging

import (
	"context"

	"github.com/machtwatch/catalyst-go-skeleton/domain/common/event"
	"github.com/machtwatch/catalyst-go-skeleton/domain/sample"
	"github.com/machtwatch/catalystdk/go/log"
)

// SampleUpdatedHandler represent sample updated messaging handler
// that provide list of use case needed.
type SampleUpdatedHandler struct {
	sampleUC sample.SampleUC
}

// NewSampleUpdatedHandler instantiate sample update handler for messaging
//
// It Accept the use case and will return the instance of handler.
func NewSampleUpdatedHandler(sampleUC sample.SampleUC) SampleUpdatedHandler {
	return SampleUpdatedHandler{
		sampleUC: sampleUC,
	}
}

// Handle handle the fulfillment updated event from messaging
func (h SampleUpdatedHandler) Handle(event *event.SampleUpdated) error {
	log.StdInfof(context.Background(), nil, "received SampleUpdated event with payload %+v - (h SampleUpdatedHandler)Handle)", event)
	// TODO: handle fulfillment events

	return nil
}
