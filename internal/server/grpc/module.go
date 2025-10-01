package grpc

import "go.uber.org/fx"

var Module = fx.Module(
	"grpc",
	fx.Provide(
		NewConfig,
		NewDataServiceServer,
		NewBinaryServiceServer,
		NewServer,
	),
	fx.Invoke(
		StartServer,
	),
)

func StartServer(*Server) {
}
