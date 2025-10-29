package helper

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func Borderize(topText, bottomText, content string, width, height int) string {
	top := renderBorderTop(topText, width)
	bottom := renderBorderBottom(bottomText, width)
	middle := borderStyle.
		BorderTop(false).
		BorderBottom(false).
		Width(width - 2).
		Height(height - 2).
		Render(content)

	return lipgloss.JoinVertical(lipgloss.Left, top, middle, bottom)
}

// renderBorderTop отрисовывает верхнюю границу рамки заданной ширины width с текстом text.
func renderBorderTop(text string, width int) string {
	border, _, _, _, _ := borderStyle.GetBorder()
	borderLeft := border.TopLeft
	borderMiddle := border.Top
	borderRight := border.TopRight

	l := lipgloss.Width

	left := customBorderStyle.Render(borderLeft)
	enclosed := encloseInBrackets(text, borderRight, borderLeft)
	right := customBorderStyle.Render(borderRight)
	remaining := customBorderStyle.Render(strings.Repeat(borderMiddle, max(0, width-l(left)-l(right)-l(enclosed))))

	return lipgloss.JoinHorizontal(lipgloss.Top, left, enclosed, remaining, right)
}

// renderBorderBottom отрисовывает нижнюю границу рамки заданной ширины width с текстом text.
func renderBorderBottom(text string, width int) string {
	border, _, _, _, _ := borderStyle.GetBorder()
	borderLeft := border.BottomLeft
	borderMiddle := border.Bottom
	borderRight := border.BottomRight

	l := lipgloss.Width

	left := customBorderStyle.Render(borderLeft)
	enclosed := encloseInBrackets(text, borderRight, borderLeft)
	right := customBorderStyle.Render(borderRight)
	remaining := customBorderStyle.Render(strings.Repeat(borderMiddle, max(0, width-l(left)-l(right)-l(enclosed))))

	return lipgloss.JoinHorizontal(lipgloss.Top, left, remaining, enclosed, right)
}

func encloseInBrackets(text, left, right string) string {
	if text != "" {
		return customBorderStyle.Render(left) + text + customBorderStyle.Render(right)
	}
	return ""
}
