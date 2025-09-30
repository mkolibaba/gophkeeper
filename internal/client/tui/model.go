package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
	"github.com/mkolibaba/gophkeeper/internal/client"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/detail"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/state"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/table"
	"go.uber.org/zap"
	"io"
	"os"
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
	titleStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Background(lipgloss.Color("105")).
			SetString("Gophkeeper")

	contentStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("141"))
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
	//spew.Fdump(b.dump, msg)

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
		}
	case tea.WindowSizeMsg:
		b.height, b.width = msg.Height, msg.Width
	}

	return b, cmd
}

// View возвращает строковое представление UI.
func (b Bubble) View() string {
	// Заголовок приложения
	title := titleStyle.
		Width(b.width).
		Render()

	// Окно со списком данных
	contentLeftBottom := contentStyle.
		Width(b.width/3*2 - contentStyle.GetHorizontalFrameSize()).
		PaddingRight(1).
		Align(lipgloss.Right).
		Render(b.dataTable.RenderInfoBar())

	contentLeft := contentStyle.
		Width(b.width/3*2 - contentStyle.GetHorizontalFrameSize()).
		Height(b.height - lipgloss.Height(title) - contentStyle.GetVerticalFrameSize() - lipgloss.Height(contentLeftBottom)).
		PaddingLeft(1).
		Render(b.dataTable.View())

	spew.Fdump(b.dump, b.dataDetail.Data)

	// Окно детального просмотра
	contentRight := contentStyle.
		Width(b.width - lipgloss.Width(contentLeft) - contentStyle.GetHorizontalFrameSize()).
		Height(b.height - lipgloss.Height(title) - contentStyle.GetVerticalFrameSize()).
		PaddingLeft(1).
		Render(b.dataDetail.View())

	content := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Top, contentLeft, contentLeftBottom),
		contentRight,
	)

	return lipgloss.JoinVertical(lipgloss.Top, title, content)
}
