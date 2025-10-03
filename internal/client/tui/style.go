package tui

import "github.com/charmbracelet/lipgloss"

var (
	HeaderColor = lipgloss.Color("171")

	HeaderStyle = lipgloss.NewStyle().
			Foreground(HeaderColor)
)
