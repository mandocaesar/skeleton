package middlewares

import (
	"context"
	"net/http"
)

var AuthKey struct{}

func SetHTTPHeadersToContext() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get("Authorization")
			ctx := context.WithValue(r.Context(), AuthKey, authorizationHeader)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
