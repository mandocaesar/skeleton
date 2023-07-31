package logger

import (
	"testing"

	"github.com/machtwatch/catalystdk/go/log"
)

func TestSetStdLog(t *testing.T) {
	type args struct {
		appName string
		caller  bool
		level   string
		useJSON bool
	}

	tests := []struct {
		name     string
		args     args
		wantErr  error
		isActive bool
	}{
		{
			name: "when_params_not_empty",
			args: args{
				appName: "test-app",
				caller:  true,
				level:   "info",
				useJSON: true,
			},
			isActive: true,
		},
		{
			name: "when_app_name_not_empty",
			args: args{
				appName: "test-app",
			},
			isActive: true,
		},
		{
			name:     "when_params_nil",
			args:     args{},
			isActive: true,
		},
		{
			name: "when_appName_params_empty_return_error",
			args: args{
				appName: "",
			},
			isActive: true,
		},
	}

	for _, test := range tests {
		if !test.isActive {
			continue
		}
		t.Run(test.name, func(t *testing.T) {
			SetStandardLog(&log.Config{
				AppName: test.args.appName,
				Caller:  test.args.caller,
				Level:   test.args.level,
				UseJSON: test.args.useJSON,
			})
		})
	}
}
