package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/mkolibaba/gophkeeper/internal/client"
	"go.uber.org/fx"
	"regexp"
)

// Bubble представляет состояние UI.
type Bubble struct {
	tabs      []TabItem
	activeTab int

	width  int // Ширина терминала
	height int // Высота терминала

	loginService  client.LoginService
	noteService   client.NoteService
	binaryService client.BinaryService
	cardService   client.CardService

	shutdowner fx.Shutdowner
}

// NewBubble создает новый экземпляр UI.
func NewBubble(
	loginService client.LoginService,
	noteService client.NoteService,
	binaryService client.BinaryService,
	cardService client.CardService,
	shutdowner fx.Shutdowner,
) Bubble {
	loginTab := NewTab("Login", []list.Item{
		ListItem{Name: "Google", Desc: "iivanov"},
		ListItem{Name: "Ozon", Desc: "+79031002030"},
		ListItem{Name: "Wildberries", Desc: "+79031002030"},
		ListItem{Name: "Госуслуги", Desc: "iivanov@gmail.com"},
		ListItem{Name: "Mail.ru", Desc: "ivanivanov"},
		ListItem{Name: "VK", Desc: "ivanivanov@mail.ru"},
	})
	noteTab := NewTab("Note", []list.Item{
		NewNoteItem("Записки о природе", "Кто никогда не видал, как растет клюква, тот может очень долго идти по болоту и не замечать, что он по клюкве идет."),
		NewNoteItem("Мысль", "Живешь ты, может быть, сам триста лет, и кто породил тебя, тот в яичке своем пересказал все, что он тоже узнал за свои триста лет жизни."),
	})
	binaryTab := NewTab("Binary", []list.Item{})
	cardTab := NewTab("Card", []list.Item{
		NewCardItem("Сбербанк", "2200123456789019"),
		NewCardItem("Т-Банк", "2201987654321000"),
	})
	settingsTab := NewTab("Settings", []list.Item{})

	return Bubble{
		tabs:          []TabItem{loginTab, noteTab, binaryTab, cardTab, settingsTab},
		loginService:  loginService,
		noteService:   noteService,
		binaryService: binaryService,
		cardService:   cardService,
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
