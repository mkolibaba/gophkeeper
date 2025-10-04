package view

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/helper"
)

type addViewKeyMap struct {
	Exit key.Binding
}

func (k addViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Exit}
}

func (k addViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Exit},
	}
}

type AddDataViewModel struct {
	baseViewModel
	keyMap addViewKeyMap
}

func InitialAddDataViewModel() *AddDataViewModel {
	keyMap := addViewKeyMap{
		Exit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "exit"),
		),
	}

	return &AddDataViewModel{
		keyMap: keyMap,
	}
}

func (m *AddDataViewModel) Init() tea.Cmd {
	return nil
}

func (m *AddDataViewModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Exit):
			return ExitAddDataView
		}
	}

	return nil
}

func (m *AddDataViewModel) View() string {
	// Строка помощи
	hm := help.New()
	hm.ShowAll = true
	helpView := lipgloss.NewStyle().PaddingLeft(1).Render(hm.View(m.keyMap))

	w := m.Width - helper.ContentStyle.GetHorizontalFrameSize()

	borderTop := helper.RenderBorderTop(helper.ContentStyle, "Add Data", w)

	h := m.Height - lipgloss.Height(borderTop) - helper.ContentStyle.GetBorderBottomSize() - lipgloss.Height(helpView)

	addDataView := helper.ContentStyle.
		BorderTop(false).
		Width(w).
		Height(h).
		PaddingLeft(1).
		Render("Тут можно будет добавлять всякий контент")

	return lipgloss.JoinVertical(lipgloss.Top, borderTop, addDataView, helpView)
}
