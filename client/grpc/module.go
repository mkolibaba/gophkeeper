package grpc

import (
	"github.com/mkolibaba/gophkeeper/client"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"grpc",
	fx.Provide(
		NewConnection,
		NewAuthorizationServiceClient,
		fx.Annotate(NewAuthorizationService, fx.As(new(client.AuthorizationService))),
		NewLoginServiceClient,
		fx.Annotate(NewLoginService, fx.As(new(client.LoginService))),
		NewNoteServiceClient,
		fx.Annotate(NewNoteService, fx.As(new(client.NoteService))),
		NewBinaryServiceClient,
		fx.Annotate(NewBinaryService, fx.As(new(client.BinaryService))),
		NewCardServiceClient,
		fx.Annotate(NewCardService, fx.As(new(client.CardService))),
	),
)
