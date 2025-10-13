package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"go.uber.org/fx"
)

func Start(
	bubble Bubble,
	shutdowner fx.Shutdowner,
	logger *log.Logger,
) {
	go func() {
		_, err := tea.NewProgram(bubble, tea.WithAltScreen()).Run()
		logger.Info("shutting down")
		if err != nil {
			logger.Error("error shutting down tui", "err", err)
		}
		shutdowner.Shutdown()
	}()
}
