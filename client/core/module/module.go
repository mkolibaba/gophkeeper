package module

import (
	"github.com/mkolibaba/gophkeeper/client/core"
	"github.com/mkolibaba/gophkeeper/client/core/grpc"
	"github.com/mkolibaba/gophkeeper/client/core/inmem"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"core",
	grpc.Module,
	inmem.Module,
	fx.Provide(
		core.NewLogger,
		core.NewDataValidator,
		core.NewConfig,
	),
	fx.Invoke(
		core.PrintConfig,
	),
)
