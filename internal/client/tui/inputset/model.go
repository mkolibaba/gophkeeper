package inputset

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/helper"
)

const (
	defaultWidth     = 50
	defaultCharLimit = defaultWidth
)

type Option func(*textinput.Model)

func WithFocus() Option {
	return func(input *textinput.Model) {
		input.Focus()
	}
}

func WithEchoModePassword() Option {
	return func(input *textinput.Model) {
		input.EchoMode = textinput.EchoPassword
		input.EchoCharacter = '•'
	}
}

func WithCharLimit(charLimit int) Option {
	return func(input *textinput.Model) {
		input.CharLimit = charLimit
	}
}

func NewInput(placeholder string, opts ...Option) textinput.Model {
	input := textinput.New()
	input.Placeholder = placeholder
	input.CharLimit = defaultCharLimit
	input.Width = defaultWidth
	input.Cursor.SetMode(cursor.CursorStatic)
	input.PromptStyle = helper.HeaderStyle

	for _, o := range opts {
		o(&input)
	}

	return input
}

type Model struct {
	inputs  []textinput.Model
	focused int
}

func NewInputSet(inputs ...textinput.Model) *Model {
	return &Model{
		inputs: inputs,
	}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "up", "down", "tab":
			var at int
			if keypress == "up" {
				at = (m.focused - 1 + len(m.inputs)) % len(m.inputs)
			} else {
				at = (m.focused + 1) % len(m.inputs)
			}
			m.setFocus(at)
		}
	}

	return m.updateInputs(msg)
}

func (m *Model) View() string {
	var lines []string
	for _, input := range m.inputs {
		lines = append(lines, input.View())
	}

	return lipgloss.JoinVertical(lipgloss.Top, lines...)
}

func (m *Model) Values() map[string]string {
	values := map[string]string{}
	for _, input := range m.inputs {
		values[input.Placeholder] = input.Value()
	}
	return values
}

func (m *Model) Reset(focusAt int) {
	for i := range m.inputs {
		m.inputs[i].Reset()
	}

	if focusAt >= 0 {
		m.setFocus(focusAt)
	}
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m *Model) setFocus(at int) {
	m.focused = at
	for i := range m.inputs {
		if at == i {
			m.inputs[i].Focus()
			continue
		}
		m.inputs[i].Blur()
	}
}
