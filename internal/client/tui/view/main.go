package view

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mkolibaba/gophkeeper/internal/client"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/detail"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/helper"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/state"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/table"
)

type mainViewKeyMap struct {
	UpDown         key.Binding
	AddData        key.Binding
	DownloadBinary key.Binding
	Quit           key.Binding
}

func (k mainViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.UpDown, k.Quit}
}

func (k mainViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.UpDown},
		{k.AddData},
		{k.DownloadBinary},
		{k.Quit},
	}
}

type MainViewModel struct {
	baseViewModel
	dataTable  table.Model
	dataDetail detail.Model
	keyMap     mainViewKeyMap
	manager    *state.Manager
}

func InitialMainViewModel(manager *state.Manager) *MainViewModel {
	keys := mainViewKeyMap{
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

	dataTable := table.NewModel()
	dataDetail := detail.NewModel()

	return &MainViewModel{
		dataTable:  dataTable,
		dataDetail: dataDetail,
		keyMap:     keys,
		manager:    manager,
	}
}

type AddDataCallMsg struct{}

func AddDataCall() tea.Msg {
	return AddDataCallMsg{}
}

func (m *MainViewModel) Init() tea.Cmd {
	return nil
}

func (m *MainViewModel) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case state.FetchDataMsg:
		m.dataTable, cmd = m.dataTable.Update(msg)
		current := m.dataTable.GetCurrentRow()
		m.dataDetail = m.dataDetail.SetData(current)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.UpDown):
			m.dataTable, cmd = m.dataTable.Update(msg)
			current := m.dataTable.GetCurrentRow()
			m.dataDetail = m.dataDetail.SetData(current)

		case key.Matches(msg, m.keyMap.Quit):
			return tea.Quit

		case key.Matches(msg, m.keyMap.DownloadBinary):
			current := m.dataTable.GetCurrentRow()
			if d, ok := current.(client.BinaryData); ok {
				m.manager.StartDownloadBinary(d)
			}

		case key.Matches(msg, m.keyMap.AddData):
			return AddDataCall
		}
	}

	return cmd
}

func (m *MainViewModel) View() string {
	// Строка помощи
	hm := help.New()
	hm.ShowAll = true
	helpView := lipgloss.NewStyle().PaddingLeft(1).Render(hm.View(m.keyMap))

	h := m.Height - lipgloss.Height(helpView)

	// Окно со списком данных
	tableView := m.renderTableView(m.Width, h)

	// Окно детального просмотра
	detailViewWidth := m.Width - lipgloss.Width(tableView)
	detailView := m.renderDetailView(detailViewWidth, h)

	return lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Top, tableView, detailView),
		helpView,
	)
}

func (m *MainViewModel) renderTableView(bubbleWidth int, height int) string {
	w := bubbleWidth/3*2 - helper.ContentStyle.GetHorizontalFrameSize()

	tableTopBorder := helper.RenderBorderTop(helper.ContentStyle, "Data", w)

	tableBottomBorder := helper.RenderBorderBottom(helper.ContentStyle, m.dataTable.RenderInfoBar(), w)

	tableView := helper.ContentStyle.
		BorderTop(false).
		BorderBottom(false).
		Width(w).
		Height(height - lipgloss.Height(tableTopBorder) - lipgloss.Height(tableBottomBorder)).
		PaddingLeft(1).
		Render(m.dataTable.View())

	return lipgloss.JoinVertical(lipgloss.Top, tableTopBorder, tableView, tableBottomBorder)
}

func (m *MainViewModel) renderDetailView(width int, height int) string {
	w := width - helper.ContentStyle.GetHorizontalFrameSize()

	detailTop := helper.RenderBorderTop(helper.ContentStyle, "Detail", w)

	detailView := helper.ContentStyle.
		BorderTop(false).
		Width(w).
		Height(height - helper.ContentStyle.GetBorderBottomSize() - lipgloss.Height(detailTop)).
		PaddingLeft(1).
		Render(m.dataDetail.View())

	return lipgloss.JoinVertical(lipgloss.Top, detailTop, detailView)
}
