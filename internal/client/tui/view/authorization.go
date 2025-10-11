package view

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/components/inputset"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/helper"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/state"
)

type AuthorizationViewModel struct {
	baseViewModel
	manager  *state.Manager
	inputSet *inputset.Model
}

func InitialAuthorizationViewModel(manager *state.Manager) *AuthorizationViewModel {
	return &AuthorizationViewModel{
		manager: manager,
		inputSet: inputset.NewInputSet(
			inputset.NewTextInput("Login"),
			inputset.NewTextInput("Password", inputset.WithEchoModePassword()),
		),
	}
}

func (m *AuthorizationViewModel) Init() tea.Cmd {
	return nil
}

func (m *AuthorizationViewModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case state.AuthorizationResultMsg:
		m.inputSet.Err = msg.Err
		m.inputSet.Reset(0)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "esc":
			return tea.Quit

		case "up", "down", "tab":
			return m.inputSet.Update(msg)

		case "enter":
			values := m.inputSet.Values()
			return m.manager.Authorize(values["Login"], values["Password"])
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
