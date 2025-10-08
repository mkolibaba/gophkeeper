package tui

import (
	"github.com/mkolibaba/gophkeeper/internal/client/tui/state"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"tui",
	fx.Provide(
		state.NewManager,
		NewBubble,
	),
	fx.Invoke(
		Start,
	),
)
