package view

import tea "github.com/charmbracelet/bubbletea"

type ExitAddDataViewMsg struct{}

func ExitAddDataView() tea.Msg {
	return ExitAddDataViewMsg{}
}
