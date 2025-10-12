package inmem

import (
	"github.com/mkolibaba/gophkeeper/internal/client"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"inmem",
	fx.Provide(
		fx.Annotate(NewUserService, fx.As(new(client.UserService))),
	),
)
