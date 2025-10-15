package server

import (
	"github.com/charmbracelet/log"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"os"
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

func NewLogger() *log.Logger {
	return log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
	})
}

func NewValidate() *validator.Validate {
	return validator.New()
}
