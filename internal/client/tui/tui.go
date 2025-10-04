package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Start(
	bubble Bubble,
	shutdowner fx.Shutdowner,
	logger *zap.Logger,
) {
	go func() {
		_, err := tea.NewProgram(bubble, tea.WithAltScreen()).Run()
		logger.Info("shutting down")
		if err != nil {
			logger.Error("error shutting down tui", zap.Error(err))
		}
		shutdowner.Shutdown()
	}()
}
