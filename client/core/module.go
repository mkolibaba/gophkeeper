package core

import (
	"fmt"
	"os"

	"charm.land/log/v2"
)

func NewLogger(config *Config) (*log.Logger, error) {
	logOutput := config.Log.Output
	if logOutput == "" {
		return nil, fmt.Errorf("log output is not set")
	}

	flags := os.O_WRONLY | os.O_CREATE
	if config.Log.Truncate {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_EXCL
	}

	out, err := os.OpenFile(logOutput, flags, 0666)
	if err != nil {
		return nil, fmt.Errorf("new logger: %w", err)
	}

	opts := log.Options{
		ReportTimestamp: true,
		Formatter:       log.JSONFormatter,
	}
	if config.Development.Enabled {
		opts.Level = log.DebugLevel
	}

	return log.NewWithOptions(out, opts), nil
}

func PrintConfig(logger *log.Logger, config *Config) {
	logger.Debug("client config", "config", config)
}
