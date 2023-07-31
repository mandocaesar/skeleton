package logger

import (
	"context"

	"github.com/machtwatch/catalystdk/go/log"
)

// SetStandardLog init catalystdk log globally for standard logging with fields
func SetStandardLog(logConfig *log.Config) {
	if err := log.SetStdLog(logConfig); err != nil {
		// when got error on setting config, it will use the default config.
		log.StdError(context.Background(), nil, err, "init catalystdk log got error")
	}
}
