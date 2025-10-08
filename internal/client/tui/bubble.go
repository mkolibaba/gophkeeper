package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
	"github.com/mkolibaba/gophkeeper/internal/client"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/helper"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/state"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/view"
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

	view View

	views map[View]view.Model
}

func NewBubble(manager *state.Manager, session *client.Session) (Bubble, error) {
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
		manager: manager,
		session: session,
		views: map[View]view.Model{
			ViewAuthorization: view.InitialAuthorizationViewModel(manager),
			ViewMain:          view.InitialMainViewModel(manager),
			ViewAddData:       view.InitialAddDataViewModel(),
		},
	}, nil
}

// Init инициализирует UI.
func (b Bubble) Init() tea.Cmd {
	return tea.Batch(
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
			b.manager.SetInSession(msg.Login, msg.Password)
			b.view = ViewMain
			return b, b.manager.FetchData()
		}

	// Вызов окна добавления данных
	case view.AddDataCallMsg:
		b.view = ViewAddData
		return b, nil

	// Выход из окна добавления данных
	case view.ExitAddDataViewMsg:
		b.view = ViewMain
		return b, nil

	// Изменение размеров окна терминала
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height
		// TODO(trivial): можно ли как-то вычислить высоту компонента заголовка?
		for i := range b.views {
			b.views[i].SetSize(b.width, b.height-2) // -1 для заголовка приложения, -1 для футера
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
	footer := b.buildFooter()
	content := b.views[b.view].View()

	return lipgloss.JoinVertical(lipgloss.Top, title, content, footer)
}

func (b Bubble) buildFooter() string {
	// TODO(trivial): попробовать другие стили

	w := lipgloss.Width

	help := lipgloss.NewStyle().
		Width(8).
		PaddingLeft(1).
		Background(lipgloss.Color("243")).
		Render("h Help")

	var user string
	if usr := b.session.GetCurrentUser(); usr != nil {
		user = lipgloss.NewStyle().
			Width(w(usr.Login) + 2).
			PaddingLeft(1).
			Background(lipgloss.Color("243")).
			Render(usr.Login)
	}

	rest := lipgloss.NewStyle().
		Width(b.width - w(help) - w(user)).
		Background(lipgloss.Color("105")).
		Render()

	return lipgloss.JoinHorizontal(lipgloss.Top, help, rest, user)
}

// spew выводит в dump состояния объектов для дебага.
func (b Bubble) spew(a ...any) {
	spew.Fdump(b.dump, a...)
}
