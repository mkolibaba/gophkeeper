package main

import (
	"github.com/mkolibaba/gophkeeper/internal/client"
	"github.com/mkolibaba/gophkeeper/internal/client/grpc"
	"github.com/mkolibaba/gophkeeper/internal/client/inmem"
	"github.com/mkolibaba/gophkeeper/internal/client/tui"
	"go.uber.org/fx"
)

func main() {
	fx.New(createApp()).Run()
}

func createApp() fx.Option {
	return fx.Options(
		client.Module,
		grpc.Module,
		inmem.Module,
		tui.Module,
	)
}
