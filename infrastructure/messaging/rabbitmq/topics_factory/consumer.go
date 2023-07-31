package topics_factory

import (
	"context"
	"reflect"

	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging/common/config"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging/rabbitmq/core"
	"github.com/machtwatch/catalystdk/go/log"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type TopicsConsumer[T proto.Message] struct {
	coreConsumer *core.Consumer
	queue        string
	topic        string
	handler      messaging.Handler[T]
	requeuer     messaging.Requeuer
}

type TopicsConsumerConfig[T proto.Message] struct {
	Queue        string
	Topic        string
	ExchangeName string
	Handler      messaging.Handler[T]
	Requeuer     messaging.Requeuer
	RmqConfig    config.RabbitMQConfig
}

// NewTopicsConsumer creates rabbitmq consumer with topic exhchange type
func NewTopicsConsumer[T proto.Message](consumerConfig TopicsConsumerConfig[T]) TopicsConsumer[T] {
	consumer, err := core.CreateConsumerWithRouting(
		consumerConfig.ExchangeName,
		amqp091.ExchangeTopic,
		consumerConfig.Queue,
		consumerConfig.Topic,
		consumerConfig.RmqConfig,
	)

	if err != nil {
		log.StdFatal(context.Background(), consumerConfig, err, "error on creating topic consumer")
	}

	return TopicsConsumer[T]{
		coreConsumer: consumer,
		topic:        consumerConfig.Topic,
		queue:        consumerConfig.Queue,
		handler:      consumerConfig.Handler,
		requeuer:     consumerConfig.Requeuer,
	}
}

// Listen listen for incoming messages from the exchange
func (c TopicsConsumer[T]) Listen() error {
	return c.coreConsumer.Listen(func(delivery amqp091.Delivery, err error) {
		var ctx = context.Background()
		var msg T

		// Obtain the type inside T
		msgType := reflect.TypeOf(msg).Elem()

		// Create non nill pointer to zero value object of the type, and transform it back into T
		msg = reflect.New(msgType).Interface().(T)

		if err := proto.Unmarshal(delivery.Body, msg); err != nil {
			log.StdError(ctx, delivery, err, "error unmarshalling event payload")
			delivery.Reject(false)
			return
		}

		if err := c.handler.Handle(msg); err != nil {
			log.StdError(ctx, msg, err, "error on handling event")

			if c.requeuer != nil {
				c.requeuer.Handle(delivery, msg, err)
				return
			}
		}

		delivery.Ack(false)
	})
}
