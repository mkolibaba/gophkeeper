package view

import tea "github.com/charmbracelet/bubbletea"

type Model interface {
	Init() tea.Cmd
	Update(tea.Msg) tea.Cmd
	View() string
	SetSize(width int, height int)
}

type baseViewModel struct {
	// Ширина компонента.
	Width int
	// Высота компонента.
	Height int
}

func (m *baseViewModel) SetSize(width int, height int) {
	m.Width = width
	m.Height = height
}
