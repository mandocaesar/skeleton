package sample

import (
	"time"

	"github.com/machtwatch/catalyst-go-skeleton/domain/common/model"
	"github.com/machtwatch/catalyst-go-skeleton/domain/common/repository"
	request "github.com/machtwatch/catalyst-go-skeleton/domain/common/request/api"
)

// TableName used to map table name with entity model name
func (e SampleEntity) TableName() string {
	return repository.TableSample
}

// SampleEntity represents a sample database entity struct.
// It includes of entity BaseModel.
type SampleEntity struct {
	model.BaseModel
	Reference     *string   `json:"reference"`
	ShippingFee   float64   `json:"shipping_fee"`
	InsuranceFee  float64   `json:"insurance_fee"`
	AdjustmentFee float64   `json:"adjustment_fee"`
	TotalPrice    float64   `json:"total_price"`
	TotalPriceOld float64   `json:"total_price_old"`
	OrderDate     time.Time `json:"order_date"`
	OrderID       int64     `json:"order_id"`
	Note          string    `json:"note"`
}

// SampleJoinModel represents a sample entity struct join model
type SampleJoinModel struct {
	model.BaseModel `json:",inline"`
	OrderNumber     string `json:"order_number" gorm:"column:xms_order_number"`
	OrderReference  string `json:"order_reference" gorm:"column:reference"`
	Source          string `json:"source" gorm:"column:source"`
	Note            string `json:"note" gorm:"column:note"`
}

// SampleListRequest represents a sample entity struct list request
type SampleListRequest struct {
	Pagination request.PaginationRequest `json:"pagination"`
	Filter     SampleFilterRequest       `json:"filter"`
	Search     string                    `json:"search"`
}

// SampleFilterRequest represents a sample entity filter request
type SampleFilterRequest struct {
	SampleIDs []int64 `json:"order_ids"`
	PriceGap  int8    `json:"price_gap"` // 0: all, 1: ada selisih, 2: tidak ada selisih
	DateFrom  string  `json:"date_from"`
	DateTo    string  `json:"date_to"`
}

// CreateOrderRequest represents a sample entity as params for creating order request
type CreateOrderRequest struct {
	ResellerID     int64   `json:"reseller_id"`
	OfficeID       int64   `json:"office_id"`
	Adjustment     float64 `json:"adjustment"`
	TotalPrice     float64 `json:"total_price"`
	Notes          string  `json:"notes"`
	IsReserveOrder bool    `json:"is_reserve_order"`
}

// CreateScheduleRequest represents a sample entity as params for creating schedule request
// format RunAt yyyy-mm-dd hh:mm:ss
type CreateScheduleRequest struct {
	RunAt   string `json:"run_at"`
	Message string `json:"message"`
}

// PublishEventRequest represents a sample entity as params for creating order request
type PublishEventRequest struct {
	ID             int64   `json:"id"`
	OrderId        int64   `json:"order_id"`
	Adjustment     float64 `json:"adjustment"`
	TotalPrice     float64 `json:"total_price"`
	Status         string  `json:"status"`
	IsReserveOrder bool    `json:"is_reserve_order"`
}
