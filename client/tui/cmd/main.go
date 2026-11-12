package main

import (
	"log/slog"

	"charm.land/log/v2"
	core "github.com/mkolibaba/gophkeeper/client/core/module"
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
		core.Module,
		tui.Module,
	)
}
