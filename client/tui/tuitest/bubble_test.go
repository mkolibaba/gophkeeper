package tuitest

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/mkolibaba/gophkeeper/client"
	"github.com/mkolibaba/gophkeeper/client/inmem"
	"github.com/mkolibaba/gophkeeper/client/mock"
	"github.com/mkolibaba/gophkeeper/client/tui"
	"github.com/mkolibaba/gophkeeper/client/tui/view/adddata"
	"github.com/mkolibaba/gophkeeper/client/tui/view/authorization"
	"github.com/mkolibaba/gophkeeper/client/tui/view/editdata"
	"github.com/mkolibaba/gophkeeper/client/tui/view/home"
	"github.com/mkolibaba/gophkeeper/client/tui/view/registration"
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"testing"
	"time"
)

func TestAuthorizationView(t *testing.T) {
	t.Parallel()

	userService := inmem.NewUserService(log.New(io.Discard))
	authMock := &mock.AuthorizationServiceMock{
		AuthorizeFunc: func(ctx context.Context, login string, password string) (string, error) {
			return "", fmt.Errorf("some error")
		},
	}
	var config client.Config
	config.Development.Enabled = false

	bubble, err := tui.NewBubble(tui.BubbleParams{
		Config: &config, // TODO: выглядит как сильная связанность
		AuthorizationView: authorization.New(authorization.Params{
			AuthorizationService: authMock,
			UserService:          userService,
		}),
		MainView: home.New(home.Params{
			LoginService:  &mock.LoginServiceMock{},
			BinaryService: &mock.BinaryServiceMock{},
			NoteService:   &mock.NoteServiceMock{},
			CardService:   &mock.CardServiceMock{},
			UserService:   userService,
		}),
		AddDataView:      adddata.New(adddata.Params{}),
		EditDataView:     editdata.New(editdata.Params{}),
		RegistrationView: registration.New(registration.Params{}),
	})
	require.NoError(t, err)

	// Инициализируем приложение.
	tm := teatest.NewTestModel(t, bubble, teatest.WithInitialTermSize(130, 40))

	// Ожидаем отрисовки формы авторизации.
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Authorization") &&
			strings.Contains(s, "Login") &&
			strings.Contains(s, "Password")
	})

	// Пытаемся авторизоваться.
	tm.Type("foo")
	tm.Send(tea.KeyMsg{Type: tea.KeyTab})
	tm.Type("bar")
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

	// Ждем отображения ошибки авторизации.
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "some error")
	})

	// Меняем моковый метод и пытаемся авторизоваться снова.
	authMock.AuthorizeFunc = func(ctx context.Context, login string, password string) (string, error) {
		return "super token", nil
	}
	userService.SetInfo("foo", "super token") // TODO: выглядит неправильно
	tm.Type("foo")
	tm.Send(tea.KeyMsg{Type: tea.KeyTab})
	tm.Type("bar")
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

	// Проверяем, что приложение перешло на домашнюю страницу.
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Data") &&
			strings.Contains(s, "Detail")
	})
}

func TestAddDataView(t *testing.T) {
	t.Parallel()

	userService := inmem.NewUserService(log.New(io.Discard))
	authMock := &mock.AuthorizationServiceMock{
		AuthorizeFunc: func(ctx context.Context, login string, password string) (string, error) {
			return "some token", nil
		},
	}
	noteServiceMock := &mock.NoteServiceMock{
		SaveFunc: func(ctx context.Context, data client.NoteData) error {
			return nil
		},
	}
	var config client.Config
	config.Development.Enabled = false

	bubble, err := tui.NewBubble(tui.BubbleParams{
		Config: &config, // TODO: выглядит как сильная связанность
		AuthorizationView: authorization.New(authorization.Params{
			AuthorizationService: authMock,
			UserService:          userService,
		}),
		MainView: home.New(home.Params{
			LoginService:  &mock.LoginServiceMock{},
			BinaryService: &mock.BinaryServiceMock{},
			NoteService:   noteServiceMock,
			CardService:   &mock.CardServiceMock{},
			UserService:   userService,
		}),
		AddDataView: adddata.New(adddata.Params{
			NoteService: noteServiceMock,
		}),
		EditDataView:     editdata.New(editdata.Params{}),
		RegistrationView: registration.New(registration.Params{}),
	})
	require.NoError(t, err)

	// Инициализируем приложение.
	tm := teatest.NewTestModel(t, bubble, teatest.WithInitialTermSize(130, 40))

	// Ожидаем отрисовки формы авторизации.
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Authorization") &&
			strings.Contains(s, "Login") &&
			strings.Contains(s, "Password")
	})

	// За счет мока сразу авторизуемся.
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

	// Проверяем, что приложение перешло на домашнюю страницу.
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Data") &&
			strings.Contains(s, "Detail")
	})

	// Заходим в форму добавления.
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}, Alt: true})
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Add Note")
	})

	// Заполняем форму.
	tm.Type("new name")
	tm.Send(tea.KeyMsg{Type: tea.KeyTab})
	tm.Type("some long long long text")
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlS})
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Data") &&
			strings.Contains(s, "Detail")
	})

	// Проверяем корректность переданных данных
	cc := noteServiceMock.SaveCalls()
	require.Len(t, cc, 1)
	c := cc[0].Data
	require.Equal(t, c.Name, "new name")
	require.Equal(t, c.Text, "some long long long text")
}

func TestEditDataView(t *testing.T) {
	t.Parallel()

	userService := inmem.NewUserService(log.New(io.Discard))
	authMock := &mock.AuthorizationServiceMock{
		AuthorizeFunc: func(ctx context.Context, login string, password string) (string, error) {
			return "some token", nil
		},
	}
	noteServiceMock := &mock.NoteServiceMock{
		GetAllFunc: func(ctx context.Context) ([]client.NoteData, error) {
			return []client.NoteData{{
				ID:   1,
				Name: "my note",
				Text: "long text",
			}}, nil
		},
		SaveFunc: func(ctx context.Context, data client.NoteData) error {
			return nil
		},
	}
	var config client.Config
	config.Development.Enabled = false

	bubble, err := tui.NewBubble(tui.BubbleParams{
		Config: &config, // TODO: выглядит как сильная связанность
		AuthorizationView: authorization.New(authorization.Params{
			AuthorizationService: authMock,
			UserService:          userService,
		}),
		MainView: home.New(home.Params{
			LoginService:  &mock.LoginServiceMock{},
			BinaryService: &mock.BinaryServiceMock{},
			NoteService:   noteServiceMock,
			CardService:   &mock.CardServiceMock{},
			UserService:   userService,
		}),
		AddDataView: adddata.New(adddata.Params{}),
		EditDataView: editdata.New(editdata.Params{
			NoteService: noteServiceMock,
		}),
		RegistrationView: registration.New(registration.Params{}),
	})
	require.NoError(t, err)

	// Инициализируем приложение.
	tm := teatest.NewTestModel(t, bubble, teatest.WithInitialTermSize(130, 40))

	// Ожидаем отрисовки формы авторизации.
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Authorization") &&
			strings.Contains(s, "Login") &&
			strings.Contains(s, "Password")
	})

	// За счет мока сразу авторизуемся.
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

	// Проверяем, что приложение перешло на домашнюю страницу и отрисовало таблицу.
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Data") &&
			strings.Contains(s, "Detail") &&
			strings.Contains(s, "my note")
	})

	// Заходим в форму редактирования.
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}, Alt: true})
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Edit my note")
	})

	// Заполняем форму.
	tm.Type(" new new")
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlS})
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Data") &&
			strings.Contains(s, "Detail")
	})

	// Проверяем корректность переданных данных
	cc := noteServiceMock.UpdateCalls()
	require.Len(t, cc, 1)
	c := cc[0].Data
	require.Equal(t, *c.Name, "my note new new")
}

func waitFor(t *testing.T, tm *teatest.TestModel, cond func(s string) bool) {
	t.Helper()

	teatest.WaitFor(
		t,
		tm.Output(),
		func(b []byte) bool {
			return cond(string(b))
		},
		teatest.WithCheckInterval(100*time.Millisecond),
		teatest.WithDuration(3*time.Second),
	)
}
