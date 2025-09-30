package detail

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mkolibaba/gophkeeper/internal/client"
)

var (
	fieldStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("171"))
)

type Model struct {
	Data client.Data
}

func NewModel() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	switch d := m.Data.(type) {
	case client.LoginData:
		return lipgloss.JoinVertical(lipgloss.Top,
			fieldStyle.Render("Type"),
			"Login",
			"",
			fieldStyle.Render("Name"),
			d.Name,
			"",
			fieldStyle.Render("Login"),
			d.Login,
			"",
			fieldStyle.Render("Password"),
			d.Password,
		)
	case client.NoteData:
		return lipgloss.JoinVertical(lipgloss.Top,
			fieldStyle.Render("Type"),
			"Note",
			"",
			fieldStyle.Render("Name"),
			d.Name,
			"",
			fieldStyle.Render("Text"),
			d.Text,
		)
	case client.BinaryData:
		return lipgloss.JoinVertical(lipgloss.Top,
			fieldStyle.Render("Type"),
			"Binary",
			"",
			fieldStyle.Render("Name"),
			d.Name,
			"",
			fieldStyle.Render("File"),
			"<binary>",
		)
	case client.CardData:
		return lipgloss.JoinVertical(lipgloss.Top,
			fieldStyle.Render("Type"),
			"Card",
			"",
			fieldStyle.Render("Name"),
			d.Name,
			"",
			fieldStyle.Render("Number"),
			d.Number,
			"",
			fieldStyle.Render("Expiry date"),
			d.ExpDate,
			"",
			fieldStyle.Render("CVV"),
			d.CVV,
			"",
			fieldStyle.Render("Cardholder"),
			d.Cardholder,
		)
	}

	return ""
}

func (m Model) SetData(data client.Data) Model {
	m.Data = data
	return m
}
