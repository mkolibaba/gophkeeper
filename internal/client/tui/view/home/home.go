package home

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mkolibaba/gophkeeper/internal/client"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/components/detail"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/components/statusbar"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/components/table"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/helper"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/view"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/view/adddata"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/view/authorization"
	"go.uber.org/fx"
	"sync"
)

// TODO: название не нравится, переименовать

// CallAddDataViewMsg отправляется при вызове пользователем окна добавления данных.
type CallAddDataViewMsg helper.DataType

func CallAddDataView(dataType helper.DataType) tea.Cmd {
	return func() tea.Msg {
		return CallAddDataViewMsg(dataType)
	}
}

type loadDataMsg []client.Data

type keyMap struct {
	UpDown         key.Binding
	AddLogin       key.Binding
	AddNote        key.Binding
	AddBinary      key.Binding
	AddCard        key.Binding
	DownloadBinary key.Binding
	Remove         key.Binding
	Help           key.Binding
	Quit           key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.UpDown, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.UpDown},
		{k.AddLogin, k.AddNote, k.AddBinary, k.AddCard},
		{k.DownloadBinary, k.Remove},
		{k.Quit},
	}
}

type Model struct {
	view.BaseModel
	dataTable     *table.Model
	dataDetail    detail.Model
	keyMap        keyMap
	showHelp      bool
	statusBar     *statusbar.Model
	loginService  client.LoginService
	binaryService client.BinaryService
	noteService   client.NoteService
	cardService   client.CardService
	userService   client.UserService
}

type Params struct {
	fx.In

	LoginService  client.LoginService
	BinaryService client.BinaryService
	NoteService   client.NoteService
	CardService   client.CardService
	UserService   client.UserService
}

func New(p Params) *Model {
	keys := keyMap{
		UpDown: key.NewBinding(
			key.WithKeys("up", "down"),
			key.WithHelp("↑/↓", "move up/down"),
		),
		AddLogin: key.NewBinding(
			key.WithKeys("alt+1"),
			key.WithHelp("alt+1", "add login"),
		),
		AddNote: key.NewBinding(
			key.WithKeys("alt+2"),
			key.WithHelp("alt+2", "add note"),
		),
		AddBinary: key.NewBinding(
			key.WithKeys("alt+3"),
			key.WithHelp("alt+3", "add binary"),
		),
		AddCard: key.NewBinding(
			key.WithKeys("alt+4"),
			key.WithHelp("alt+4", "add card"),
		),
		Remove: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("ctrl+r", "remove"),
		),
		DownloadBinary: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "download binary"),
		),
		Help: key.NewBinding(
			key.WithKeys("h"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}

	dataTable := table.New()
	dataDetail := detail.New()
	statusBar := statusbar.New()

	return &Model{
		dataTable:     dataTable,
		dataDetail:    dataDetail,
		statusBar:     statusBar,
		keyMap:        keys,
		loginService:  p.LoginService,
		binaryService: p.BinaryService,
		noteService:   p.NoteService,
		cardService:   p.CardService,
		userService:   p.UserService,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case loadDataMsg:
		m.statusBar.CurrentUser = m.userService.Get().Login
		m.dataTable.ProcessFetchedData(msg)
		m.dataDetail.Data = m.dataTable.GetCurrentRow()

	case adddata.AddDataResultMsg:
		// По процессу условие всегда true.
		if msg.Err == nil {
			return tea.Batch(
				m.LoadData(),
				m.statusBar.NotifyOk(fmt.Sprintf("Added %s successfully", msg.Name)),
			)
		}
		return m.statusBar.NotifyError(fmt.Sprintf("Adding %s failed. See logs", msg.Name))

	case authorization.AuthorizationResultMsg:
		// По процессу условие всегда true.
		if msg.Err == nil {
			return m.LoadData()
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.UpDown):
			cmd = m.dataTable.Update(msg)
			m.dataDetail.Data = m.dataTable.GetCurrentRow()

		case key.Matches(msg, m.keyMap.Quit):
			return tea.Quit

		case key.Matches(msg, m.keyMap.DownloadBinary):
			current := m.dataTable.GetCurrentRow()
			if d, ok := current.(client.BinaryData); ok {
				return m.startDownloadBinary(d)
			}

		case key.Matches(msg, m.keyMap.Remove):
			current := m.dataTable.GetCurrentRow()
			return m.removeData(current)

		case key.Matches(msg, m.keyMap.AddLogin):
			return CallAddDataView(helper.DataTypeLogin)

		case key.Matches(msg, m.keyMap.AddNote):
			return CallAddDataView(helper.DataTypeNote)

		case key.Matches(msg, m.keyMap.AddBinary):
			return CallAddDataView(helper.DataTypeBinary)

		case key.Matches(msg, m.keyMap.AddCard):
			return CallAddDataView(helper.DataTypeCard)

		case key.Matches(msg, m.keyMap.Help):
			m.showHelp = !m.showHelp
		}
	}

	return tea.Batch(
		cmd,
		m.statusBar.Update(msg),
	)
}

func (m *Model) View() string {
	statusBar := m.statusBar.View()

	var helpView string
	h := m.Height - lipgloss.Height(statusBar)

	if m.showHelp {
		// Строка помощи
		hm := help.New()
		hm.ShowAll = true
		helpView = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingTop(1).
			Render(hm.View(m.keyMap))

		h -= lipgloss.Height(helpView)
	}

	// Окно со списком данных
	tableView := m.renderTableView(m.Width, h)

	// Окно детального просмотра
	detailViewWidth := m.Width - lipgloss.Width(tableView)
	detailView := m.renderDetailView(detailViewWidth, h)

	return lipgloss.JoinVertical(lipgloss.Top,
		removeEmptyStrings(
			lipgloss.JoinHorizontal(lipgloss.Top, tableView, detailView),
			statusBar,
			helpView,
		)...,
	)
}

func (m *Model) SetSize(width int, height int) {
	m.BaseModel.SetSize(width, height)
	m.statusBar.Width = width
}

func (m *Model) renderTableView(bubbleWidth int, height int) string {
	return helper.Borderize(
		"Data",
		m.dataTable.RenderInfoBar(),
		lipgloss.NewStyle().
			PaddingLeft(1).
			Render(m.dataTable.View()),
		bubbleWidth/3*2,
		height,
	)
}

func (m *Model) renderDetailView(width int, height int) string {
	return helper.Borderize(
		"Detail",
		"",
		lipgloss.NewStyle().
			Padding(0, 1). // TODO: при длинном сообщении нет паддинга на всех строках
			Render(m.dataDetail.View()),
		width,
		height,
	)
}

func (m *Model) LoadData() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		var result []client.Data

		ch := make(chan client.Data)
		var wg sync.WaitGroup

		wg.Go(func() {
			elems, err := m.loginService.GetAll(ctx)
			collect(elems, err, ch)
		})
		wg.Go(func() {
			elems, err := m.noteService.GetAll(ctx)
			collect(elems, err, ch)
		})
		wg.Go(func() {
			elems, err := m.binaryService.GetAll(ctx)
			collect(elems, err, ch)
		})
		wg.Go(func() {
			elems, err := m.cardService.GetAll(ctx)
			collect(elems, err, ch)
		})

		go func() {
			wg.Wait()
			close(ch)
		}()

		for el := range ch {
			result = append(result, el)
		}

		return loadDataMsg(result)
	}
}

func (m *Model) NotifyOk(format string, a ...any) tea.Cmd {
	return m.statusBar.NotifyOk(fmt.Sprintf(format, a...))
}

func (m *Model) NotifyError(format string, a ...any) tea.Cmd {
	return m.statusBar.NotifyError(fmt.Sprintf(format, a...))
}

func (m *Model) startDownloadBinary(data client.BinaryData) tea.Cmd {
	return func() tea.Msg {
		err := m.binaryService.Download(context.Background(), data.Name)
		if err != nil {
			return m.statusBar.NotifyError(fmt.Sprintf("Download %s failed: %v", data.Name, err))
		}

		return m.statusBar.NotifyOk(fmt.Sprintf("Downloaded %s successfully", data.Name))
	}
}

func (m *Model) removeData(data client.Data) tea.Cmd {
	return func() tea.Msg {
		var (
			ctx = context.Background()
			err error
		)

		switch data := data.(type) {
		case client.LoginData:
			err = m.loginService.Remove(ctx, data.Name)
		case client.NoteData:
			err = m.noteService.Remove(ctx, data.Name)
		case client.BinaryData:
			err = m.binaryService.Remove(ctx, data.Name)
		case client.CardData:
			err = m.cardService.Remove(ctx, data.Name)
		}

		if err != nil {
			return m.NotifyError("Removing %s failed: %v", data.GetName(), err)
		}

		return tea.Batch(
			m.NotifyOk("Removed %s successfully", data.GetName()),
			m.LoadData(),
		)()
	}
}

func removeEmptyStrings(strs ...string) []string {
	n := 0
	for _, s := range strs {
		if s != "" {
			strs[n] = s
			n++
		}
	}
	return strs[:n]
}

func collect[S ~[]E, E client.Data](s S, err error, out chan client.Data) {
	if err != nil {
		return
	}

	for _, el := range s {
		out <- el
	}
}
