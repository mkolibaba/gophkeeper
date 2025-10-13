package tuitest

import (
	"context"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/mkolibaba/gophkeeper/internal/client/inmem"
	"github.com/mkolibaba/gophkeeper/internal/client/mock"
	"github.com/mkolibaba/gophkeeper/internal/client/tui"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/view/adddata"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/view/authorization"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/view/home"
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
		AuthorizeFunc: func(ctx context.Context, login string, password string) error {
			return fmt.Errorf("some error")
		},
	}

	bubble, err := tui.NewBubble(tui.BubbleParams{
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
		AddDataView: adddata.New(adddata.Params{}),
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
	authMock.AuthorizeFunc = func(ctx context.Context, login string, password string) error {
		return nil
	}
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
