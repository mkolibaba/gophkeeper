package table

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mkolibaba/gophkeeper/client"
	"github.com/mkolibaba/gophkeeper/client/tui/helper"
	"regexp"
)

var (
	columnStyle = lipgloss.NewStyle().Inline(true)

	maskingCardNumberRegexp = regexp.MustCompile(`(\d{6})\d{6}(\d{4})`)
	spacingCardNumberRegexp = regexp.MustCompile(`(.{4})(.{4})(.{4})(.{4})`)
)

const (
	typeWidth = 10
	nameWidth = 30
)

type Row struct {
	DataType      helper.DataType
	Name          string
	Value         string
	RenderedValue string
}

type Model struct {
	cursor       int
	data         []client.Data
	renderedRows []Row

	valueWidth int
}

func New() *Model {
	return &Model{
		valueWidth: 70,
	}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.cursor = (m.cursor - 1 + len(m.data)) % len(m.data)
		case "down":
			m.cursor = (m.cursor + 1) % len(m.data)
		}
	}

	return nil
}

func (m *Model) View() string {
	if len(m.renderedRows) == 0 {
		return "No data"
	}

	rows := []string{m.renderHeader()}
	for i, row := range m.renderedRows {
		rows = append(rows, m.renderRow(row.DataType, row.Name, row.RenderedValue, i == m.cursor))
	}
	return lipgloss.JoinVertical(lipgloss.Top, rows...)
}

func (m *Model) GetCurrentRow() client.Data {
	if len(m.data) == 0 {
		return nil
	}

	return m.data[m.cursor]
}

func (m *Model) SetWidth(width int) {
	m.valueWidth = width - typeWidth - nameWidth
}

func (m *Model) RenderInfoBar() string {
	// TODO: выводить информацию, соответствующую действительности (viewport)
	//return fmt.Sprintf("%d/%d", m.cursor+1, len(m.renderedRows))
	return fmt.Sprintf("1-%d of %d", len(m.renderedRows), len(m.renderedRows))
}

func (m Model) renderHeader() string {
	return lipgloss.NewStyle().
		Inline(true).
		Render(lipgloss.JoinHorizontal(lipgloss.Left,
			helper.HeaderStyle.Width(typeWidth).Render("Type"),
			helper.HeaderStyle.Width(nameWidth).Render("Name"),
			helper.HeaderStyle.Width(m.valueWidth).Render("Value"),
		))
}

func (m Model) renderRow(t helper.DataType, name, value string, selected bool) string {
	columns := []string{
		columnStyle.Width(typeWidth).Render(string(t)),
		columnStyle.Width(nameWidth).Render(name),
		columnStyle.Width(m.valueWidth).Render(value),
	}

	rowStyle := lipgloss.NewStyle().
		Inline(true)
	if selected {
		rowStyle = rowStyle.Background(helper.HeaderColor)
	}
	return rowStyle.
		Render(lipgloss.JoinHorizontal(lipgloss.Left, columns...))
}

func (m *Model) ProcessFetchedData(msg []client.Data) {
	m.data = msg
	m.cursor = min(max(0, m.cursor), len(m.data)-1)

	m.renderedRows = make([]Row, 0, len(m.data))
	for _, el := range m.data {
		switch el := el.(type) {
		case client.LoginData:
			m.renderedRows = append(m.renderedRows, Row{
				DataType:      helper.DataTypeLogin,
				Name:          el.Name,
				Value:         el.Login,
				RenderedValue: el.Login,
			})
		case client.NoteData:
			m.renderedRows = append(m.renderedRows, Row{
				DataType:      helper.DataTypeNote,
				Name:          el.Name,
				Value:         el.Text,
				RenderedValue: m.trimNoteText(el.Text),
			})
		case client.BinaryData:
			m.renderedRows = append(m.renderedRows, Row{
				DataType:      helper.DataTypeBinary,
				Name:          el.Name,
				Value:         "<binary>",
				RenderedValue: "<binary>",
			})
		case client.CardData:
			m.renderedRows = append(m.renderedRows, Row{
				DataType:      helper.DataTypeCard,
				Name:          el.Name,
				Value:         el.Number,
				RenderedValue: maskCardNumber(el.Number),
			})
		}
	}
}

func (m *Model) trimNoteText(text string) string {
	asRunes := []rune(text) // TODO: может есть лучше решение?
	if len(asRunes) > m.valueWidth {
		return string(asRunes[:m.valueWidth-3]) + "..."
	}
	return text
}

func maskCardNumber(number string) string {
	masked := maskingCardNumberRegexp.ReplaceAllString(number, "$1******$2")
	return spacingCardNumberRegexp.ReplaceAllString(masked, "$1 $2 $3 $4")
}
