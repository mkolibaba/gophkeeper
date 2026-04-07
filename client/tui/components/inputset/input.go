package inputset

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/filepicker"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mkolibaba/gophkeeper/client/tui/helper"
)

const (
	defaultWidth     = 50
	defaultCharLimit = defaultWidth
)

type Input interface {
	Init() tea.Cmd
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

func NewTextInput(placeholder string, opts ...Option) Input {
	input := textinput.New()
	input.Placeholder = placeholder
	input.CharLimit = defaultCharLimit
	input.SetWidth(defaultWidth)

	s := input.Styles()
	s.Cursor.Blink = false
	s.Focused.Prompt = promptStyle
	s.Blurred.Prompt = promptStyle
	input.SetStyles(s)

	for _, o := range opts {
		o(&input)
	}

	return &TextInput{input}
}

func (t TextInput) Init() tea.Cmd {
	return nil
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

func NewTextArea(placeholder string, opts ...TextAreaOption) Input {
	area := textarea.New()
	area.Placeholder = placeholder
	area.CharLimit = 2000
	area.SetWidth(defaultWidth)
	area.SetHeight(15)
	area.SetPromptFunc(2, func(info textarea.PromptInfo) string {
		if info.LineNumber == 0 {
			return promptStyle.Render("> ")
		}
		return "  "
	})

	s := area.Styles()
	s.Cursor.Blink = false
	area.SetStyles(s)
	area.ShowLineNumbers = false

	for _, o := range opts {
		o(&area)
	}

	return &TextArea{area}
}

func (t TextArea) Init() tea.Cmd {
	return nil
}

func (i TextArea) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	i.Model, cmd = i.Model.Update(msg)
	return &i, cmd
}

func (i TextArea) Placeholder() string {
	return i.Model.Placeholder
}

type FilePicker struct {
	filePicker   filepicker.Model
	textInput    textinput.Model
	focused      bool
	pickingMode  bool
	selectedFile string
	disabled     bool
}

func NewFilePicker(placeholder string, opts ...FilePickerOption) Input {
	input := textinput.New()
	input.Placeholder = placeholder
	input.CharLimit = 255
	input.SetWidth(100)

	s := input.Styles()
	s.Cursor.Blink = false
	s.Focused.Prompt = promptStyle
	s.Blurred.Prompt = promptStyle
	input.SetStyles(s)

	picker := filepicker.New()
	picker.CurrentDirectory, _ = os.Getwd()
	picker.SetHeight(15)
	picker.Styles.Selected = lipgloss.NewStyle().Background(helper.HeaderColor)
	picker.Cursor = " "

	fp := &FilePicker{
		filePicker: picker,
		textInput:  input,
	}

	for _, o := range opts {
		o(fp)
	}

	return fp
}

func (f *FilePicker) Init() tea.Cmd {
	return f.filePicker.Init()
}

func (i *FilePicker) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd

	// TODO(minor): неплохо было бы пропускать этот инпут при навигации
	if i.disabled {
		return i, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+p":
			i.pickingMode = !i.pickingMode
			if i.pickingMode {
				i.textInput.Blur()
			} else {
				i.textInput.Focus()
			}
		}

	default:
		i.filePicker, cmd = i.filePicker.Update(msg) // TODO: причесать это. тут скорее всего будет filePicker.Init()
		return i, cmd
	}

	if i.pickingMode {
		i.filePicker, cmd = i.filePicker.Update(msg)
	} else {
		i.textInput, cmd = i.textInput.Update(msg)
	}

	if didSelect, path := i.filePicker.DidSelectFile(msg); didSelect {
		i.textInput.SetValue(path)
	}

	return i, cmd
}

func (i FilePicker) View() string {
	view := i.textInput.View()

	if i.pickingMode {
		view += fmt.Sprintf("\n  %s %s\n",
			lipgloss.NewStyle().Foreground(helper.HeaderColor).Render("Directory: "),
			i.filePicker.CurrentDirectory) + i.filePicker.View()
	}

	return view
}

func (i FilePicker) Placeholder() string {
	return i.textInput.Placeholder
}

func (i FilePicker) Value() string {
	return i.textInput.Value()
}

func (i *FilePicker) Focus() tea.Cmd {
	i.focused = true
	i.textInput.Focus()
	return nil
}

func (i *FilePicker) Blur() {
	i.focused = false
	i.pickingMode = false
	i.textInput.Blur()
}

func (i *FilePicker) Reset() {
	i.filePicker.CurrentDirectory, _ = os.Getwd()
	i.textInput.Reset()
}
