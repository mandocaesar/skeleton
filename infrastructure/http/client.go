package http

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config"
)

// DefaultHTTPClient will return a new http client as the project's default http client
// This default http client is configured using env variables
func DefaultHTTPClient() *resty.Client {
	return resty.New().
		SetRetryCount(config.HTTP_REQUEST_RETRY_COUNT).
		SetRetryWaitTime(time.Duration(config.HTTP_REQUEST_RETRY_WAIT_TIME_MS) * time.Millisecond).
		SetRetryMaxWaitTime(time.Duration(config.HTTP_REQUEST_MAX_WAIT_TIME_MS) * time.Millisecond)
}
