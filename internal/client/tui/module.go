package tui

import "go.uber.org/fx"

var Module = fx.Module(
	"tui",
	fx.Provide(
		NewBubble,
	),
	fx.Invoke(
		Start,
	),
)
