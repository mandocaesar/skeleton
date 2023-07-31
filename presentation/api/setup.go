package api

import (
	"github.com/go-chi/chi/v5"
	hcdeliveryhttp "github.com/machtwatch/catalyst-go-skeleton/app/healthcheck/delivery/http"
	sampledeliveryhttp "github.com/machtwatch/catalyst-go-skeleton/app/sample/delivery/http"
	"github.com/machtwatch/catalyst-go-skeleton/domain"
	"github.com/machtwatch/catalyst-go-skeleton/presentation/middlewares"
)

// Set instantiate all handlers with the application usecase and middlwares
func Set(r *chi.Mux, usecaseCollection domain.UsecaseCollection, apiMiddlwares middlewares.APIMiddleware) {
	sampledeliveryhttp.NewSampleHandler(r, usecaseCollection.SampleUC, apiMiddlwares)
	sampledeliveryhttp.NewSampleSchedulerHandler(r, usecaseCollection.SampleUC, apiMiddlwares)
	hcdeliveryhttp.NewHealthCheckHandler(r)
}
