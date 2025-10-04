package main

import (
	"github.com/mkolibaba/gophkeeper/internal/client"
	"github.com/mkolibaba/gophkeeper/internal/client/grpc"
	"github.com/mkolibaba/gophkeeper/internal/client/tui"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"os"
)

func main() {
	fx.New(createApp()).Run()
}

func createApp() fx.Option {
	return fx.Options(
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		fx.Provide(
			newLogger,
			client.NewDataValidator,
		),
		grpc.Module,
		tui.Module,
	)
}

func newLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	if logOutput := os.Getenv("LOG_OUTPUT"); logOutput != "" {
		cfg.OutputPaths = []string{logOutput}
	}
	return cfg.Build()
}
