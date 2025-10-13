package inputset

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/helper"
)

var (
	promptStyle = helper.HeaderStyle

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("169"))
)

type Option func(*textinput.Model)

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

type Model struct {
	Err error

	inputs  []Input
	focused int
}

// TODO: указать, что по умолчанию фокус устанавливается на нулевом инпуте
func NewInputSet(inputs ...Input) *Model {
	m := &Model{
		inputs: inputs,
	}
	m.setFocus(0)

	return m
}

func (m *Model) Init() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		cmds[i] = m.inputs[i].Init()
	}
	return tea.Batch(cmds...)
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "shift+tab", "tab":
			var at int
			if keypress == "shift+tab" {
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

	if m.Err != nil {
		lines = append(lines, "", errorStyle.Render(m.Err.Error()))
	}

	return lipgloss.JoinVertical(lipgloss.Top, lines...)
}

func (m *Model) Values() map[string]string {
	values := map[string]string{}
	for _, input := range m.inputs {
		values[input.Placeholder()] = input.Value()
	}
	return values
}

func (m *Model) Reset() {
	for i := range m.inputs {
		m.inputs[i].Reset()
	}
	m.setFocus(0)
}

func (m *Model) Current() Input {
	return m.inputs[m.focused]
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
