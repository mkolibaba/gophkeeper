package helper

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
	"unicode/utf8"
)

func Borderize(topText, bottomText, content string) string {
	width := lipgloss.Width(content)
	height := lipgloss.Height(content)

	top := renderBorderTop(topText, width)
	bottom := renderBorderBottom(bottomText, width)
	middle := borderStyle.
		BorderTop(false).
		BorderBottom(false).
		Width(width).
		Height(height).
		Render(content)

	return lipgloss.JoinVertical(lipgloss.Left, top, middle, bottom)
}

// renderBorderTop отрисовывает верхнюю границу рамки заданной ширины width с текстом text.
func renderBorderTop(text string, width int) string {
	s := CustomBorderStyle // TODO: нужно брать из аргумента

	border, _, _, _, _ := borderStyle.GetBorder()
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

// renderBorderBottom отрисовывает нижнюю границу рамки заданной ширины width с текстом text.
func renderBorderBottom(text string, width int) string {
	// TODO: нужно брать style из агрументов
	s := CustomBorderStyle

	border, _, _, _, _ := borderStyle.GetBorder()
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
