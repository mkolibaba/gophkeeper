package inmem

import (
	"github.com/mkolibaba/gophkeeper/client/core"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"inmem",
	fx.Provide(
		fx.Annotate(NewUserService, fx.As(new(core.UserService))),
	),
)
