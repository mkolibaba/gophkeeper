package table

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mkolibaba/gophkeeper/internal/client"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/orchestrator"
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
	case orchestrator.LoadDataMsg:
		m = m.processFetchedData(msg)
		m.cursor = min(max(0, m.cursor), len(m.data)-1) // TODO: написать зачем это нужно (удаление последнего эелемента)
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
	if len(m.renderedRows) == 0 {
		return "No data"
	}

	rows := []string{m.renderHeader()}
	for i, row := range m.renderedRows {
		rows = append(rows, m.renderRow(row.DataType, row.Name, row.RenderedValue, i == m.cursor))
	}
	return lipgloss.JoinVertical(lipgloss.Top, rows...)
}

func (m Model) GetCurrentRow() client.Data {
	if len(m.data) == 0 {
		return nil
	}

	return m.data[m.cursor]
}

func (m Model) RenderInfoBar() string {
	// TODO: выводить информацию, соответствующую действительности
	//return fmt.Sprintf("%d/%d", m.cursor+1, len(m.renderedRows))
	return fmt.Sprintf("1-%d of %d", len(m.renderedRows), len(m.renderedRows))
}

func (m Model) renderHeader() string {
	return lipgloss.NewStyle().
		Inline(true).
		Render(lipgloss.JoinHorizontal(lipgloss.Left,
			columnStyle.Width(10).Foreground(lipgloss.Color("171")).Render("Type"),
			columnStyle.Width(30).Foreground(lipgloss.Color("171")).Render("Name"),
			columnStyle.Width(52).Foreground(lipgloss.Color("171")).Render("Value"), // TODO: должна как-то определяться родительская ширина
		))
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

func (m Model) processFetchedData(msg orchestrator.LoadDataMsg) Model {
	m.data = msg

	m.renderedRows = make([]Row, 0, len(m.data))
	for _, el := range m.data {
		switch el := el.(type) {
		case client.LoginData:
			m.renderedRows = append(m.renderedRows, Row{
				DataType:      DataTypeLogin,
				Name:          el.Name,
				Value:         el.Login,
				RenderedValue: el.Login,
			})
		case client.NoteData:
			m.renderedRows = append(m.renderedRows, Row{
				DataType:      DataTypeNote,
				Name:          el.Name,
				Value:         el.Text,
				RenderedValue: trimNoteText(el.Text),
			})
		case client.BinaryData:
			m.renderedRows = append(m.renderedRows, Row{
				DataType:      DataTypeBinary,
				Name:          el.Name,
				Value:         "<binary>",
				RenderedValue: "<binary>",
			})
		case client.CardData:
			m.renderedRows = append(m.renderedRows, Row{
				DataType:      DataTypeCard,
				Name:          el.Name,
				Value:         el.Number,
				RenderedValue: maskCardNumber(el.Number),
			})
		}
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
