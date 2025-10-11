package helper

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
	"unicode/utf8"
)

func Borderize(style lipgloss.Style, topText, bottomText, content string) string {
	width := style.GetWidth()

	top := renderBorderTop(style, topText, width)
	bottom := renderBorderBottom(style, bottomText, width)
	middle := style.
		BorderTop(false).
		BorderBottom(false).
		//Height(style.GetHeight() - lipgloss.Height(top) - lipgloss.Height(bottom)). // TODO: сделать это здесь
		Height(style.GetHeight()).
		Render(content)

	return lipgloss.JoinVertical(lipgloss.Left, top, middle, bottom)
}

// renderBorderTop отрисовывает верхнюю границу рамки заданной ширины width для
// стиля style с текстом text.
func renderBorderTop(style lipgloss.Style, text string, width int) string {
	s := CustomBorderStyle // TODO: нужно брать из аргумента

	border, _, _, _, _ := style.GetBorder()
	borderLeft := border.TopLeft
	borderMiddle := border.Top
	borderRight := border.TopRight

	if text == "" {
		return s.Render(borderLeft + strings.Repeat(borderMiddle, width) + borderRight)
	}

	leftText := borderLeft + borderRight

	rightRepeat := width -
		utf8.RuneCountInString(leftText) -
		utf8.RuneCountInString(text)
	// TODO: возможно, это workaround из-за некорректной логики
	if rightRepeat <= 0 {
		rightRepeat = 1
	}

	rightText := borderLeft + strings.Repeat(borderMiddle, rightRepeat) + borderRight

	left := s.Render(leftText)
	right := s.Render(rightText)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, text, right)
}

// renderBorderBottom отрисовывает нижнюю границу рамки заданной ширины width для
// стиля style с текстом text.
func renderBorderBottom(style lipgloss.Style, text string, width int) string {
	// TODO: нужно брать style из агрументов
	s := CustomBorderStyle

	border, _, _, _, _ := style.GetBorder()
	borderLeft := border.BottomLeft
	borderMiddle := border.Bottom
	borderRight := border.BottomRight

	if text == "" {
		return s.Render(borderLeft + strings.Repeat(borderMiddle, width) + borderRight)
	}

	rightText := borderLeft + borderRight

	leftRepeat := width -
		utf8.RuneCountInString(rightText) -
		utf8.RuneCountInString(text)
	// TODO: возможно, это workaround из-за некорректной логики
	if leftRepeat <= 0 {
		leftRepeat = 1
	}

	leftText := borderLeft + strings.Repeat(borderMiddle, leftRepeat) + borderRight

	left := s.Render(leftText)
	right := s.Render(rightText)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, text, right)
}
