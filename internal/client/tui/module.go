package tui

import (
	"github.com/mkolibaba/gophkeeper/internal/client/tui/orchestrator"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"tui",
	fx.Provide(
		orchestrator.New,
		NewBubble,
	),
	fx.Invoke(
		Start,
	),
)
