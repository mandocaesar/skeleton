package topics_factory

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config"
	rabbitmqCfg "github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging/common/config"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/messaging/rabbitmq/core"
	"github.com/machtwatch/catalystdk/go/log"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type TopicsProducer struct {
	coreProducer *core.Producer
}

// NewTopicsProducer creates messages producer with topic exchange type
func NewTopicsProducer(exchangeName string, rabbitmqCfg rabbitmqCfg.RabbitMQConfig) TopicsProducer {
	producer, err := core.CreateProducer(exchangeName, amqp091.ExchangeTopic, rabbitmqCfg)
	if err != nil {
		log.StdFatal(context.Background(), rabbitmqCfg, err, "error on creating topic producer")
	}

	return TopicsProducer{
		coreProducer: producer,
	}
}

// Publish publish event to queue
func (p TopicsProducer) Publish(ctx context.Context, event proto.Message, routingKey string) error {
	body, err := proto.Marshal(event)
	if err != nil {
		log.StdDebug(ctx, map[string]interface{}{"event": event, "routingKey": routingKey}, err, "proto.Marshal failed - (p TopicsProducer) Publish")
		return err
	}

	err = p.coreProducer.Publish(body, routingKey)
	if err != nil {
		log.StdDebug(ctx, map[string]interface{}{"body": body, "routingKey": routingKey}, err, "p.coreProducer.Publish failed, retrying.. - (p TopicsProducer) Publish")
		err = p.Republish(ctx, body, routingKey)
		if err != nil {
			log.StdDebug(ctx, map[string]interface{}{"body": body, "routingKey": routingKey}, err, "p.Republish failed - (p TopicsProducer) Publish")
			return err
		}
	}

	return err
}

func (p TopicsProducer) Republish(ctx context.Context, body []byte, routingKey string) error {
	var (
		err                 error
		maxRetry            = uint64(config.PUBLISH_MAX_RETRY_COUNT)
		backoffFactor       = float64(config.PUBLISH_BACKOFF_FACTOR)
		delayDuration       = time.Duration(config.PUBLISH_RETRY_WAIT_TIME_SEC) * time.Second
		maxRetryElapsedTime = time.Duration(config.PUBLISH_MAX_RETRY_ELAPSED_TIME_SEC) * time.Second
		maxDelayDuration    = time.Duration(config.PUBLISH_MAX_RETRY_WAIT_TIME_SEC)
	)

	// this config should be generated at runtime.
	expBackoffCfg := backoff.NewExponentialBackOff()
	expBackoffCfg.InitialInterval = delayDuration
	expBackoffCfg.MaxInterval = maxDelayDuration
	expBackoffCfg.MaxElapsedTime = maxRetryElapsedTime
	expBackoffCfg.Multiplier = backoffFactor

	backoffCfg := backoff.WithMaxRetries(backoff.WithContext(expBackoffCfg, ctx), maxRetry)

	backoff.Retry(func() error {
		err = p.coreProducer.Publish(body, routingKey)
		if err != nil {
			log.StdDebugf(ctx, map[string]interface{}{
				"body": body, "routingKey": routingKey}, err, "error on republishing event, retrying...")
		}

		return err
	}, backoffCfg)

	if err != nil {
		log.StdDebug(ctx, map[string]interface{}{"body": body, "routingKey": routingKey}, err, "error on republishing event, max retry count/time achieved")
	}

	return err
}
