package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"github.com/machtwatch/catalyst-go-skeleton/domain/common/response"
	"github.com/machtwatch/catalyst-go-skeleton/domain/sample"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/tracer"
	"github.com/machtwatch/catalyst-go-skeleton/presentation/middlewares"
	"github.com/machtwatch/catalystdk/go/trace"
)

const (
	segmentHandleCreateOrder     = tracer.SegmentHandler + "SampleHandler.CreateOrder"
	segmentHandleCreateScheduler = tracer.SegmentHandler + "SampleHandler.CreateScheduler"
)

// SampleHandler represent sample handler that provide list of use case needed.
type SampleHandler struct {
	usecase sample.SampleUC
}

// NewSampleHandler instantiate Sample handler and register the correspond http routes.
//
// It Accept the sample usecase and will return the instance of sample handler.
func NewSampleHandler(r *chi.Mux, usecase sample.SampleUC, apiMiddlewares middlewares.APIMiddleware) {
	handler := &SampleHandler{
		usecase,
	}
	r.Route("/api/sample", func(r chi.Router) {
		r.Get("/get", handler.CreateOrder)
		// Post Schedule is a sample for triggering schedule
		r.Post("/schedule", handler.CreateSchedule)
	})
}

// CreateOrder create a sample order by the given http request parameter.
//
// It will write success http response if success. Otherwise it will return error
// http response.
func (h *SampleHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var (
		model sample.CreateOrderRequest
		res   response.Response[bool]
	)

	ctx := r.Context()
	ctx, span := trace.StartTransactionFromContext(ctx, segmentHandleCreateOrder)
	defer span.End()

	var token string = r.Header.Get("authorization")
	if err := render.DecodeJSON(r.Body, &model); err != nil {
		res = res.InternalError(ctx, err, "Invalid request body")
		res.WriteResponse(w)
		return
	}

	validate := validator.New()
	if err := validate.Struct(model); err != nil {
		res = res.InvalidPayload(ctx, fmt.Sprintf("Invalid request body: %v", err.Error()))
		res.WriteResponse(w)
		return
	}

	res = h.usecase.CreateOrder(ctx, model, token)
	res.WriteResponse(w)
}

// CreateSchedule create a sample scheduler to rundeck by the given http request parameter.
//
// It will write success http response if success. Otherwise it will return error
// http response.
func (h *SampleHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	var (
		model sample.CreateScheduleRequest
		res   response.Response[interface{}]
	)

	ctx := r.Context()
	ctx, span := trace.StartTransactionFromContext(ctx, segmentHandleCreateScheduler)
	defer span.End()

	if err := render.DecodeJSON(r.Body, &model); err != nil {
		res = res.InternalError(ctx, err, "Invalid request body")
		res.WriteResponse(w)
		return
	}

	res = h.usecase.CreateSchedule(ctx, model)
	res.WriteResponse(w)
}
