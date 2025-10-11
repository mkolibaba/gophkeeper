package helper

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
	"unicode/utf8"
)

// RenderBorderTop отрисовывает верхнюю границу рамки заданной ширины width для
// стиля style с текстом text.
func RenderBorderTop(style lipgloss.Style, text string, width int) string {
	border, _, _, _, _ := style.GetBorder()
	borderLeft := border.TopLeft
	borderMiddle := border.Top
	borderRight := border.TopRight

	leftText := borderLeft + borderRight

	rightRepeat := width -
		utf8.RuneCountInString(leftText) -
		utf8.RuneCountInString(text)
	// TODO: возможно, это workaround из-за некорректной логики
	if rightRepeat <= 0 {
		rightRepeat = 1
	}

	rightText := borderLeft + strings.Repeat(borderMiddle, rightRepeat) + borderRight

	//s := style.
	//	UnsetBorderBottom().
	//	UnsetBorderTop().
	//	UnsetBorderLeft().
	//	UnsetBorderRight()
	s := CustomBorderStyle
	left := s.Render(leftText)
	right := s.Render(rightText)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, text, right)
}

// RenderBorderBottom отрисовывает нижнюю границу рамки заданной ширины width для
// стиля style с текстом text.
func RenderBorderBottom(style lipgloss.Style, text string, width int) string {
	border, _, _, _, _ := style.GetBorder()
	borderLeft := border.BottomLeft
	borderMiddle := border.Bottom
	borderRight := border.BottomRight

	rightText := borderLeft + borderRight

	leftRepeat := width -
		utf8.RuneCountInString(rightText) -
		utf8.RuneCountInString(text)
	// TODO: возможно, это workaround из-за некорректной логики
	if leftRepeat <= 0 {
		leftRepeat = 1
	}

	leftText := borderLeft + strings.Repeat(borderMiddle, leftRepeat) + borderRight

	//s := style.
	//	UnsetBorderBottom().
	//	UnsetBorderTop().
	//	UnsetBorderLeft().
	//	UnsetBorderRight()
	// TODO: нужно брать style из агрументов
	s := CustomBorderStyle
	left := s.Render(leftText)
	right := s.Render(rightText)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, text, right)
}
