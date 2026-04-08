package main

import (
	"log/slog"

	"charm.land/log/v2"
	"github.com/mkolibaba/gophkeeper/client"
	"github.com/mkolibaba/gophkeeper/client/grpc"
	"github.com/mkolibaba/gophkeeper/client/inmem"
	"github.com/mkolibaba/gophkeeper/client/tui"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {
	fx.New(createApp()).Run()
}

func createApp() fx.Option {
	return fx.Options(
		fx.WithLogger(func(logger *log.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{Logger: slog.New(logger)}
		}),
		client.Module,
		grpc.Module,
		inmem.Module,
		tui.Module,
	)
}
