package grpc

import (
	"github.com/mkolibaba/gophkeeper/server/grpc/interceptors"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"grpc",
	fx.Provide(
		interceptors.NewAuthInterceptor,
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
