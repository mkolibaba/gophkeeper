package grpc

import (
	"github.com/mkolibaba/gophkeeper/internal/client"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"grpc",
	fx.Provide(
		NewConfig,
		NewConnection,
		fx.Annotate(NewAuthorizationService, fx.As(new(client.AuthorizationService))),
		fx.Annotate(NewLoginService, fx.As(new(client.LoginService))),
		fx.Annotate(NewNoteService, fx.As(new(client.NoteService))),
		fx.Annotate(NewBinaryService, fx.As(new(client.BinaryService))),
		fx.Annotate(NewCardService, fx.As(new(client.CardService))),
	),
)
