package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var (
	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)

	highlightColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	inactiveTabStyle  = lipgloss.NewStyle().
				Border(inactiveTabBorder, true).
				BorderForeground(highlightColor).
				Padding(0, 1).
				Width(20)

	activeTabBorder = tabBorderWithBottom("┘", " ", "└")
	activeTabStyle  = inactiveTabStyle.Border(activeTabBorder, true)

	windowStyle = lipgloss.NewStyle().
			BorderForeground(highlightColor).
			Padding(2, 0).
			Border(lipgloss.NormalBorder()).
			UnsetBorderTop()

	bold = lipgloss.NewStyle().Bold(true).
		Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"})
	border = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(highlightColor)
)

// View возвращает строковое представление UI.
func (b Bubble) View() string {
	var renderedTabs []string

	for i, t := range b.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(b.tabs)-1, i == b.activeTab
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t.Name))
	}

	// Левое окошко

	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	var doc strings.Builder
	doc.WriteString(tabRow)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Width(lipgloss.Width(tabRow) - windowStyle.GetHorizontalFrameSize()).
		Render(b.tabs[b.activeTab].List.View()))
	listView := doc.String()

	// Правое окошко

	var item list.DefaultItem = NewCardItem("1", "2")
	if b.tabs[b.activeTab].List.SelectedItem() != nil {
		item = b.tabs[b.activeTab].List.SelectedItem().(list.DefaultItem)
	}
	infoViewContent := lipgloss.JoinVertical(lipgloss.Top,
		bold.Render("Name"),
		item.Title(),
		"",
		bold.Render("Login"),
		item.Description(),
		"",
		bold.Render("Password"),
		"some password",
		"",
		bold.Render("meta 1"),
		"some meta 1",
	)

	infoView := border.
		Width(b.width - lipgloss.Width(listView) - border.GetHorizontalFrameSize() - 2).
		Height(b.height - 3).
		Render(infoViewContent)

	ui := lipgloss.JoinHorizontal(lipgloss.Top, listView, infoView)

	return docStyle.Render(ui)
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}
