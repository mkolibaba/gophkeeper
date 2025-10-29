package tui

import (
	"github.com/mkolibaba/gophkeeper/client/tui/view/adddata"
	"github.com/mkolibaba/gophkeeper/client/tui/view/authorization"
	"github.com/mkolibaba/gophkeeper/client/tui/view/editdata"
	"github.com/mkolibaba/gophkeeper/client/tui/view/home"
	"github.com/mkolibaba/gophkeeper/client/tui/view/registration"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"tui",
	fx.Provide(
		authorization.New,
		home.New,
		adddata.New,
		registration.New,
		editdata.New,
		NewBubble,
	),
	fx.Invoke(
		Start,
	),
)
