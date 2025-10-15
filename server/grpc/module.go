package grpc

import (
	"github.com/go-playground/validator/v10"
	gophkeeperv1 "github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeper"
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
		RegisterValidationRules,
	),
)

func StartServer(*Server) {
}

func RegisterValidationRules(validate *validator.Validate) {
	validate.RegisterStructValidationMapRules(map[string]string{
		"login":    "required",
		"password": "required",
	}, gophkeeperv1.UserCredentials{})
}
