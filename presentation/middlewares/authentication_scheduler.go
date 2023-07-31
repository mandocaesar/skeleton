package middlewares

import (
	"net/http"
	"time"

	"github.com/machtwatch/catalyst-go-skeleton/domain/common/response"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config/secret"
)

// this authentication is for rundeck to access this service
// set this authenticate in routes that will be called by rundeck
func (a *ApiMiddleware) AuthenticateScheduler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		rundeckToken := r.Header.Get("X-Rundeck-Auth")
		if rundeckToken != secret.RUNDECK_JOB_BEARER_TOKEN {
			responseError(w, response.Response[bool]{
				Code:       response.CodeUnauthorized,
				Message:    "Unauthorized",
				ServerTime: time.Now().Unix(),
			}, response.HttpStatus[response.CodeUnauthorized])
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
