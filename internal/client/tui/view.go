package tui

import (
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
)

// View возвращает строковое представление UI.
func (b Bubble) View() string {
	var doc strings.Builder

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

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(windowStyle.Width(lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize()).Render(b.tabs[b.activeTab].List.View()))
	return docStyle.Render(doc.String())
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}
