package grpc

import (
	"github.com/go-playground/validator/v10"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"
	"github.com/mkolibaba/gophkeeper/server/grpc/interceptors"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"grpc",
	fx.Provide(
		interceptors.NewAuthInterceptor,
		interceptors.NewLoggerInterceptor,
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
	validate.RegisterStructValidation(func(sl validator.StructLevel) {
		in := sl.Current().Interface().(gophkeeperv1.UserCredentials)
		if !in.HasLogin() || len(in.GetLogin()) == 0 {
			sl.ReportError(in.GetLogin(), "login", "login", "", "")
		}
		if !in.HasPassword() || len(in.GetPassword()) == 0 {
			sl.ReportError(in.GetPassword(), "password", "password", "", "")
		}
	}, gophkeeperv1.UserCredentials{})
}
