package response

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config"
	"github.com/machtwatch/catalystdk/go/log"
	"github.com/machtwatch/catalystdk/go/trace"
	"go.opentelemetry.io/otel/codes"
)

type MetaPagination struct {
	Page      int    `json:"page"`
	Size      int    `json:"size"`
	SortBy    string `json:"sort_by,omitempty"`
	SortType  string `json:"sort_type,omitempty"`
	TotalData int    `json:"total_data"`
	TotalPage int    `json:"total_pages"`
}

type MetaResponse struct {
	Code       string          `json:"code"`
	Message    string          `json:"message"`
	Error      string          `json:"error,omitempty"`
	Pagination *MetaPagination `json:"pagination,omitempty"`
}

type Response[T any] struct {
	RequestID  string         `json:"request_id"`
	Code       string         `json:"code"`
	Message    string         `json:"message"`
	Errors     *[]interface{} `json:"errors"`
	ServerTime int64          `json:"server_time"`
	Data       *T             `json:"data"`
}

const (
	SUCCESS = "SUCCESS"
)

type IntegrationServiceResponse struct {
	Code       string        `json:"code"`
	Data       interface{}   `json:"data"`
	Errors     []interface{} `json:"errors"`
	Message    string        `json:"message"`
	ServerTime int           `json:"server_time"`
}

var HttpStatus = map[string]int{
	CodeSuccess:             http.StatusOK,
	CodeInternalServerError: http.StatusInternalServerError,
	CodeBadRequest:          http.StatusBadRequest,
	CodeNotFound:            http.StatusNotFound,
	CodeUnauthorized:        http.StatusUnauthorized,
}

func (res *Response[T]) Success(ctx context.Context, data T, msgs ...string) Response[T] {

	var msgStr string = "Success"
	if len(msgs) > 0 {
		msgStr = msgs[0]
	}

	span := trace.GetSpanFromContext(ctx)
	if span != nil && span.IsRecording() {
		span.SetStatus(codes.Ok, msgStr)
	}

	return Response[T]{
		RequestID:  log.GetCtxRequestID(ctx),
		Code:       CodeSuccess,
		Message:    msgStr,
		Data:       &data,
		ServerTime: time.Now().Unix(),
	}
}

func (res *Response[T]) InternalError(ctx context.Context, err error, msgStr string) Response[T] {
	var errStr string
	if config.DEBUG {
		errStr = err.Error()
	}

	trace.SetErrorSpanFromContext(ctx, err, errStr)

	return Response[T]{
		RequestID:  log.GetCtxRequestID(ctx),
		Code:       CodeInternalServerError,
		Message:    msgStr,
		Errors:     &[]interface{}{errStr},
		ServerTime: time.Now().Unix(),
	}
}

func (res *Response[T]) InvalidPayload(ctx context.Context, msgStr string) Response[T] {

	err := errors.New(msgStr)
	trace.SetErrorSpanFromContext(ctx, err, msgStr)

	return Response[T]{
		RequestID:  log.GetCtxRequestID(ctx),
		Code:       CodeBadRequest,
		Message:    msgStr,
		ServerTime: time.Now().Unix(),
	}

}

func (res *Response[T]) NotFound(ctx context.Context, msgStr string) Response[T] {

	err := errors.New(msgStr)
	trace.SetErrorSpanFromContext(ctx, err, msgStr)

	return Response[T]{
		RequestID:  log.GetCtxRequestID(ctx),
		Code:       CodeNotFound,
		Message:    msgStr,
		ServerTime: time.Now().Unix(),
	}
}

func (res *Response[T]) Unauthorized(ctx context.Context, msgStr string) Response[T] {

	err := errors.New(msgStr)
	trace.SetErrorSpanFromContext(ctx, err, msgStr)

	return Response[T]{
		RequestID:  log.GetCtxRequestID(ctx),
		Code:       CodeUnauthorized,
		Message:    msgStr,
		ServerTime: time.Now().Unix(),
	}
}

func (res *Response[T]) WriteResponse(w http.ResponseWriter) {
	resCode := HttpStatus[res.Code]
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Status-Code", strconv.Itoa(resCode))
	w.WriteHeader(resCode)

	json.NewEncoder(w).Encode(res)
}
