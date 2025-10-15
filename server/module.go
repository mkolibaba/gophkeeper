package server

import (
	"github.com/charmbracelet/log"
	"go.uber.org/fx"
	"os"
)

var Module = fx.Module(
	"server",
	fx.Provide(
		newLogger,
		NewDataValidator,
		NewConfig,
	),
)

func newLogger() *log.Logger {
	return log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
	})
}
