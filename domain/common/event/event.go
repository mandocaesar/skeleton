package event

type Event interface {
	Marshal() ([]byte, error)
}

const (
	// Exchanges
	DOMAIN_EVENT_EXCHANGE = "catalyst-go-skeleton.domain-event"
	DOMAIN_EVENT_DLX      = "catalyst-go-skeleton.domain-event.dlx"
	PARKING_LOT_EXCHANGE  = "catalyst-go-skeleton.parking-lot"

	// Topics
	SAMPLE_UPDATED_TOPIC = "xms.sample.updated"
	SAMPLE_TOPIC         = DOMAIN_EVENT_EXCHANGE + ".sample"

	// Queues
	SAMPLE_UPDATED_QUEUE   = "sample-updated-queue"
	DOMAIN_EVENT_DLX_QUEUE = DOMAIN_EVENT_DLX + "-queue"
	PARKING_LOT_QUEUE      = PARKING_LOT_EXCHANGE + "-queue"
)
