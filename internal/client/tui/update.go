package tui

import tea "github.com/charmbracelet/bubbletea"

// Update обновляет UI в зависимости от события.
func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return b, tea.Quit
		case "right", "l", "n", "tab":
			b.activeTab = min(b.activeTab+1, len(b.tabs)-1)
			b.tabs[b.activeTab] = b.tabs[b.activeTab].UpdateItems()
			return b, nil // блокируем дефолтное обновление list
		case "left", "h", "p", "shift+tab":
			b.activeTab = max(b.activeTab-1, 0)
			b.tabs[b.activeTab] = b.tabs[b.activeTab].UpdateItems()
			return b, nil // блокируем дефолтное обновление list
		}
	case tea.WindowSizeMsg:
		b.height, b.width = docStyle.GetFrameSize()
		for i := range b.tabs {
			b.tabs[i].List.SetSize(msg.Width-b.height, msg.Height-b.width-7)
		}
	}

	var cmd tea.Cmd
	b.tabs[b.activeTab] = b.tabs[b.activeTab].UpdateItems()
	b.tabs[b.activeTab].List, cmd = b.tabs[b.activeTab].List.Update(msg)
	return b, cmd
}
