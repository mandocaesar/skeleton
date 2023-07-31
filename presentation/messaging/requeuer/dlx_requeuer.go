package requeuer

import (
	"context"
	"fmt"

	"github.com/machtwatch/catalyst-go-skeleton/domain/common/event"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging/rabbitmq/direct_factory"
	"github.com/machtwatch/catalystdk/go/log"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type DLXRequeuer struct {
	parkingLotProducer messaging.Producer
}

func NewDLXRequeuer() DLXRequeuer {
	return DLXRequeuer{
		parkingLotProducer: direct_factory.NewDirectProducer(event.PARKING_LOT_EXCHANGE, config.GetRmqConfig()),
	}
}

// Handle handles event requeue mechanism via Dead Letter Exchange
func (r DLXRequeuer) Handle(delivery amqp091.Delivery, msg proto.Message, err error) error {
	ctx := context.Background()

	if delivery.Headers["x-death"] != nil {

		deathCount, err := r.getDeathCount(delivery)
		if err != nil {
			log.StdError(ctx, delivery, err, "failed to get message death count")
			return delivery.Ack(false)
		}

		if deathCount >= config.EVENT_MAX_RETRY_COUNT {
			log.StdError(ctx, delivery, err, "event max retry count reached")

			// Publish to parking lot queue
			r.parkingLotProducer.Publish(ctx, msg, event.PARKING_LOT_QUEUE)

			return delivery.Ack(false)
		}
	}

	// Reject with requeue false, message will be delivered to Dead Letter Exchange if configured
	return delivery.Reject(false)
}

func (r DLXRequeuer) getDeathCount(delivery amqp091.Delivery) (count int64, err error) {
	xDeaths, ok := delivery.Headers["x-death"].([]interface{})
	if !ok || len(xDeaths) == 0 {
		return 0, fmt.Errorf("failed to parse x-death header")
	}

	deathHeader, ok := xDeaths[0].(amqp091.Table)
	if !ok {
		return 0, fmt.Errorf("failed to parse x-death payload")
	}

	count, ok = deathHeader["count"].(int64)
	if !ok {
		return 0, fmt.Errorf("failed to parse x-death count")
	}

	return count, nil
}
