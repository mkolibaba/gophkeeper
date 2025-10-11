package inputset

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Input interface {
	View() string
	Update(tea.Msg) (Input, tea.Cmd)
	Placeholder() string
	Value() string
	Focus() tea.Cmd
	Blur()
	Reset()
}

type TextInput struct {
	textinput.Model
}

func (i TextInput) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	i.Model, cmd = i.Model.Update(msg)
	return &i, cmd
}

func (i TextInput) Placeholder() string {
	return i.Model.Placeholder
}

type TextArea struct {
	textarea.Model
}

func (i TextArea) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	i.Model, cmd = i.Model.Update(msg)
	return &i, cmd
}

func (i TextArea) Placeholder() string {
	return i.Model.Placeholder
}
