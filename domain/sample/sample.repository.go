package sample

import (
	"context"

	"github.com/machtwatch/catalyst-go-skeleton/domain/common/repository"
	"github.com/machtwatch/catalyst-go-skeleton/domain/common/response"
)

type SampleRepo interface {
	repository.IBaseRepo[SampleEntity]
	GetSampleRequest(ctx context.Context, req SampleListRequest) (samples []SampleJoinModel, pagination response.MetaPagination, err error)
}
