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

func TestRegistrationView(t *testing.T) {
	t.Parallel()

	userService := inmem.NewUserService(log.New(io.Discard))
	authMock := &mock.AuthorizationServiceMock{
		RegisterFunc: func(ctx context.Context, login string, password string) (string, error) {
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
		AddDataView:  adddata.New(adddata.Params{}),
		EditDataView: editdata.New(editdata.Params{}),
		RegistrationView: registration.New(registration.Params{
			AuthorizationService: authMock,
			UserService:          userService,
		}),
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

	// Переходим на форму регистрации.
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlR})
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Registration") &&
			strings.Contains(s, "Repeat password")
	})

	// Пытаемся зарегистрироваться.
	tm.Type("foo")
	tm.Send(tea.KeyMsg{Type: tea.KeyTab})
	tm.Type("bar")
	tm.Send(tea.KeyMsg{Type: tea.KeyTab})
	tm.Type("bar")
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

	// Ждем отображения ошибки регистрации.
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "some error")
	})

	// Меняем моковый метод и пытаемся зарегистрироваться снова.
	authMock.RegisterFunc = func(ctx context.Context, login string, password string) (string, error) {
		return "super token", nil
	}
	userService.SetInfo("foo", "super token") // TODO: выглядит неправильно
	tm.Type("foo")
	tm.Send(tea.KeyMsg{Type: tea.KeyTab})
	tm.Type("bar")
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

	// TODO: по-хорошему нужно вынести в отдельный тест
	// Вызываем строку помощи.
	tm.Type("h")
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "alt+2")
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

func TestAddDataView_Binary(t *testing.T) {
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
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}, Alt: true})
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Add Binary")
	})
}

func TestEditDataView_Note(t *testing.T) {
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

func TestEditDataView_Card(t *testing.T) {
	t.Parallel()

	userService := inmem.NewUserService(log.New(io.Discard))
	authMock := &mock.AuthorizationServiceMock{
		AuthorizeFunc: func(ctx context.Context, login string, password string) (string, error) {
			return "some token", nil
		},
	}
	cardServiceMock := &mock.CardServiceMock{
		GetAllFunc: func(ctx context.Context) ([]client.CardData, error) {
			return []client.CardData{{
				ID:     1,
				Name:   "my card",
				Number: "112233",
			}}, nil
		},
		SaveFunc: func(ctx context.Context, data client.CardData) error {
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
			NoteService:   &mock.NoteServiceMock{},
			CardService:   cardServiceMock,
			UserService:   userService,
		}),
		AddDataView: adddata.New(adddata.Params{}),
		EditDataView: editdata.New(editdata.Params{
			CardService: cardServiceMock,
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
			strings.Contains(s, "my card")
	})

	// Заходим в форму редактирования.
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}, Alt: true})
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Edit my card")
	})

	// Заполняем форму.
	tm.Type(" new new")
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlS})
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Data") &&
			strings.Contains(s, "Detail")
	})

	// Проверяем корректность переданных данных
	cc := cardServiceMock.UpdateCalls()
	require.Len(t, cc, 1)
	c := cc[0].Data
	require.Equal(t, *c.Name, "my card new new")
}

func TestEditDataView_Login(t *testing.T) {
	t.Parallel()

	userService := inmem.NewUserService(log.New(io.Discard))
	authMock := &mock.AuthorizationServiceMock{
		AuthorizeFunc: func(ctx context.Context, login string, password string) (string, error) {
			return "some token", nil
		},
	}
	loginServiceMock := &mock.LoginServiceMock{
		GetAllFunc: func(ctx context.Context) ([]client.LoginData, error) {
			return []client.LoginData{{
				ID:    1,
				Name:  "my login123",
				Login: "testuser",
			}}, nil
		},
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
		AddDataView: adddata.New(adddata.Params{}),
		EditDataView: editdata.New(editdata.Params{
			LoginService: loginServiceMock,
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
			strings.Contains(s, "my login123")
	})

	// Заходим в форму редактирования.
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}, Alt: true})
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Edit my login123")
	})

	// Заполняем форму.
	tm.Type(" new new")
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlS})
	waitFor(t, tm, func(s string) bool {
		return strings.Contains(s, "Data") &&
			strings.Contains(s, "Detail")
	})

	// Проверяем корректность переданных данных
	cc := loginServiceMock.UpdateCalls()
	require.Len(t, cc, 1)
	c := cc[0].Data
	require.Equal(t, *c.Name, "my login123 new new")
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
