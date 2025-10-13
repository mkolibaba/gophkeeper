package authorization

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mkolibaba/gophkeeper/internal/client"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/components/inputset"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/helper"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/view"
	"go.uber.org/fx"
)

type Model struct {
	view.BaseModel
	inputSet             *inputset.Model
	authorizationService client.AuthorizationService
	userService          client.UserService
}

type Params struct {
	fx.In

	AuthorizationService client.AuthorizationService
	UserService          client.UserService
}

func New(p Params) *Model {
	return &Model{
		authorizationService: p.AuthorizationService,
		userService:          p.UserService,
		inputSet: inputset.NewInputSet(
			inputset.NewTextInput("Login"),
			inputset.NewTextInput("Password", inputset.WithEchoModePassword()),
		),
	}
}

func (m *Model) Init() tea.Cmd {
	return m.inputSet.Init()
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case AuthorizationResultMsg:
		m.inputSet.Err = msg.Err
		m.inputSet.Reset()

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "esc":
			return tea.Quit

		case "up", "down", "tab":
			return m.inputSet.Update(msg)

		case "enter":
			return m.authorize()
		}
	}

	return m.inputSet.Update(msg)
}

func (m *Model) View() string {
	return helper.Borderize(
		"Authorization",
		"",
		lipgloss.NewStyle().
			PaddingTop(1).
			PaddingLeft(1).
			Render(m.inputSet.View()),
		m.Width,
		m.Height/2,
	)
}

type AuthorizationResultMsg struct {
	Err error
}

func (m *Model) authorize() tea.Cmd {
	values := m.inputSet.Values()
	return func() tea.Msg {
		login, password := values["Login"], values["Password"]
		err := m.authorizationService.Authorize(context.Background(), login, password)
		if err == nil {
			m.userService.Set(client.User{
				Login:    login,
				Password: password,
			})
		}

		return AuthorizationResultMsg{
			Err: err,
		}
	}
}
