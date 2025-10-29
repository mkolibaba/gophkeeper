package view

import tea "github.com/charmbracelet/bubbletea"

// View - тип состояния UI.
type View uint8

const (
	// ViewAuthorization - авторизация пользователя.
	ViewAuthorization View = iota

	// ViewHome - основное окно приложения.
	ViewHome

	// ViewAddData - окно добавления данных.
	ViewAddData

	// ViewRegistration - регистрация нового пользователя.
	ViewRegistration

	// ViewEditData - окно редактирования данных.
	ViewEditData
)

// Model - представление состояния UI.
type Model interface {
	Init() tea.Cmd
	Update(tea.Msg) tea.Cmd
	View() string
	SetSize(width int, height int)
}

type BaseModel struct {
	// Ширина компонента.
	Width int
	// Высота компонента.
	Height int
}

func (m *BaseModel) SetSize(width int, height int) {
	m.Width = width
	m.Height = height
}
