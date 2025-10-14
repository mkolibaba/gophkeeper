package client

import (
	"fmt"
	"github.com/charmbracelet/log"
	"go.uber.org/fx"
	"os"
	"strings"
)

var Module = fx.Module(
	"client",
	fx.Provide(
		newLogger,
		NewDataValidator,
	),
)

func newLogger() (*log.Logger, error) {
	logOutput := os.Getenv("LOG_OUTPUT")
	if logOutput == "" {
		return nil, fmt.Errorf("environment variable LOG_OUTPUT is not set")
	}

	flags := os.O_WRONLY | os.O_CREATE
	if strings.ToLower(os.Getenv("LOG_TRUNCATE")) == "true" {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_EXCL
	}

	out, err := os.OpenFile(logOutput, flags, 0666)
	if err != nil {
		return nil, fmt.Errorf("new logger: %w", err)
	}

	return log.NewWithOptions(out, log.Options{
		ReportTimestamp: true,
		Level:           log.DebugLevel,
		Formatter:       log.JSONFormatter,
	}), nil
}
