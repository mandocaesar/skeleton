package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HealthcheckHandler struct {
}

// NewHandler instantiate Health Check handler.
//
// It returns the instance of health check handler.
func NewHealthCheckHandler(r *chi.Mux) {
	handler := &HealthcheckHandler{}

	r.Route("/healthcheck", func(r chi.Router) {
		r.Get("/", handler.HandleHealthCheck)
	})
}

// Handle write http.StatusOK to response of successful health check
func (h *HealthcheckHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
