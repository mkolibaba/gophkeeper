package view

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/helper"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/state"
)

var (
	authErrorRenderer = lipgloss.NewStyle().
				Foreground(lipgloss.Color("169")).
				Render

	authViewRenderer = lipgloss.NewStyle().
				PaddingTop(1).
				Render
)

type AuthorizationViewModel struct {
	baseViewModel
	inputs  []textinput.Model
	focused int
	err     error
	manager *state.Manager
}

func InitialAuthorizationViewModel(manager *state.Manager) *AuthorizationViewModel {
	loginInput := textinput.New()
	loginInput.Placeholder = "Login"
	loginInput.Focus()
	loginInput.CharLimit = 20
	loginInput.Width = 20
	loginInput.Cursor.SetMode(cursor.CursorStatic)
	loginInput.PromptStyle = helper.HeaderStyle

	passwordInput := textinput.New()
	passwordInput.Placeholder = "Password"
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '•'
	passwordInput.CharLimit = 20
	passwordInput.Width = 20
	passwordInput.Cursor.SetMode(cursor.CursorStatic)
	passwordInput.PromptStyle = helper.HeaderStyle

	return &AuthorizationViewModel{
		inputs:  []textinput.Model{loginInput, passwordInput},
		manager: manager,
	}
}

func (m *AuthorizationViewModel) Init() tea.Cmd {
	return nil
}

func (m *AuthorizationViewModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case state.AuthorizationResultMsg:
		m.err = msg.Err
		m.inputs[0].SetValue("")
		m.inputs[1].SetValue("")
		m.setFocus(0)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "esc":
			return tea.Quit

		case "up", "down", "tab":
			var at int
			if keypress == "up" {
				at = (m.focused - 1 + len(m.inputs)) % len(m.inputs)
			} else {
				at = (m.focused + 1) % len(m.inputs)
			}
			m.setFocus(at)

		case "enter":
			login, password := m.inputs[0].Value(), m.inputs[1].Value()
			return m.manager.Authorize(login, password)
		}
	}

	return m.updateInputs(msg)
}

func (m *AuthorizationViewModel) View() string {
	w := m.Width - helper.ContentStyle.GetHorizontalFrameSize()

	borderTop := helper.RenderBorderTop(helper.ContentStyle, "Authorization", w)

	h := m.Height - lipgloss.Height(borderTop) - helper.ContentStyle.GetBorderBottomSize()

	oldV := func() string {
		lines := []string{m.inputs[0].View(), m.inputs[1].View()}
		if m.err != nil {
			lines = append(lines, "", authErrorRenderer(m.err.Error()))
		}
		return authViewRenderer(lipgloss.JoinVertical(lipgloss.Top, lines...))
	}

	v := oldV()

	authorizationView := helper.ContentStyle.
		BorderTop(false).
		Width(w).
		Height(h / 2).
		PaddingLeft(1).
		Render(v)

	return lipgloss.JoinVertical(lipgloss.Top, borderTop, authorizationView)
}

func (m *AuthorizationViewModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m *AuthorizationViewModel) setFocus(at int) {
	m.focused = at
	for i := range m.inputs {
		if at == i {
			m.inputs[i].Focus()
			continue
		}
		m.inputs[i].Blur()
	}
}
