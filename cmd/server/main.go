package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/mkolibaba/gophkeeper/internal/server/grpc"
	"github.com/mkolibaba/gophkeeper/internal/server/sqlite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
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
			zap.NewDevelopment,
			validator.New,
		),
		sqlite.Module,
		grpc.Module,
	)
}
