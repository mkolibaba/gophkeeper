package tui

import (
	"github.com/mkolibaba/gophkeeper/internal/client/tui/view/adddata"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/view/authorization"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/view/home"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"tui",
	fx.Provide(
		authorization.New,
		home.New,
		adddata.New,
		NewBubble,
	),
	fx.Invoke(
		Start,
	),
)
