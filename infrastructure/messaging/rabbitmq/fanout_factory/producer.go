package fanout_factory

import (
	"context"

	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging/common/config"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging/rabbitmq/core"
	"github.com/machtwatch/catalystdk/go/log"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type FanoutProducer struct {
	coreProducer *core.Producer
}

// NewFanoutProducer creates messages producer with fanout exchange type
func NewFanoutProducer(exchangeName string, config config.RabbitMQConfig) FanoutProducer {
	producer, err := core.CreateProducer(exchangeName, amqp091.ExchangeFanout, config)
	if err != nil {
		log.StdFatal(context.Background(), config, err, "error on creating fanout producer")
	}

	return FanoutProducer{
		coreProducer: producer,
	}
}

// Publish publish event to queue
func (p FanoutProducer) Publish(event proto.Message, routingKey string) error {
	body, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	return p.coreProducer.Publish(body, routingKey)
}
