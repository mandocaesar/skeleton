package model

import "time"

// EventType is the type of any event, used as its unique identifier.
type EventType string

// AggregateType type of the Aggregate
type AggregateType string

type EventStoreEntityModel interface {
	TableName() string
}

type EventStoreModel struct {
	EventID       string        `gorm:"primarykey" json:"event_id"`
	AggregateID   string        `json:"aggregate_id"`
	EventType     EventType     `json:"event_type"`
	AggregateType AggregateType `json:"aggregate_typ"`
	Version       uint64        `json:"version"`
	Data          []byte        `json:"data"`
	Metadata      []byte        `json:"metadata"`
	Status        string        `json:"status"`
	Timestamp     time.Time     `json:"timestamp"`
}
