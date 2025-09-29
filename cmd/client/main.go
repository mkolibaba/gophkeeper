package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/mkolibaba/gophkeeper/internal/client/tui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	loginTab := tui.NewTab("Login", []list.Item{
		tui.ListItem{Name: "Google", Desc: "iivanov"},
		tui.ListItem{Name: "Ozon", Desc: "+79031002030"},
		tui.ListItem{Name: "Wildberries", Desc: "+79031002030"},
		tui.ListItem{Name: "Госуслуги", Desc: "iivanov@gmail.com"},
		tui.ListItem{Name: "Mail.ru", Desc: "ivanivanov"},
		tui.ListItem{Name: "VK", Desc: "ivanivanov@mail.ru"},
	})
	noteTab := tui.NewTab("Note", []list.Item{
		tui.NewNoteItem("Записки о природе", "Кто никогда не видал, как растет клюква, тот может очень долго идти по болоту и не замечать, что он по клюкве идет."),
		tui.NewNoteItem("Мысль", "Живешь ты, может быть, сам триста лет, и кто породил тебя, тот в яичке своем пересказал все, что он тоже узнал за свои триста лет жизни."),
	})
	binaryTab := tui.NewTab("Binary", []list.Item{})
	cardTab := tui.NewTab("Card", []list.Item{
		tui.NewCardItem("Сбербанк", "2200123456789019"),
		tui.NewCardItem("Т-Банк", "2201987654321000"),
	})
	settingsTab := tui.NewTab("Settings", []list.Item{})

	m := tui.NewBubble(loginTab, noteTab, binaryTab, cardTab, settingsTab)

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
