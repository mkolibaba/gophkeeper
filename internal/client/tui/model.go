package tui

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
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
	"strings"
	"unicode/utf8"
)

// Bubble представляет состояние UI.
type Bubble struct {
	width  int // Ширина терминала
	height int // Высота терминала

	dump io.Writer

	dataTable  table.Model
	dataDetail detail.Model

	manager *state.Manager

	view view

	authModel authorizationModel
}

type view int

const (
	ViewAuthorization view = iota
	ViewMain
	ViewAddData
)

type keyMap struct {
	UpDown         key.Binding
	AddData        key.Binding
	DownloadBinary key.Binding
	Quit           key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.UpDown, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.UpDown},
		{k.AddData, k.DownloadBinary},
		{k.Quit},
	}
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

	keys = keyMap{
		UpDown: key.NewBinding(
			key.WithKeys("up", "down"),
			key.WithHelp("↑/↓", "move up/down"),
		),
		AddData: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add data"),
		),
		DownloadBinary: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "download binary"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
)

// NewBubble создает новый экземпляр UI.
func NewBubble(
	authorizationService client.AuthorizationService,
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
		authorizationService,
		loginService,
		noteService,
		binaryService,
		cardService,
		logger,
	)

	return Bubble{
		dump:       dump,
		dataTable:  dataTable,
		dataDetail: dataDetail,
		manager:    manager,
		authModel:  initialAuthorizationModel(manager),
	}, nil
}

// Init инициализирует UI.
func (b Bubble) Init() tea.Cmd {
	return tea.Batch(
		//b.manager.FetchData(),
		b.authModel.Init(),
	)
}

// Update обновляет UI в зависимости от события.
func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	b.spew(msg)

	// Корневые события
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.height, b.width = msg.Height, msg.Width
		return b, cmd

	case changeViewMsg:
		b.view = msg.view
		if b.view == ViewMain {
			return b, b.manager.FetchData()
		}
	}

	switch b.view {
	case ViewAuthorization:
		b.authModel, cmd = b.authModel.Update(msg)
		return b, cmd
	case ViewMain:
	case ViewAddData:
	}

	// TODO: все, что ниже, поместить во viewmain case
	switch msg := msg.(type) {
	case state.FetchDataMsg:
		b.dataTable, cmd = b.dataTable.Update(msg)
		current := b.dataTable.GetCurrentRow()
		b.dataDetail = b.dataDetail.SetData(current)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.UpDown):
			b.dataTable, cmd = b.dataTable.Update(msg)
			current := b.dataTable.GetCurrentRow()
			b.dataDetail = b.dataDetail.SetData(current)

		case key.Matches(msg, keys.Quit):
			return b, tea.Quit

		case key.Matches(msg, keys.DownloadBinary):
			current := b.dataTable.GetCurrentRow()
			if d, ok := current.(client.BinaryData); ok {
				b.manager.StartDownloadBinary(d)
			}

		case key.Matches(msg, keys.AddData):
			if b.view == ViewMain {
				b.view = ViewAddData
			} else if b.view == ViewAddData {
				b.view = ViewMain
			}
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

	contentHeight := b.height - lipgloss.Height(title)

	var content string
	switch b.view {
	case ViewAuthorization:
		content = b.renderAuthorizationView(contentHeight)
	case ViewMain:
		content = b.renderMainView(contentHeight)
	case ViewAddData:
		content = b.renderAddDataView(contentHeight)
	}

	return lipgloss.JoinVertical(lipgloss.Top, title, content)
}

func (b Bubble) renderAuthorizationView(height int) string {
	// Вид приложения:
	// +------------------------+
	// | Title                  |
	// +--------------+---------+
	// | Authorization View     |
	// |                        |
	// +--------------+---------+

	w := b.width - contentStyle.GetHorizontalFrameSize()

	borderTop := b.renderBorderTop(contentStyle, "Authorization", w)

	h := height - lipgloss.Height(borderTop) - contentStyle.GetBorderBottomSize()

	v := b.authModel.View()

	authorizationView := contentStyle.
		BorderTop(false).
		Width(w).
		Height(h / 2).
		PaddingLeft(1).
		Render(v)

	return lipgloss.JoinVertical(lipgloss.Top, borderTop, authorizationView)
}

func (b Bubble) renderAddDataView(height int) string {
	// Вид приложения:
	// +------------------------+
	// | Title                  |
	// +--------------+---------+
	// | Add View               |
	// |                        |
	// +--------------+---------+

	w := b.width - contentStyle.GetHorizontalFrameSize()

	borderTop := b.renderBorderTop(contentStyle, "Add Data", w)

	addDataView := contentStyle.
		BorderTop(false).
		Width(w).
		Height(height - lipgloss.Height(borderTop) - contentStyle.GetBorderBottomSize()).
		PaddingLeft(1).
		Render("Тут можно будет добавлять всякий контент")

	return lipgloss.JoinVertical(lipgloss.Top, borderTop, addDataView)
}

func (b Bubble) renderMainView(height int) string {
	// Вид приложения:
	// +------------------------+
	// | Title                  |
	// +--------------+---------+
	// | Table        | Detail  |
	// | View         | View    |
	// +--------------+---------+

	// Строка помощи
	hm := help.New()
	hm.ShowAll = true
	helpView := lipgloss.NewStyle().PaddingLeft(1).Render(hm.View(keys))

	h := height - lipgloss.Height(helpView)

	// Окно со списком данных
	tableView := b.renderTableView(h)

	// Окно детального просмотра
	detailViewWidth := b.width - lipgloss.Width(tableView)
	detailView := b.renderDetailView(detailViewWidth, h)

	return lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Top, tableView, detailView),
		helpView,
	)
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
		Height(height - contentStyle.GetBorderBottomSize() - lipgloss.Height(detailTop)).
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

func (b Bubble) spew(a ...any) {
	spew.Fdump(b.dump, a...)
}

// -------------
// Authorization
// -------------

type changeViewMsg struct {
	view view
}

func changeView(view view) tea.Cmd {
	return func() tea.Msg {
		return changeViewMsg{view: view}
	}
}

type authorizationModel struct {
	inputs  []textinput.Model
	focused int
	err     error
	manager *state.Manager
}

func initialAuthorizationModel(manager *state.Manager) authorizationModel {
	loginInput := textinput.New()
	loginInput.Placeholder = "Login"
	loginInput.Focus()
	loginInput.CharLimit = 20
	loginInput.Width = 20
	loginInput.Cursor.SetMode(cursor.CursorStatic)
	loginInput.PromptStyle = HeaderStyle

	passwordInput := textinput.New()
	passwordInput.Placeholder = "Password"
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '•'
	passwordInput.CharLimit = 20
	passwordInput.Width = 20
	passwordInput.Cursor.SetMode(cursor.CursorStatic)
	passwordInput.PromptStyle = HeaderStyle

	return authorizationModel{
		inputs:  []textinput.Model{loginInput, passwordInput},
		manager: manager,
	}
}

func (m authorizationModel) Init() tea.Cmd {
	return nil
}

func (m authorizationModel) Update(msg tea.Msg) (authorizationModel, tea.Cmd) {
	switch msg := msg.(type) {
	case state.AuthorizationResultMsg:
		m.err = msg.Err
		if m.err == nil {
			return m, changeView(ViewMain)
		}
		m.inputs[0].SetValue("")
		m.inputs[1].SetValue("")
		m = m.setFocus(0)
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "up", "down", "tab":
			if keypress == "up" {
				m = m.setFocus((m.focused - 1 + len(m.inputs)) % len(m.inputs))
			} else {
				m = m.setFocus((m.focused + 1) % len(m.inputs))
			}
		case "enter":
			login, password := m.inputs[0].Value(), m.inputs[1].Value()
			return m, m.manager.Authorize(login, password)
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m authorizationModel) View() string {
	lines := []string{m.inputs[0].View(), m.inputs[1].View()}
	if m.err != nil {
		lines = append(lines, "", authErrorRenderer(m.err.Error()))
	}
	return authViewRenderer(lipgloss.JoinVertical(lipgloss.Top, lines...))
}

func (m authorizationModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m authorizationModel) setFocus(at int) authorizationModel {
	m.focused = at
	for i := range m.inputs {
		if at == i {
			m.inputs[i].Focus()
			continue
		}
		m.inputs[i].Blur()
	}
	return m
}

var (
	authErrorRenderer = lipgloss.NewStyle().
				Foreground(lipgloss.Color("169")).
				Render

	authViewRenderer = lipgloss.NewStyle().
				PaddingTop(1).
				Render
)
