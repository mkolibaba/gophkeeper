package main

import (
	"github.com/mkolibaba/gophkeeper/internal/server"
	"github.com/mkolibaba/gophkeeper/internal/server/grpc"
	"github.com/mkolibaba/gophkeeper/internal/server/sqlite"
	"go.uber.org/fx"
)

func main() {
	fx.New(createApp()).Run()
}

func createApp() fx.Option {
	return fx.Options(
		server.Module,
		sqlite.Module,
		grpc.Module,
	)
}
