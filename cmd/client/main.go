package main

import (
	"github.com/charmbracelet/log"
	"github.com/mkolibaba/gophkeeper/internal/client"
	"github.com/mkolibaba/gophkeeper/internal/client/grpc"
	"github.com/mkolibaba/gophkeeper/internal/client/inmem"
	"github.com/mkolibaba/gophkeeper/internal/client/tui"
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
		client.Module,
		grpc.Module,
		inmem.Module,
		tui.Module,
	)
}
