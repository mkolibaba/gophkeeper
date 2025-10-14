package grpc

import "go.uber.org/fx"

var Module = fx.Module(
	"grpc",
	fx.Provide(
		NewConfig,
		NewAuthorizationServiceServer,
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
