package middlewares

import (
	"net/http"

	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/cache"
)

// APIMiddleware instantiate middleware by the given params and returns a new instance of middleware
type APIMiddleware interface {
	// Authentication authenticates requests using the jwt token
	Authentication(next http.Handler) http.Handler
	// AuthenticateScheduler authenticates requests from scheduler
	AuthenticateScheduler(next http.Handler) http.Handler
}

type ApiMiddleware struct {
	cache cache.Cache
}

func NewApiMiddleware(cache cache.Cache) APIMiddleware {
	return &ApiMiddleware{cache: cache}
}
