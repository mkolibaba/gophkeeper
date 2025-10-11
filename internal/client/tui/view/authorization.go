package view

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mkolibaba/gophkeeper/internal/client"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/components/inputset"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/helper"
)

type AuthorizationViewModel struct {
	baseViewModel
	inputSet             *inputset.Model
	authorizationService client.AuthorizationService
}

func InitialAuthorizationViewModel(authorizationService client.AuthorizationService) *AuthorizationViewModel {
	return &AuthorizationViewModel{
		authorizationService: authorizationService,
		inputSet: inputset.NewInputSet(
			inputset.NewTextInput("Login"),
			inputset.NewTextInput("Password", inputset.WithEchoModePassword()),
		),
	}
}

func (m *AuthorizationViewModel) Init() tea.Cmd {
	return m.inputSet.Init()
}

func (m *AuthorizationViewModel) Update(msg tea.Msg) tea.Cmd {
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

func (m *AuthorizationViewModel) View() string {
	w := m.Width - helper.ContentStyle.GetHorizontalFrameSize()

	borderTop := helper.RenderBorderTop(helper.ContentStyle, "Authorization", w)

	h := m.Height - lipgloss.Height(borderTop) - helper.ContentStyle.GetBorderBottomSize()

	authorizationView := helper.ContentStyle.
		BorderTop(false).
		Width(w).
		Height(h / 2).
		PaddingTop(1).
		PaddingLeft(1).
		Render(m.inputSet.View())

	return lipgloss.JoinVertical(lipgloss.Top, borderTop, authorizationView)
}

type AuthorizationResultMsg struct {
	Login    string
	Password string
	Err      error
}

func (m *AuthorizationViewModel) authorize() tea.Cmd {
	values := m.inputSet.Values()
	return func() tea.Msg {
		_, err := m.authorizationService.Authorize(context.Background(), values["Login"], values["Password"])
		return AuthorizationResultMsg{
			Login:    values["Login"],
			Password: values["Password"],
			Err:      err,
		}
	}
}
