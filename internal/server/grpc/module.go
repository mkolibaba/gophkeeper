package grpc

import "go.uber.org/fx"

var Module = fx.Module(
	"grpc",
	fx.Provide(
		NewConfig,
		NewAuthorizationService,
		NewLoginServiceServer,
		NewNoteServiceServer,
		NewBinaryServiceServer,
		NewCardServiceServer,
		NewServer,
	),
	fx.Invoke(
		StartServer,
	),
)

func StartServer(*Server) {
}
