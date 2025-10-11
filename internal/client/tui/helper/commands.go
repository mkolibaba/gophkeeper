package helper

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mkolibaba/gophkeeper/internal/client"
)

type LoadDataMsg []client.Data

func LoadData(data []client.Data) tea.Cmd {
	return func() tea.Msg {
		return LoadDataMsg(data)
	}
}
