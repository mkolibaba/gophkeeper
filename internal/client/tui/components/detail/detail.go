package detail

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/mkolibaba/gophkeeper/internal/client"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/helper"
)

var (
	fieldStyle = helper.HeaderStyle
)

type Model struct {
	Data client.Data
}

func NewModel() Model {
	return Model{}
}

func (m Model) View() string {
	var lines []string

	switch d := m.Data.(type) {
	case client.LoginData:
		lines = []string{
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
			"",
			fieldStyle.Render("Website"),
			d.Website,
			"",
			fieldStyle.Render("Notes"),
			d.Notes,
		}
	case client.NoteData:
		lines = []string{
			fieldStyle.Render("Type"),
			"Note",
			"",
			fieldStyle.Render("Name"),
			d.Name,
			"",
			fieldStyle.Render("Text"),
			d.Text,
		}
	case client.BinaryData:
		lines = []string{
			fieldStyle.Render("Type"),
			"Binary",
			"",
			fieldStyle.Render("Name"),
			d.Name,
			"",
			fieldStyle.Render("File name"),
			d.Filename,
			"",
			fieldStyle.Render("Size"),
			fmt.Sprintf("%d", d.Size),
			"",
			fieldStyle.Render("Notes"),
			d.Notes,
		}
	case client.CardData:
		lines = []string{
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
			"",
			fieldStyle.Render("Notes"),
			d.Notes,
		}
	case nil:
		lines = []string{"No data"}
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}
