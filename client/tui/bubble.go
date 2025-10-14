package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
	"github.com/mkolibaba/gophkeeper/client/tui/helper"
	"github.com/mkolibaba/gophkeeper/client/tui/view"
	"github.com/mkolibaba/gophkeeper/client/tui/view/adddata"
	"github.com/mkolibaba/gophkeeper/client/tui/view/authorization"
	"github.com/mkolibaba/gophkeeper/client/tui/view/home"
	"go.uber.org/fx"
	"io"
	"os"
)

// Bubble представляет корневой объект UI.
type Bubble struct {
	// Ширина терминала.
	width int
	// Высота терминала.
	height int

	// Writer для дебага.
	dump io.Writer

	// Текущий view интерфейса.
	view view.View

	// Все возможные view интерфейса.
	views map[view.View]view.Model
}

type BubbleParams struct {
	fx.In

	AuthorizationView *authorization.Model
	MainView          *home.Model
	AddDataView       *adddata.Model
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
		dump: dump,
		views: map[view.View]view.Model{
			view.ViewAuthorization: p.AuthorizationView,
			view.ViewHome:          p.MainView,
			view.ViewAddData:       p.AddDataView,
		},
	}, nil
}

// Init инициализирует UI.
func (b Bubble) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("Gophkeeper"),
		b.views[view.ViewAuthorization].Init(),
	)
}

// Update обновляет UI в зависимости от события.
func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	b.spew(msg)

	// Корневые события
	switch msg := msg.(type) {
	// Авторизация
	case authorization.AuthorizationResultMsg:
		if msg.Err == nil {
			b.view = view.ViewHome
		}

	// Добавление данных
	case adddata.AddDataResultMsg:
		if msg.Err == nil {
			b.view = view.ViewHome
			homeView := b.views[view.ViewHome].(*home.Model)
			return b, tea.Batch(
				homeView.LoadData(),
				homeView.NotifyOk("Added %s successfully", msg.Name),
			)
		}

	// Вызов окна добавления данных
	case home.CallAddDataViewMsg:
		b.view = view.ViewAddData
		addDataView := b.views[view.ViewAddData].(*adddata.Model)
		addDataView.ResetFor(helper.DataType(msg))
		return b, addDataView.Init()

	// Выход из окна добавления данных
	case adddata.ExitMsg:
		b.view = view.ViewHome

	// Изменение размеров окна терминала
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height
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
	if b.dump != nil {
		spew.Fdump(b.dump, a...)
	}
}
