package tui

import (
	"github.com/mkolibaba/gophkeeper/internal/client/tui/view"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"tui",
	fx.Provide(
		view.InitialAuthorizationViewModel,
		view.InitialMainViewModel,
		view.InitialAddDataViewModel,
		NewBubble,
	),
	fx.Invoke(
		Start,
	),
)
