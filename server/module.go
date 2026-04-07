package server

import (
	"os"

	"charm.land/log/v2"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"server",
	fx.Provide(
		NewLogger,
		NewValidate,
		NewConfig,
	),
	fx.Invoke(
		RegisterDataValidationRules,
	),
)

func NewLogger(config *Config) *log.Logger {
	opts := log.Options{
		ReportTimestamp: true,
	}

	if config.Development.Enabled {
		opts.Level = log.DebugLevel
	}

	return log.NewWithOptions(os.Stderr, opts)
}

func NewValidate() *validator.Validate {
	return validator.New()
}
