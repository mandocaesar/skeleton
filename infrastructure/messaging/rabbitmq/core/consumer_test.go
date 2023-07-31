//go:build rabbitmq
// +build rabbitmq

// can only be run with -tags=rabbitmq
// example: go test ./... -tags=rabbitmq
package core

import (
	"testing"

	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging/common/config"
	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

func TestConsumerCanListen(t *testing.T) {
	_, err := CreateConsumer(
		"direct_exchange",
		amqp091.ExchangeDirect,
		"direct_queue",
		config.RabbitMQConfig{
			AmqpURI:     "amqp://guest:guest@localhost:5672/",
			Durable:     true,
			Exclusive:   false,
			AutoDeleted: false,
			Internal:    false,
			NoWait:      false,
			Arguments:   nil,
		},
	)
	assert.Equal(t, nil, err)
}
