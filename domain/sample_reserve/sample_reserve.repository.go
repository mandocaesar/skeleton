package samplereserve

import (
	"context"

	"github.com/machtwatch/catalyst-go-skeleton/domain/common/repository"
)

type SampleReserveRepo interface {
	repository.IBaseRepo[SampleReserveEntity]
	GetByOrderID(ctx context.Context, orderID int64, filter interface{}) ([]SampleReserveEntity, error)
}
