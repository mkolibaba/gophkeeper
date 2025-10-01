package table

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mkolibaba/gophkeeper/internal/client"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/state"
	"regexp"
)

type DataType string

const (
	DataTypeLogin  = DataType("Login")
	DataTypeNote   = DataType("Note")
	DataTypeBinary = DataType("Binary")
	DataTypeCard   = DataType("Card")
)

var (
	columnStyle = lipgloss.NewStyle().
			Width(30).
			Inline(true).
			Padding(0, 1)

	maskingCardNumberRegexp = regexp.MustCompile(`(\d{6})\d{6}(\d{4})`)
	spacingCardNumberRegexp = regexp.MustCompile(`(.{4})(.{4})(.{4})(.{4})`)
)

type Row struct {
	DataType      DataType
	Name          string
	Value         string
	RenderedValue string
}

type Model struct {
	cursor       int
	data         []client.Data
	renderedRows []Row
	manager      *state.Manager
}

func NewModel() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case state.FetchDataMsg:
		m = m.processFetchedData(msg)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.cursor = (m.cursor - 1 + len(m.data)) % len(m.data)
		case "down":
			m.cursor = (m.cursor + 1) % len(m.data)
		}
	}

	return m, nil
}

func (m Model) View() string {
	rows := []string{m.renderRow("Type", "Name", "Value", false)}
	for i, row := range m.renderedRows {
		rows = append(rows, m.renderRow(row.DataType, row.Name, row.RenderedValue, i == m.cursor))
	}
	return lipgloss.JoinVertical(lipgloss.Top, rows...)
}

func (m Model) GetCurrentRow() client.Data {
	return m.data[m.cursor]
}

func (m Model) RenderInfoBar() string {
	return fmt.Sprintf("%d/%d", m.cursor+1, len(m.renderedRows))
}

func (m Model) renderRow(t DataType, name, value string, selected bool) string {
	columns := []string{
		columnStyle.Width(10).Render(string(t)),
		columnStyle.Width(30).Render(name),
		columnStyle.Width(52).Render(value), // TODO: должна как-то определяться родительская ширина
	}

	rowStyle := lipgloss.NewStyle().
		Inline(true)
	if selected {
		rowStyle = rowStyle.Background(lipgloss.Color("171"))
	}
	return rowStyle.
		Render(lipgloss.JoinHorizontal(lipgloss.Left, columns...))
}

func (m Model) processFetchedData(msg state.FetchDataMsg) Model {
	m.data = []client.Data{}
	for _, login := range msg.Logins {
		m.data = append(m.data, login)
	}
	for _, note := range msg.Notes {
		m.data = append(m.data, note)
	}
	for _, binary := range msg.Binaries {
		m.data = append(m.data, binary)
	}
	for _, card := range msg.Cards {
		m.data = append(m.data, card)
	}

	m.renderedRows = make([]Row, 0, len(msg.Logins)+len(msg.Notes)+len(msg.Binaries)+len(msg.Cards))
	for _, login := range msg.Logins {
		m.renderedRows = append(m.renderedRows, Row{
			DataType:      DataTypeLogin,
			Name:          login.Name,
			Value:         login.Login,
			RenderedValue: login.Login,
		})
	}
	for _, note := range msg.Notes {
		m.renderedRows = append(m.renderedRows, Row{
			DataType:      DataTypeNote,
			Name:          note.Name,
			Value:         note.Text,
			RenderedValue: trimNoteText(note.Text),
		})
	}
	for _, binary := range msg.Binaries {
		m.renderedRows = append(m.renderedRows, Row{
			DataType:      DataTypeBinary,
			Name:          binary.Name,
			Value:         "<binary>",
			RenderedValue: "<binary>",
		})
	}
	for _, card := range msg.Cards {
		m.renderedRows = append(m.renderedRows, Row{
			DataType:      DataTypeCard,
			Name:          card.Name,
			Value:         card.Number,
			RenderedValue: maskCardNumber(card.Number),
		})
	}
	return m
}

func trimNoteText(text string) string {
	maxLength := 50
	asRunes := []rune(text) // TODO: может есть лучше решение?
	if len(asRunes) > maxLength {
		return string(asRunes[:maxLength-3]) + "..."
	}
	return text
}

func maskCardNumber(number string) string {
	masked := maskingCardNumberRegexp.ReplaceAllString(number, "$1******$2")
	return spacingCardNumberRegexp.ReplaceAllString(masked, "$1 $2 $3 $4")
}
