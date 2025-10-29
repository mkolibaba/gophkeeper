package main

import (
	"github.com/charmbracelet/log"
	"github.com/mkolibaba/gophkeeper/server"
	"github.com/mkolibaba/gophkeeper/server/grpc"
	"github.com/mkolibaba/gophkeeper/server/jwt"
	"github.com/mkolibaba/gophkeeper/server/sqlite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"log/slog"
)

func main() {
	fx.New(createApp()).Run()
}

func createApp() fx.Option {
	return fx.Options(
		fx.WithLogger(func(logger *log.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{Logger: slog.New(logger)}
		}),
		server.Module,
		sqlite.Module,
		grpc.Module,
		fx.Provide(
			fx.Annotate(jwt.NewAuthorizationService, fx.As(new(server.AuthorizationService))),
		),
	)
}
