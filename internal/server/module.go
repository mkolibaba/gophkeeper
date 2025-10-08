package server

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"server",
	fx.Provide(
		zap.NewDevelopment,
		NewDataValidator,
		NewAuthService,
	),
)
