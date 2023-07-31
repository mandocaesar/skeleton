package eventstore

import (
	"context"
	"encoding/json"

	"github.com/machtwatch/catalyst-go-skeleton/domain/common/model"
	"github.com/machtwatch/catalystdk/go/log"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

type PgEventStore struct {
	db *gorm.DB
}

type EventStoreData struct {
	model.EventStoreModel
	Data  string `json:"data"`
	Error string `json:"error"`
}

func NewPgEventStore(db *gorm.DB) *PgEventStore {
	return &PgEventStore{
		db: db,
	}
}

// SaveEvent saves a starting (i.e : status = "Attempt") event to the event store.
func (p *PgEventStore) SaveEvent(ctx context.Context, event proto.Message, routingKey string) error {
	data, err := json.Marshal(event)
	if err != nil {
		log.StdDebug(ctx, event, err, "json.Marshal error on marshal sample event - uc.PublishEvent()")
		return err
	}

	// p.db.AutoMigrate(&EventStoreData{})
	eventData := EventStoreData{
		Data: string(data),
		EventStoreModel: model.EventStoreModel{
			Status: "Attempt",
		},
	}

	if err := p.db.WithContext(ctx).Create(&eventData).Error; err != nil {
		log.StdDebug(ctx, map[string]interface{}{"estore": event}, err, "r.writeConn.WithContext() got error - pgEventStore.SaveEvent")
		return err
	}

	return nil
}

// SaveEvent saves a failed event along with the cause to the event store.
func (p *PgEventStore) SaveFailedEvent(ctx context.Context, event proto.Message, routingKey string, errString string) error {
	data, err := json.Marshal(event)
	if err != nil {
		log.StdDebug(ctx, map[string]interface{}{"estore": event}, err, "json.Marshal got error - pgEventStore.SaveFailedEvent")
		return err
	}

	eventData := EventStoreData{
		Data:  string(data),
		Error: errString,
		EventStoreModel: model.EventStoreModel{
			Status: "Failed",
		},
	}

	if err := p.db.WithContext(ctx).Create(&eventData).Error; err != nil {
		log.StdDebug(ctx, map[string]interface{}{"estore": event}, err, "r.writeConn.WithContext() got error - pgEventStore.SaveFailedEvent")
		return err
	}

	return err
}
