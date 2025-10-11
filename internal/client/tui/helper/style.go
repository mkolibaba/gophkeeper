package helper

import "github.com/charmbracelet/lipgloss"

var (
	HeaderColor = lipgloss.Color("171")

	HeaderStyle = lipgloss.NewStyle().
			Foreground(HeaderColor)

	BorderColor = lipgloss.Color("141")

	// Стиль стандартных рамок приложения.
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(BorderColor)

	customBorderStyle = lipgloss.NewStyle().
				Foreground(BorderColor)

	TitleStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Background(lipgloss.Color("105")).
			SetString("Gophkeeper")
)
