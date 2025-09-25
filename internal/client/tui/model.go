package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"regexp"
)

// Bubble представляет состояние UI.
type Bubble struct {
	tabs      []TabItem
	activeTab int

	width  int // Ширина терминала
	height int // Высота терминала
}

// NewBubble создает новый экземпляр UI.
func NewBubble(tabs ...TabItem) Bubble {
	return Bubble{
		tabs: tabs,
	}
}

type TabItem struct {
	Name string
	List list.Model
}

func NewTab(name string, items []list.Item) TabItem {
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.SetShowStatusBar(false)

	return TabItem{
		Name: name,
		List: l,
	}
}

// ListItem представляет элемент списка.
type ListItem struct {
	Name string
	Desc string
}

func (i ListItem) Title() string {
	return i.Name
}

func (i ListItem) Description() string {
	return i.Desc
}

func (i ListItem) FilterValue() string {
	return i.Name
}

type NoteItem struct {
	ListItem
}

func NewNoteItem(name, text string) NoteItem {
	return NoteItem{ListItem{name, text}}
}

func (i NoteItem) Description() string {
	maxLength := 30
	asRunes := []rune(i.Desc) // TODO: может есть лучше решение?
	if len(asRunes) > maxLength {
		return string(asRunes[:maxLength-3]) + "..."
	}
	return i.Desc
}

var maskingCardNumberRegexp = regexp.MustCompile(`(\d{6})\d{6}(\d{4})`)
var spacingCardNumberRegexp = regexp.MustCompile(`(.{4})(.{4})(.{4})(.{4})`)

// CardItem представляет элемент списка вкладки Card.
type CardItem struct {
	ListItem
}

func NewCardItem(name, number string) CardItem {
	return CardItem{ListItem{name, number}}
}

func (c CardItem) Description() string {
	masked := maskingCardNumberRegexp.ReplaceAllString(c.Desc, "$1******$2")
	return spacingCardNumberRegexp.ReplaceAllString(masked, "$1 $2 $3 $4")
}
