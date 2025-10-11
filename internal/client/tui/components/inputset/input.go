package inputset

import (
	"fmt"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

type Input interface {
	Init() tea.Cmd
	View() string
	Update(tea.Msg) (Input, tea.Cmd)
	Placeholder() string // TODO: можно переименовать на name
	Value() string
	Focus() tea.Cmd
	Blur() // TODO: ожно убрать и сделать focus(bool)
	Reset()
}

type TextInput struct {
	textinput.Model
}

func NewTextInput(placeholder string, opts ...Option) Input {
	input := textinput.New()
	input.Placeholder = placeholder
	input.CharLimit = defaultCharLimit
	input.Width = defaultWidth
	input.Cursor.SetMode(cursor.CursorStatic)
	input.PromptStyle = promptStyle

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

func NewTextArea(placeholder string) Input {
	area := textarea.New()
	area.Placeholder = placeholder
	area.CharLimit = 2000 // TODO: в константы
	area.SetWidth(defaultWidth)
	area.SetHeight(15)
	area.SetPromptFunc(2, func(lineIdx int) string {
		if lineIdx == 0 {
			return promptStyle.Render("> ")
		}
		return "  "
	})
	area.Cursor.SetMode(cursor.CursorStatic)
	area.ShowLineNumbers = false

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
}

func NewFilePicker(placeholder string) Input {
	input := textinput.New()
	input.Placeholder = placeholder
	input.CharLimit = 255
	input.Width = 100
	input.Cursor.SetMode(cursor.CursorStatic)
	input.PromptStyle = promptStyle

	picker := filepicker.New()
	picker.CurrentDirectory, _ = os.Getwd()
	picker.SetHeight(15)
	picker.Styles.Selected = lipgloss.NewStyle().Background(lipgloss.Color("171"))
	picker.Cursor = " "

	return &FilePicker{
		filePicker: picker,
		textInput:  input,
	}
}

func (f *FilePicker) Init() tea.Cmd {
	return f.filePicker.Init()
}

func (i *FilePicker) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd

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

	// Did the user select a file?
	if didSelect, path := i.filePicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		i.textInput.SetValue(path)
	}

	return i, cmd
}

func (i FilePicker) View() string {
	view := i.textInput.View()

	if i.pickingMode {
		view += fmt.Sprintf("\n  %s %s\n", lipgloss.NewStyle().
			Foreground(lipgloss.Color("171")).Render("Directory:"), i.filePicker.CurrentDirectory) + i.filePicker.View()
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
