package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"go.uber.org/fx"
)

func Start(
	bubble Bubble,
	shutdowner fx.Shutdowner,
	logger *log.Logger,
) {
	go func() {
		_, err := tea.NewProgram(bubble).Run()
		logger.Info("shutting down")
		if err != nil {
			logger.Error("error shutting down tui", "err", err)
		}
		shutdowner.Shutdown()
	}()
}
