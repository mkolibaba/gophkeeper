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
	filepicker.Model
	placeholder  string
	focused      bool
	pickingMode  bool
	selectedFile string
}

func NewFilePicker(placeholder string) Input {
	picker := filepicker.New()
	picker.CurrentDirectory, _ = os.Getwd()
	picker.SetHeight(15)

	return &FilePicker{
		Model:       picker,
		placeholder: placeholder,
	}
}

func (i FilePicker) Update(msg tea.Msg) (Input, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !i.focused {
			return &i, nil
		}

		switch msg.String() {
		case "ctrl+p":
			i.pickingMode = !i.pickingMode
		}
	}

	var cmd tea.Cmd
	i.Model, cmd = i.Model.Update(msg)

	// Did the user select a file?
	if didSelect, path := i.Model.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		i.selectedFile = path
	}

	return &i, cmd
}

func (i FilePicker) View() string {
	if i.pickingMode {
		return fmt.Sprintf("  %s\n", i.Model.CurrentDirectory) + i.Model.View()
	}

	content := i.selectedFile
	if content == "" {
		placeholderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		content = placeholderStyle.Render("File path (Ctrl+P to select)")
	}

	return promptStyle.Render("> ") + content
}

func (i FilePicker) Placeholder() string {
	return i.placeholder
}

func (i FilePicker) Value() string {
	return i.selectedFile
}

func (i *FilePicker) Focus() tea.Cmd {
	i.focused = true
	return nil
}

func (i *FilePicker) Blur() {
	i.focused = false
	i.pickingMode = false
}

func (i *FilePicker) Reset() {
	i.Model.CurrentDirectory, _ = os.Getwd()
}
