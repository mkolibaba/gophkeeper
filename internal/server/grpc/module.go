package grpc

import "go.uber.org/fx"

var Module = fx.Module(
	"grpc",
	fx.Provide(
		NewConfig,
		NewDataServiceServer,
		NewServer,
	),
	fx.Invoke(
		func(*Server) {},
	),
)
