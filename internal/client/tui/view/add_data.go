package view

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mkolibaba/gophkeeper/internal/client"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/components/inputset"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/helper"
)

type AddDataResultMsg struct {
	Name string
	Err  error
}

type addViewKeyMap struct {
	ToggleFilepicker key.Binding
	SelectFile       key.Binding
	Send             key.Binding
	Exit             key.Binding
}

func (k addViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Send, k.Exit}
}

func (k addViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		// TODO: добавить помощь по навигации по инпутам
		{k.ToggleFilepicker},
		{k.SelectFile},
		{k.Send},
		{k.Exit},
	}
}

type AddDataViewModel struct {
	baseViewModel
	keyMap        addViewKeyMap
	dataType      helper.DataType
	inputSet      *inputset.Model
	send          func(map[string]string) error
	loginService  client.LoginService
	noteService   client.NoteService
	binaryService client.BinaryService
	cardService   client.CardService
}

func InitialAddDataViewModel(
	loginService client.LoginService,
	noteService client.NoteService,
	binaryService client.BinaryService,
	cardService client.CardService,
) *AddDataViewModel {
	keyMap := addViewKeyMap{
		ToggleFilepicker: key.NewBinding(
			key.WithKeys("ctrl+p"),
			key.WithHelp("ctrl+p", "toggle file picker"),
		),
		SelectFile: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select file"),
		),
		Send: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "save"),
		),
		Exit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "exit"),
		),
	}

	return &AddDataViewModel{
		keyMap:        keyMap,
		loginService:  loginService,
		noteService:   noteService,
		binaryService: binaryService,
		cardService:   cardService,
	}
}

func (m *AddDataViewModel) Init() tea.Cmd {
	return m.inputSet.Init()
}

func (m *AddDataViewModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case AddDataCallMsg:
		m.ResetFor(msg.t)
		return m.Init()

	case AddDataResultMsg:
		m.inputSet.Err = msg.Err
		m.inputSet.Reset()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Send):
			return m.save()

		case key.Matches(msg, m.keyMap.Exit):
			return ExitAddDataView
		}

	}

	cmd := m.inputSet.Update(msg)

	filepickerSelected := m.inputSet.Current().Placeholder() == "File path"
	m.keyMap.ToggleFilepicker.SetEnabled(filepickerSelected)
	m.keyMap.SelectFile.SetEnabled(filepickerSelected)

	return cmd
}

func (m *AddDataViewModel) View() string {
	// Строка помощи
	hm := help.New()
	hm.ShowAll = true
	helpView := lipgloss.NewStyle().PaddingLeft(1).Render(hm.View(m.keyMap))

	addDataView := helper.Borderize(
		fmt.Sprintf("Add %s", m.dataType),
		"",
		lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingTop(1).
			Render(m.inputSet.View()),
		m.Width,
		m.Height-lipgloss.Height(helpView),
	)

	return lipgloss.JoinVertical(lipgloss.Top, addDataView, helpView)
}

func (m *AddDataViewModel) ResetFor(t helper.DataType) {
	m.dataType = t

	switch m.dataType {
	case helper.DataTypeLogin:
		m.inputSet = inputset.NewInputSet(
			inputset.NewTextInput("Name"),
			inputset.NewTextInput("Login"),
			inputset.NewTextInput("Password", inputset.WithEchoModePassword()),
			inputset.NewTextInput("Website"),
			inputset.NewTextInput("Notes"),
		)
		m.send = func(values map[string]string) error {
			return m.loginService.Save(context.Background(), client.LoginData{
				Name:     values["Name"],
				Login:    values["Login"],
				Password: values["Password"],
				Website:  values["Website"],
				Notes:    values["Notes"],
			})
		}
	case helper.DataTypeNote:
		m.inputSet = inputset.NewInputSet(
			inputset.NewTextInput("Name"),
			inputset.NewTextArea("Text"),
		)
		m.send = func(values map[string]string) error {
			return m.noteService.Save(context.Background(), client.NoteData{
				Name: values["Name"],
				Text: values["Text"],
			})
		}
	case helper.DataTypeBinary:
		m.inputSet = inputset.NewInputSet(
			inputset.NewTextInput("Name"),
			inputset.NewFilePicker("File path"),
			inputset.NewTextInput("Notes"),
		)
		m.send = func(values map[string]string) error {
			return m.binaryService.Save(context.Background(), client.BinaryData{
				Name:     values["Name"],
				Filename: values["File path"],
				Notes:    values["Notes"],
			})
		}
	case helper.DataTypeCard:
		m.inputSet = inputset.NewInputSet(
			inputset.NewTextInput("Name"),
			inputset.NewTextInput("Number"),
			inputset.NewTextInput("Expiration date"),
			inputset.NewTextInput("CVV"),
			inputset.NewTextInput("Cardholder"),
			inputset.NewTextInput("Notes"),
		)
		m.send = func(values map[string]string) error {
			return m.cardService.Save(context.Background(), client.CardData{
				Name:       values["Name"],
				Number:     values["Number"],
				ExpDate:    values["Expiration date"],
				CVV:        values["CVV"],
				Cardholder: values["Cardholder"],
				Notes:      values["Notes"],
			})
		}
	}
}

func (m *AddDataViewModel) save() tea.Cmd {
	values := m.inputSet.Values()
	return func() tea.Msg {
		return AddDataResultMsg{
			Name: values["Name"],
			Err:  m.send(values),
		}
	}
}
