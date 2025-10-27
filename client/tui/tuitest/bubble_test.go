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
	loginServiceMock := &mock.LoginServiceMock{
		SaveFunc: func(ctx context.Context, data client.LoginData) error {
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
			LoginService:  loginServiceMock,
			BinaryService: &mock.BinaryServiceMock{},
			NoteService:   &mock.NoteServiceMock{},
			CardService:   &mock.CardServiceMock{},
			UserService:   userService,
		}),
		AddDataView: adddata.New(adddata.Params{
			LoginService: loginServiceMock,
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

	// Заходим в форму добавления логина.
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}, Alt: true})
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Add Login")
	})

	// Заполняем форму.
	tm.Type("new name")
	tm.Send(tea.KeyMsg{Type: tea.KeyTab})
	tm.Type("some login")
	tm.Send(tea.KeyMsg{Type: tea.KeyTab})
	tm.Type("my password")
	tm.Send(tea.KeyMsg{Type: tea.KeyTab})
	tm.Type("example.com")
	tm.Send(tea.KeyMsg{Type: tea.KeyTab})
	tm.Type("long note")
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlS})
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Data") &&
			strings.Contains(s, "Detail")
	})

	// Проверяем корректность переданных данных
	cc := loginServiceMock.SaveCalls()
	require.Len(t, cc, 1)
	c := cc[0].Data
	require.Equal(t, c.Name, "new name")
	require.Equal(t, c.Login, "some login")
	require.Equal(t, c.Password, "my password")
	require.Equal(t, c.Website, "example.com")
	require.Equal(t, c.Notes, "long note")
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
