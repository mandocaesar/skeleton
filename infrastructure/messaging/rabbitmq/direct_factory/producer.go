package direct_factory

import (
	"context"

	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging/common/config"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging/rabbitmq/core"
	"github.com/machtwatch/catalystdk/go/log"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type DirectProducer struct {
	coreProducer *core.Producer
}

// NewDirectProducer creates messages producer with direct exchange type
func NewDirectProducer(exchangeName string, config config.RabbitMQConfig) DirectProducer {
	producer, err := core.CreateProducer(exchangeName, amqp091.ExchangeDirect, config)
	if err != nil {
		log.StdFatal(context.Background(), config, err, "error on creating direct producer")
	}

	return DirectProducer{
		coreProducer: producer,
	}
}

// Publish publish event to queue
func (p DirectProducer) Publish(ctx context.Context, event proto.Message, routingKey string) error {
	body, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	return p.coreProducer.Publish(body, routingKey)
}
