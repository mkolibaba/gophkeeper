package table

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mkolibaba/gophkeeper/internal/client"
	"go.uber.org/zap"
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
	cursor int
	rows   []Row

	loginService  client.LoginService
	noteService   client.NoteService
	binaryService client.BinaryService
	cardService   client.CardService

	logger *zap.Logger
}

func NewModel(
	loginService client.LoginService,
	noteService client.NoteService,
	binaryService client.BinaryService,
	cardService client.CardService,
	logger *zap.Logger,
) Model {
	return Model{
		loginService:  loginService,
		noteService:   noteService,
		binaryService: binaryService,
		cardService:   cardService,
		logger:        logger,
	}
}

func (m Model) Init() tea.Cmd {
	return FetchData(
		m.loginService,
		m.noteService,
		m.binaryService,
		m.cardService,
		m.logger,
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case FetchDataMsg:
		m = m.processFetchedData(msg)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.rows)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	rows := []string{m.renderRow("Type", "Name", "Value", false)}
	for i, row := range m.rows {
		rows = append(rows, m.renderRow(row.DataType, row.Name, row.RenderedValue, i == m.cursor))
	}
	return lipgloss.JoinVertical(lipgloss.Top, rows...)
}

func (m Model) renderRow(t DataType, name, value string, selected bool) string {
	columns := []string{
		columnStyle.Width(10).Render(string(t)),
		columnStyle.Render(name),
		columnStyle.Render(value),
	}

	rowStyle := lipgloss.NewStyle().
		Inline(true)
	if selected {
		rowStyle = rowStyle.Background(lipgloss.Color("171"))
	}
	return rowStyle.
		Render(lipgloss.JoinHorizontal(lipgloss.Left, columns...))
}

func (m Model) processFetchedData(msg FetchDataMsg) Model {
	m.rows = make([]Row, 0, len(msg.Logins)+len(msg.Notes)+len(msg.Binaries)+len(msg.Cards))
	for _, login := range msg.Logins {
		m.rows = append(m.rows, Row{
			DataType:      DataTypeLogin,
			Name:          login.Name,
			Value:         login.Login,
			RenderedValue: login.Login,
		})
	}
	for _, note := range msg.Notes {
		m.rows = append(m.rows, Row{
			DataType:      DataTypeNote,
			Name:          note.Name,
			Value:         note.Text,
			RenderedValue: trimNoteText(note.Text),
		})
	}
	for _, binary := range msg.Binaries {
		m.rows = append(m.rows, Row{
			DataType:      DataTypeBinary,
			Name:          binary.Name,
			Value:         "<binary>",
			RenderedValue: "<binary>",
		})
	}
	for _, card := range msg.Cards {
		m.rows = append(m.rows, Row{
			DataType:      DataTypeCard,
			Name:          card.Name,
			Value:         card.Number,
			RenderedValue: maskCardNumber(card.Number),
		})
	}
	return m
}

func trimNoteText(text string) string {
	maxLength := 30
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
