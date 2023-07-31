package fanout_factory

import (
	"context"

	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging/common/config"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging/rabbitmq/core"
	"github.com/machtwatch/catalystdk/go/log"
	"github.com/rabbitmq/amqp091-go"
)

type FanoutConsumer struct {
	coreConsumer *core.Consumer
}

func NewFanoutConsumer(exchangeName, queueName string, config config.RabbitMQConfig) FanoutConsumer {
	consumer, err := core.CreateConsumer(
		exchangeName,
		amqp091.ExchangeFanout,
		queueName,
		config,
	)
	if err != nil {
		log.StdFatal(context.Background(), config, err, "error on creating fanout consumer")
	}

	return FanoutConsumer{
		coreConsumer: consumer,
	}
}
