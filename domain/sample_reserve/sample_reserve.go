package samplereserve

import (
	"time"

	"github.com/machtwatch/catalyst-go-skeleton/domain/common/model"
	"github.com/machtwatch/catalyst-go-skeleton/domain/common/repository"
)

// TableName used to map table name with entity model name
func (e SampleReserveEntity) TableName() string {
	return repository.TableSampleReserve
}

type SampleReserveEntity struct {
	model.BaseModel
	OrderDate   time.Time
	OfficeID    int64 `json:"office_id"`
	OrderItemID int64 `json:"order_item_id"`
	Qty         int   `json:"qty"`
}

type SampleReserveResponse struct {
	SampleReserveEntity `json:",inline"`
	IsDropship          bool `json:"is_dropship"`
	IsNonDropship       bool `json:"is_non_dropship"`
}
