package messaging

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type Producer interface {
	Publish(ctx context.Context, event proto.Message, routingKey string) error
}

type Consumer interface {
	Listen() error
}

type Handler[T proto.Message] interface {
	Handle(T) error
}

type Requeuer interface {
	// Handle must ack/nack/reject delivery
	Handle(delivery amqp091.Delivery, msg proto.Message, err error) error
}
