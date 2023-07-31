package direct_factory

import (
	"context"

	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging/common/config"
	core "github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging/rabbitmq/core"
	"github.com/machtwatch/catalystdk/go/log"
	"github.com/rabbitmq/amqp091-go"
)

type DirectConsumer struct {
	coreConsumer *core.Consumer
}

func NewDirectConsumer(exchangeName, queueName string, config config.RabbitMQConfig) DirectConsumer {
	consumer, err := core.CreateConsumer(
		exchangeName,
		amqp091.ExchangeDirect,
		queueName,
		config,
	)
	if err != nil {
		log.StdFatal(context.Background(), config, err, "error on creating direct consumer")
	}

	return DirectConsumer{
		coreConsumer: consumer,
	}
}
