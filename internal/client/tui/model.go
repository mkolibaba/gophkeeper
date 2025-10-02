package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mkolibaba/gophkeeper/internal/client"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/detail"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/state"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/table"
	"go.uber.org/zap"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

// Bubble представляет состояние UI.
type Bubble struct {
	width  int // Ширина терминала
	height int // Высота терминала

	loginService  client.LoginService
	noteService   client.NoteService
	binaryService client.BinaryService
	cardService   client.CardService

	dump io.Writer

	dataTable  table.Model
	dataDetail detail.Model

	manager *state.Manager
}

var (
	borderColor = lipgloss.Color("141")

	titleStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Background(lipgloss.Color("105")).
			SetString("Gophkeeper")

	contentStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(borderColor)

	customBorderStyle = lipgloss.NewStyle().
				Foreground(borderColor)
)

// NewBubble создает новый экземпляр UI.
func NewBubble(
	loginService client.LoginService,
	noteService client.NoteService,
	binaryService client.BinaryService,
	cardService client.CardService,
	logger *zap.Logger,
) (Bubble, error) {
	var dump *os.File
	if dumpPath, ok := os.LookupEnv("SPEW_DUMP_OUTPUT"); ok {
		var err error
		dump, err = os.OpenFile(dumpPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			return Bubble{}, err
		}
	}

	dataTable := table.NewModel()
	dataDetail := detail.NewModel()

	manager := state.NewManager(
		loginService,
		noteService,
		binaryService,
		cardService,
		logger,
	)

	return Bubble{
		loginService:  loginService,
		noteService:   noteService,
		binaryService: binaryService,
		cardService:   cardService,
		dump:          dump,
		dataTable:     dataTable,
		dataDetail:    dataDetail,
		manager:       manager,
	}, nil
}

// Init инициализирует UI.
func (b Bubble) Init() tea.Cmd {
	return b.manager.FetchData()
}

// Update обновляет UI в зависимости от события.
func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case state.FetchDataMsg:
		b.dataTable, cmd = b.dataTable.Update(msg)
		current := b.dataTable.GetCurrentRow()
		b.dataDetail = b.dataDetail.SetData(current)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "up", "down":
			b.dataTable, cmd = b.dataTable.Update(msg)
			current := b.dataTable.GetCurrentRow()
			b.dataDetail = b.dataDetail.SetData(current)
		case "ctrl+c", "q":
			return b, tea.Quit
		case "d":
			current := b.dataTable.GetCurrentRow()
			if d, ok := current.(client.BinaryData); ok {
				b.manager.StartDownloadBinary(d)
			}
		}
	case tea.WindowSizeMsg:
		b.height, b.width = msg.Height, msg.Width
	}

	return b, cmd
}

// View возвращает строковое представление UI.
func (b Bubble) View() string {
	// Вид приложения:
	// +------------------------+
	// | Title                  |
	// +--------------+---------+
	// | Table        | Detail  |
	// | View         | View    |
	// +--------------+---------+

	// Заголовок приложения
	title := titleStyle.
		Width(b.width).
		Render()

	viewsHeight := b.height - lipgloss.Height(title)

	// Окно со списком данных
	tableView := b.renderTableView(viewsHeight)

	// Окно детального просмотра
	detailViewWidth := b.width - lipgloss.Width(tableView)
	detailView := b.renderDetailView(detailViewWidth, viewsHeight)

	content := lipgloss.JoinHorizontal(lipgloss.Top, tableView, detailView)

	return lipgloss.JoinVertical(lipgloss.Top, title, content)
}

func (b Bubble) renderTableView(height int) string {
	w := b.width/3*2 - contentStyle.GetHorizontalFrameSize()

	tableTopBorder := b.renderBorderTop(contentStyle, "Data", w)

	tableBottomBorder := b.renderBorderBottom(contentStyle, b.dataTable.RenderInfoBar(), w)

	tableView := contentStyle.
		BorderTop(false).
		BorderBottom(false).
		Width(w).
		Height(height - lipgloss.Height(tableTopBorder) - lipgloss.Height(tableBottomBorder)).
		PaddingLeft(1).
		Render(b.dataTable.View())

	return lipgloss.JoinVertical(lipgloss.Top, tableTopBorder, tableView, tableBottomBorder)
}

func (b Bubble) renderDetailView(width int, height int) string {
	w := width - contentStyle.GetHorizontalFrameSize()

	detailTop := b.renderBorderTop(contentStyle, "Detail", w)

	detailView := contentStyle.
		BorderTop(false).
		Width(w).
		Height(height - contentStyle.GetVerticalFrameSize()).
		PaddingLeft(1).
		Render(b.dataDetail.View())

	return lipgloss.JoinVertical(lipgloss.Top, detailTop, detailView)
}

func (b Bubble) renderBorderTop(style lipgloss.Style, text string, width int) string {
	border, _, _, _, _ := style.GetBorder()
	borderLeft := border.TopLeft
	borderMiddle := border.Top
	borderRight := border.TopRight

	leftText := borderLeft + borderRight

	rightRepeat := width -
		utf8.RuneCountInString(leftText) -
		utf8.RuneCountInString(text)
	// TODO: fix workaround
	if rightRepeat <= 0 {
		rightRepeat = 1
	}

	rightText := borderLeft + strings.Repeat(borderMiddle, rightRepeat) + borderRight

	left := customBorderStyle.Render(leftText)
	right := customBorderStyle.Render(rightText)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, text, right)
}

func (b Bubble) renderBorderBottom(style lipgloss.Style, text string, width int) string {
	border, _, _, _, _ := style.GetBorder()
	borderLeft := border.BottomLeft
	borderMiddle := border.Bottom
	borderRight := border.BottomRight

	rightText := borderLeft + borderRight

	leftRepeat := width -
		utf8.RuneCountInString(rightText) -
		utf8.RuneCountInString(text)
	// TODO: fix workaround
	if leftRepeat <= 0 {
		leftRepeat = 1
	}

	leftText := borderLeft + strings.Repeat(borderMiddle, leftRepeat) + borderRight

	left := customBorderStyle.Render(leftText)
	right := customBorderStyle.Render(rightText)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, text, right)
}
