package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
	"github.com/mkolibaba/gophkeeper/internal/client"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/helper"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/state"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/view"
	"go.uber.org/fx"
	"io"
	"os"
)

// View представляет состояние UI.
type View int

const (
	// ViewAuthorization - авторизация.
	ViewAuthorization View = iota

	// ViewMain - основное окно приложения.
	ViewMain

	// ViewAddData - окно добавления данных.
	ViewAddData
)

// Bubble представляет корневой объект UI.
type Bubble struct {
	// Ширина терминала.
	width int
	// Высота терминала.
	height int

	// Writer для дебага.
	dump io.Writer

	manager *state.Manager
	session *client.Session

	// Текущий view интерфейса.
	view View

	views map[View]view.Model
}

type BubbleParams struct {
	fx.In

	Manager       *state.Manager
	Session       *client.Session
	BinaryService client.BinaryService
	LoginService  client.LoginService
	NoteService   client.NoteService
	CardService   client.CardService
}

func NewBubble(p BubbleParams) (Bubble, error) {
	var dump *os.File
	if dumpPath, ok := os.LookupEnv("SPEW_DUMP_OUTPUT"); ok {
		var err error
		dump, err = os.OpenFile(dumpPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			return Bubble{}, err
		}
	}

	return Bubble{
		dump:    dump,
		manager: p.Manager,
		session: p.Session,
		views: map[View]view.Model{
			ViewAuthorization: view.InitialAuthorizationViewModel(p.Manager),
			ViewMain:          view.InitialMainViewModel(p.Session, p.BinaryService),
			ViewAddData:       view.InitialAddDataViewModel(p.LoginService, p.NoteService),
		},
	}, nil
}

// Init инициализирует UI.
func (b Bubble) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("Gophkeeper"),
		b.views[ViewAuthorization].Init(),
	)
}

// Update обновляет UI в зависимости от события.
func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	b.spew(msg)

	// Корневые события
	switch msg := msg.(type) {
	// Авторизация
	case state.AuthorizationResultMsg:
		if msg.Err == nil {
			b.session.SetCurrentUser(client.User{Login: msg.Login, Password: msg.Password})
			b.view = ViewMain
			return b, b.manager.FetchData
		}

	// Вызов окна добавления данных
	case view.AddDataCallMsg:
		b.view = ViewAddData
		return b, b.views[ViewAddData].Update(msg)

	// Выход из окна добавления данных
	case view.ExitAddDataViewMsg:
		b.view = ViewMain
		return b, b.manager.FetchData // TODO: как будто немного неправильно

	// Изменение размеров окна терминала
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height
		// TODO(trivial): можно ли как-то вычислить высоту компонента заголовка?
		for i := range b.views {
			b.views[i].SetSize(b.width, b.height-1) // -1 для заголовка приложения
		}
		return b, nil
	}

	return b, b.views[b.view].Update(msg)
}

// View возвращает строковое представление UI.
func (b Bubble) View() string {
	title := helper.TitleStyle.
		Width(b.width).
		Render()
	content := b.views[b.view].View()

	return lipgloss.JoinVertical(lipgloss.Top, title, content)
}

// spew выводит в dump состояния объектов для дебага.
func (b Bubble) spew(a ...any) {
	spew.Fdump(b.dump, a...)
}
