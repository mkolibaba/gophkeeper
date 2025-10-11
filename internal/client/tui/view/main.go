package view

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
	"github.com/mkolibaba/gophkeeper/internal/client/tui/helper"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/orchestrator"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/table"
)

type mainViewKeyMap struct {
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

func (k mainViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.UpDown, k.Quit}
}

func (k mainViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.UpDown},
		{k.AddLogin, k.AddNote, k.AddBinary, k.AddCard},
		{k.DownloadBinary, k.Remove},
		{k.Quit},
	}
}

type MainViewModel struct {
	baseViewModel
	dataTable     table.Model
	dataDetail    detail.Model
	keyMap        mainViewKeyMap
	session       *client.Session
	showHelp      bool
	statusBar     *statusbar.Model
	orchestrator  *orchestrator.Orchestrator
	binaryService client.BinaryService
}

func InitialMainViewModel(
	session *client.Session,
	binaryService client.BinaryService,
	orchestrator *orchestrator.Orchestrator,
) *MainViewModel {
	keys := mainViewKeyMap{
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

	return &MainViewModel{
		dataTable:     dataTable,
		dataDetail:    dataDetail,
		statusBar:     statusBar,
		keyMap:        keys,
		session:       session,
		orchestrator:  orchestrator,
		binaryService: binaryService,
	}
}

type AddDataCallMsg struct {
	t DataType
}

func AddDataCall(t DataType) tea.Cmd {
	return func() tea.Msg {
		return AddDataCallMsg{
			t: t,
		}
	}
}

func (m *MainViewModel) Init() tea.Cmd {
	return nil
}

func (m *MainViewModel) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case helper.LoadDataMsg:
		m.statusBar.CurrentUser = m.session.GetCurrentUser().Login // TODO: это хак, сделать лучше
		m.dataTable, cmd = m.dataTable.Update(msg)
		m.dataDetail.Data = m.dataTable.GetCurrentRow()

	case statusbar.NotificationMsg:
		return m.statusBar.Update(msg)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.UpDown):
			m.dataTable, cmd = m.dataTable.Update(msg)
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
			return AddDataCall(DataTypeLogin)

		case key.Matches(msg, m.keyMap.AddNote):
			return AddDataCall(DataTypeNote)

		case key.Matches(msg, m.keyMap.AddBinary):
			return AddDataCall(DataTypeBinary)

		case key.Matches(msg, m.keyMap.AddCard):
			return AddDataCall(DataTypeCard)

		case key.Matches(msg, m.keyMap.Help):
			m.showHelp = !m.showHelp
		}
	}

	return cmd
}

func (m *MainViewModel) View() string {
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

func (m *MainViewModel) SetSize(width int, height int) {
	m.baseViewModel.SetSize(width, height)
	m.statusBar.Width = width
}

func (m *MainViewModel) renderTableView(bubbleWidth int, height int) string {
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

func (m *MainViewModel) renderDetailView(width int, height int) string {
	return helper.Borderize(
		"Detail",
		"",
		lipgloss.NewStyle().
			Padding(0, 1).
			Render(m.dataDetail.View()),
		width,
		height,
	)
}

func (m *MainViewModel) startDownloadBinary(data client.BinaryData) tea.Cmd {
	return func() tea.Msg {
		err := m.binaryService.Download(context.Background(), data.Name)
		if err != nil {
			return statusbar.NotifyError(fmt.Sprintf("Download %s failed: %v", data.Name, err))
		}

		return statusbar.NotifyOk(fmt.Sprintf("Downloaded %s successfully", data.Name))
	}
}

func (m *MainViewModel) removeData(data client.Data) tea.Cmd {
	return func() tea.Msg {
		err := m.orchestrator.Remove(context.Background(), data)
		if err != nil {
			return statusbar.NotifyError(fmt.Sprintf("Removing %s failed: %v", data.GetName(), err))
		}

		return tea.Sequence(
			statusbar.NotifyOk(fmt.Sprintf("Removed %s successfully", data.GetName())),
			helper.LoadData(m.orchestrator.GetAll(context.Background())),
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
