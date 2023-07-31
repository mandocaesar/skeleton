package messaging

import (
	"github.com/machtwatch/catalyst-go-skeleton/app/sample/delivery/messaging"
	"github.com/machtwatch/catalyst-go-skeleton/domain"
	"github.com/machtwatch/catalyst-go-skeleton/domain/common/event"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging/rabbitmq/direct_factory"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging/rabbitmq/fanout_factory"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging/rabbitmq/topics_factory"
	"github.com/machtwatch/catalyst-go-skeleton/presentation/messaging/requeuer"
	"github.com/machtwatch/catalystdk/go/log"
	"github.com/rabbitmq/amqp091-go"
)

func Set(usecaseCollection domain.UsecaseCollection) {
	var (
		dlxRequeuer = requeuer.NewDLXRequeuer()
		baseConfig  = config.GetRmqConfig()
		topicConfig = config.GetRmqConfig(amqp091.Table{"x-dead-letter-exchange": event.DOMAIN_EVENT_DLX})
		dlxConfig   = config.GetRmqConfig(amqp091.Table{
			"x-dead-letter-exchange": event.DOMAIN_EVENT_EXCHANGE,
			"x-message-ttl":          config.EVENT_REQUEUE_DELAY_MS,
		})
	)

	// Create dead letter queue for rejected events
	fanout_factory.NewFanoutConsumer(event.DOMAIN_EVENT_DLX, event.DOMAIN_EVENT_DLX_QUEUE, dlxConfig)

	// Create parking lot queue for unhandled events
	direct_factory.NewDirectConsumer(event.PARKING_LOT_EXCHANGE, event.PARKING_LOT_QUEUE, baseConfig)

	sampleUpdatedConfig := topics_factory.TopicsConsumerConfig[*event.SampleUpdated]{
		Queue:        "sample-updated-queue",
		Topic:        event.SAMPLE_UPDATED_TOPIC,
		ExchangeName: event.DOMAIN_EVENT_EXCHANGE,
		Handler:      messaging.NewSampleUpdatedHandler(usecaseCollection.SampleUC),
		Requeuer:     dlxRequeuer,
		RmqConfig:    topicConfig,
	}
	err := topics_factory.NewTopicsConsumer(sampleUpdatedConfig).Listen()
	if err != nil {
		log.Fatalf("error on creating event listener: %v", err)
	}
}
