package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/machtwatch/catalyst-go-skeleton/domain/common/response"
	"github.com/machtwatch/catalyst-go-skeleton/domain/sample"
	"github.com/machtwatch/catalyst-go-skeleton/presentation/middlewares"
	"github.com/machtwatch/catalystdk/go/trace"
)

// NewSampleSchedulerHandler instantiate Sample handler and register the correspond http routes.
//
// It Accept the sample usecase and will return the instance of sample handler.
// Post scheduler/sample is a sample for endpoint that called by rundeck
// and this endpoint using middleware authenticateScheduler
func NewSampleSchedulerHandler(r *chi.Mux, usecase sample.SampleUC, apiMiddlewares middlewares.APIMiddleware) {
	handler := &SampleHandler{
		usecase,
	}
	r.Route("/scheduler/sample", func(r chi.Router) {
		r.Use(apiMiddlewares.AuthenticateScheduler)
		r.Post("/", handler.ReceiveMessageByScheduler)
	})
}

// This is sample webhook scheduler. just receive what rundeck sent and compare the message
func (h *SampleHandler) ReceiveMessageByScheduler(w http.ResponseWriter, r *http.Request) {
	var (
		model map[string]interface{}
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

	res = h.usecase.ReceiveMessageByScheduler(ctx, model)
	res.WriteResponse(w)
}
