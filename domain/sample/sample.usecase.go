package sample

import (
	"context"

	"github.com/machtwatch/catalyst-go-skeleton/domain/common/response"
)

type SampleUC interface {
	CreateOrder(ctx context.Context, req CreateOrderRequest, token string) response.Response[bool]
	CreateSchedule(ctx context.Context, req CreateScheduleRequest) (res response.Response[interface{}])
	ReceiveMessageByScheduler(ctx context.Context, payload map[string]interface{}) (res response.Response[interface{}])
}
