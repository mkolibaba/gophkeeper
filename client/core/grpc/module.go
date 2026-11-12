package grpc

import (
	"github.com/mkolibaba/gophkeeper/client/core"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"grpc",
	fx.Provide(
		NewConnection,
		NewAuthorizationServiceClient,
		fx.Annotate(NewAuthorizationService, fx.As(new(core.AuthorizationService))),
		NewLoginServiceClient,
		fx.Annotate(NewLoginService, fx.As(new(core.LoginService))),
		NewNoteServiceClient,
		fx.Annotate(NewNoteService, fx.As(new(core.NoteService))),
		NewBinaryServiceClient,
		fx.Annotate(NewBinaryService, fx.As(new(core.BinaryService))),
		NewCardServiceClient,
		fx.Annotate(NewCardService, fx.As(new(core.CardService))),
	),
)
