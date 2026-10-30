package editdata

import (
	"context"
	"fmt"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mkolibaba/gophkeeper/client/core"
	"github.com/mkolibaba/gophkeeper/client/tui/components/inputset"
	"github.com/mkolibaba/gophkeeper/client/tui/helper"
	"github.com/mkolibaba/gophkeeper/client/tui/view"
	"go.uber.org/fx"
)

type keyMap struct {
	Send key.Binding
	Exit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Send, k.Exit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Send},
		{k.Exit},
	}
}

type Model struct {
	view.BaseModel
	keyMap        keyMap
	inputSet      *inputset.Model
	dataName      string
	send          func(map[string]string) error
	loginService  core.LoginService
	noteService   core.NoteService
	binaryService core.BinaryService
	cardService   core.CardService
}

type Params struct {
	fx.In

	LoginService  core.LoginService
	NoteService   core.NoteService
	BinaryService core.BinaryService
	CardService   core.CardService
}

func New(p Params) *Model {
	return &Model{
		keyMap: keyMap{
			Send: key.NewBinding(
				key.WithKeys("ctrl+s"),
				key.WithHelp("ctrl+s", "save"),
			),
			Exit: key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("esc", "exit"),
			),
		},
		loginService:  p.LoginService,
		noteService:   p.NoteService,
		binaryService: p.BinaryService,
		cardService:   p.CardService,
	}
}

func (m *Model) Init() tea.Cmd {
	return m.inputSet.Init()
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case EditDataResultMsg:
		m.inputSet.Err = msg.Err

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Send):
			return m.save()

		case key.Matches(msg, m.keyMap.Exit):
			return Exit
		}
	}

	return m.inputSet.Update(msg)
}

func (m *Model) View() string {
	// Строка помощи
	hm := help.New()
	hm.ShowAll = true
	helpView := lipgloss.NewStyle().PaddingLeft(1).Render(hm.View(m.keyMap))

	editDataView := helper.Borderize(
		fmt.Sprintf("Edit %s", m.dataName),
		"",
		lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingTop(1).
			Render(m.inputSet.View()),
		m.Width,
		m.Height-lipgloss.Height(helpView),
	)

	return lipgloss.JoinVertical(lipgloss.Top, editDataView, helpView)
}

func (m *Model) ResetFor(data core.Data) {
	m.dataName = data.GetName()
	switch data := data.(type) {
	case core.LoginData:
		m.inputSet = inputset.NewInputSet(
			inputset.NewTextInput("Name", inputset.WithValue(data.Name)),
			inputset.NewTextInput("Login", inputset.WithValue(data.Login)),
			inputset.NewTextInput("Password", inputset.WithValue(data.Password), inputset.WithEchoModePassword()),
			inputset.NewTextInput("Website", inputset.WithValue(data.Website)),
			inputset.NewTextInput("Notes", inputset.WithValue(data.Notes)),
		)
		m.send = func(values map[string]string) error {
			name := values["Name"]
			login := values["Login"]
			password := values["Password"]
			website := values["Website"]
			notes := values["Notes"]
			return m.loginService.Update(context.Background(), core.LoginDataUpdate{
				ID:       data.ID,
				Name:     &name,
				Login:    &login,
				Password: &password,
				Website:  &website,
				Notes:    &notes,
			})
		}
	case core.NoteData:
		m.inputSet = inputset.NewInputSet(
			inputset.NewTextInput("Name", inputset.WithValue(data.Name)),
			inputset.NewTextArea("Text", inputset.WithTextAreaValue(data.Text)),
		)
		m.send = func(values map[string]string) error {
			name, text := values["Name"], values["Text"]
			return m.noteService.Update(context.Background(), core.NoteDataUpdate{
				ID:   data.ID,
				Name: &name,
				Text: &text,
			})
		}
	case core.BinaryData:
		m.inputSet = inputset.NewInputSet(
			inputset.NewTextInput("Name", inputset.WithValue(data.Name)),
			inputset.NewFilePicker("File path", inputset.WithFilePickerDisabled()),
			inputset.NewTextInput("Notes", inputset.WithValue(data.Notes)),
		)
		m.send = func(values map[string]string) error {
			name, notes := values["Name"], values["Notes"]
			return m.binaryService.Update(context.Background(), core.BinaryDataUpdate{
				ID:    data.ID,
				Name:  &name,
				Notes: &notes,
			})
		}
	case core.CardData:
		m.inputSet = inputset.NewInputSet(
			inputset.NewTextInput("Name", inputset.WithValue(data.Name)),
			inputset.NewTextInput("Number", inputset.WithValue(data.Number)),
			inputset.NewTextInput("Expiration date", inputset.WithValue(data.ExpDate)),
			inputset.NewTextInput("CVV", inputset.WithValue(data.CVV)),
			inputset.NewTextInput("Cardholder", inputset.WithValue(data.Cardholder)),
			inputset.NewTextInput("Notes", inputset.WithValue(data.Notes)),
		)
		m.send = func(values map[string]string) error {
			name := values["Name"]
			number := values["Number"]
			expDate := values["Expiration date"]
			cvv := values["CVV"]
			cardholder := values["Cardholder"]
			notes := values["Notes"]
			return m.cardService.Update(context.Background(), core.CardDataUpdate{
				ID:         data.ID,
				Name:       &name,
				Number:     &number,
				ExpDate:    &expDate,
				CVV:        &cvv,
				Cardholder: &cardholder,
				Notes:      &notes,
			})
		}
	}
}

type ExitMsg struct{}

func Exit() tea.Msg {
	return ExitMsg{}
}

type EditDataResultMsg struct {
	Name string
	Err  error
}

func (m *Model) save() tea.Cmd {
	values := m.inputSet.Values()
	return func() tea.Msg {
		return EditDataResultMsg{
			Name: values["Name"],
			Err:  m.send(values),
		}
	}
}
