package apm

import (
	"github.com/machtwatch/catalystdk/go/log"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// NewApplication returns new newrelic Application.
func NewApplication(appName string, license string) *newrelic.Application {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(license),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)

	if err != nil {
		log.Fatalf("error on connecting to newrelic: %v", err)
	}

	return app
}
