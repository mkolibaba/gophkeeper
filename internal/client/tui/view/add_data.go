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

// TODO: вынести в одно место
type DataType string

const (
	DataTypeLogin  = DataType("Login")
	DataTypeNote   = DataType("Note")
	DataTypeBinary = DataType("Binary")
	DataTypeCard   = DataType("Card")
)

type addViewKeyMap struct {
	Send key.Binding
	Exit key.Binding
}

func (k addViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Send, k.Exit}
}

func (k addViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		// TODO: добавить помощь по навигации по инпутам
		{k.Send},
		{k.Exit},
	}
}

type AddDataViewModel struct {
	baseViewModel
	keyMap        addViewKeyMap
	dataType      DataType
	inputSet      *inputset.Model
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
	return nil
}

func (m *AddDataViewModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case AddDataCallMsg:
		m.ResetFor(msg.t)

	case addDataErrMsg:
		m.inputSet.Err = msg.err
		m.inputSet.Reset(0)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Send):
			return m.send()

		case key.Matches(msg, m.keyMap.Exit):
			return ExitAddDataView
		}

	}

	return m.inputSet.Update(msg)
}

func (m *AddDataViewModel) View() string {
	// Строка помощи
	hm := help.New()
	hm.ShowAll = true
	helpView := lipgloss.NewStyle().PaddingLeft(1).Render(hm.View(m.keyMap))

	w := m.Width - helper.ContentStyle.GetHorizontalFrameSize()

	borderTop := helper.RenderBorderTop(helper.ContentStyle, fmt.Sprintf("Add %s", m.dataType), w)

	h := m.Height - lipgloss.Height(borderTop) - helper.ContentStyle.GetBorderBottomSize() - lipgloss.Height(helpView)

	content := m.inputSet.View()

	addDataView := helper.ContentStyle.
		BorderTop(false).
		Width(w).
		Height(h).
		PaddingLeft(1).
		PaddingTop(1).
		Render(content)

	return lipgloss.JoinVertical(lipgloss.Top, borderTop, addDataView, helpView)
}

func (m *AddDataViewModel) ResetFor(t DataType) {
	// TODO: добавить metadata
	m.dataType = t

	switch m.dataType {
	case DataTypeLogin:
		m.inputSet = inputset.NewInputSet(
			inputset.NewInput("Name", inputset.WithFocus(), inputset.WithPromptStyle(helper.HeaderStyle)),
			inputset.NewInput("Login", inputset.WithPromptStyle(helper.HeaderStyle)),
			inputset.NewInput("Password", inputset.WithEchoModePassword(), inputset.WithPromptStyle(helper.HeaderStyle)),
		)
	case DataTypeNote:
		m.inputSet = inputset.NewInputSet(
			inputset.NewInput("Name", inputset.WithFocus(), inputset.WithPromptStyle(helper.HeaderStyle)),
			inputset.NewInput("Text", inputset.WithCharLimit(2000), inputset.WithPromptStyle(helper.HeaderStyle)), // TODO: тут нужен textarea
		)
	case DataTypeBinary:
		m.inputSet = inputset.NewInputSet(
			inputset.NewInput("Name", inputset.WithFocus(), inputset.WithPromptStyle(helper.HeaderStyle)),
			inputset.NewInput("File path", inputset.WithPromptStyle(helper.HeaderStyle)),
		)
	case DataTypeCard:
		m.inputSet = inputset.NewInputSet(
			inputset.NewInput("Name", inputset.WithFocus(), inputset.WithPromptStyle(helper.HeaderStyle)),
			inputset.NewInput("Number", inputset.WithPromptStyle(helper.HeaderStyle)),
			inputset.NewInput("Expiration date", inputset.WithPromptStyle(helper.HeaderStyle)),
			inputset.NewInput("CVV", inputset.WithPromptStyle(helper.HeaderStyle)),
			inputset.NewInput("Cardholder", inputset.WithPromptStyle(helper.HeaderStyle)),
		)
	}
}

func (m *AddDataViewModel) send() tea.Cmd {
	// TODO: после вебинара понять что должно быть в metadata
	values := m.inputSet.Values()

	switch m.dataType {
	case DataTypeLogin:
		data := client.LoginData{
			Name:     values["Name"],
			Login:    values["Login"],
			Password: values["Password"],
			Metadata: nil, // TODO: не забыть
		}
		return func() tea.Msg {
			err := m.loginService.Save(context.Background(), data)
			if err != nil {
				return addDataErrMsg{err: err}
			}

			return tea.Sequence(
				ExitAddDataView,
				func() tea.Msg {
					return notificationMsg{
						text: fmt.Sprintf("Add %s successfully", data.Name),
						t:    notificationOk,
					}
				},
			)()
		}
	case DataTypeNote:
		data := client.NoteData{
			Name:     values["Name"],
			Text:     values["Text"],
			Metadata: nil, // TODO: не забыть
		}
		return func() tea.Msg {
			err := m.noteService.Save(context.Background(), data)
			if err != nil {
				return addDataErrMsg{err: err}
			}

			return tea.Sequence(
				ExitAddDataView,
				func() tea.Msg {
					return notificationMsg{
						text: fmt.Sprintf("Add %s successfully", data.Name),
						t:    notificationOk,
					}
				},
			)()
		}
	case DataTypeBinary:
		data := client.BinaryData{
			Name:     values["Name"],
			FileName: values["File path"],
			Metadata: nil,
		}
		return func() tea.Msg {
			err := m.binaryService.Save(context.Background(), data)
			if err != nil {
				return addDataErrMsg{err: err}
			}

			return tea.Sequence(
				ExitAddDataView,
				func() tea.Msg {
					return notificationMsg{
						text: fmt.Sprintf("Add %s successfully", data.Name),
						t:    notificationOk,
					}
				},
			)()
		}
	case DataTypeCard:
		data := client.CardData{
			Name:       values["Name"],
			Number:     values["Number"],
			ExpDate:    values["Expiration date"],
			CVV:        values["CVV"],
			Cardholder: values["Cardholder"],
			Metadata:   nil,
		}
		return func() tea.Msg {
			err := m.cardService.Save(context.Background(), data)
			if err != nil {
				return addDataErrMsg{err: err}
			}

			return tea.Sequence(
				ExitAddDataView,
				func() tea.Msg {
					return notificationMsg{
						text: fmt.Sprintf("Add %s successfully", data.Name),
						t:    notificationOk,
					}
				},
			)()
		}
	}

	return nil
}

type addDataErrMsg struct {
	err error
}
