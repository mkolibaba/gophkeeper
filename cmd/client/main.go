package main

import (
	"github.com/mkolibaba/gophkeeper/internal/client"
	"github.com/mkolibaba/gophkeeper/internal/client/grpc"
	"github.com/mkolibaba/gophkeeper/internal/client/inmem"
	"github.com/mkolibaba/gophkeeper/internal/client/tui"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
)

func main() {
	fx.New(createApp()).Run()
}

func createApp() fx.Option {
	return fx.Options(
		fx.Provide(
			newLogger,
		),
		client.Module,
		grpc.Module,
		inmem.Module,
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
