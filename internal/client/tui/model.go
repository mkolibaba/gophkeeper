package tui

import (
	"context"
	"github.com/charmbracelet/bubbles/list"
	"github.com/mkolibaba/gophkeeper/internal/client"
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
}

// NewBubble создает новый экземпляр UI.
func NewBubble(
	loginService client.LoginService,
	noteService client.NoteService,
	binaryService client.BinaryService,
	cardService client.CardService,
) Bubble {
	loginTab := NewTab("Login", LoginFetcher(loginService))
	noteTab := NewTab("Note", NoteFetcher(noteService))
	binaryTab := NewTab("Binary", BinaryFetcher(binaryService))
	cardTab := NewTab("Card", CardFetcher(cardService))
	settingsTab := NewTab("Settings", func() []list.Item {
		return nil
	})

	return Bubble{
		tabs:          []TabItem{loginTab, noteTab, binaryTab, cardTab, settingsTab},
		loginService:  loginService,
		noteService:   noteService,
		binaryService: binaryService,
		cardService:   cardService,
	}
}

type Fetcher func() []list.Item

func LoginFetcher(loginService client.LoginService) Fetcher {
	return func() []list.Item {
		logins, err := loginService.GetAll(context.Background(), "demo")
		if err != nil {
			panic(err)
		}

		var loginItems []list.Item
		for _, login := range logins {
			loginItems = append(loginItems, ListItem{Name: login.Name, Desc: login.Login})
		}

		return loginItems
	}
}

func NoteFetcher(noteService client.NoteService) Fetcher {
	return func() []list.Item {
		notes, err := noteService.GetAll(context.Background(), "demo")
		if err != nil {
			panic(err)
		}

		var noteItems []list.Item
		for _, note := range notes {
			noteItems = append(noteItems, NewNoteItem(note.Name, note.Text))
		}

		return noteItems
	}
}

func BinaryFetcher(binaryService client.BinaryService) Fetcher {
	return func() []list.Item {
		binaries, err := binaryService.GetAll(context.Background(), "demo")
		if err != nil {
			panic(err)
		}

		var binaryItems []list.Item
		for _, binary := range binaries {
			binaryItems = append(binaryItems, ListItem{Name: binary.Name})
		}

		return binaryItems
	}
}

func CardFetcher(cardService client.CardService) Fetcher {
	return func() []list.Item {
		cards, err := cardService.GetAll(context.Background(), "demo")
		if err != nil {
			panic(err)
		}

		var cardItems []list.Item
		for _, card := range cards {
			cardItems = append(cardItems, NewCardItem(card.Name, card.Number))
		}

		return cardItems
	}
}

type TabItem struct {
	Name    string
	List    list.Model
	Fetcher Fetcher
}

func NewTab(name string, fetcher Fetcher) TabItem {
	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.SetShowStatusBar(false)

	return TabItem{
		Name:    name,
		List:    l,
		Fetcher: fetcher,
	}
}

func (t TabItem) UpdateItems() TabItem {
	t.List.SetItems(t.Fetcher())
	return t
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
