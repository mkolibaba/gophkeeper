package statusbar

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time"
)

type NotificationType uint8

const (
	NotificationNone NotificationType = iota
	NotificationOk
	NotificationError
)

var notificationColors = map[NotificationType]lipgloss.Color{
	NotificationNone:  lipgloss.Color("105"),
	NotificationOk:    lipgloss.Color("34"),
	NotificationError: lipgloss.Color("169"),
}

type clearNotificationMsg struct{}

type Model struct {
	Width            int
	notificationText string
	notificationType NotificationType
	CurrentUser      string
	ttl              time.Duration
	until            time.Time
}

func New() *Model {
	return &Model{
		ttl: 3 * time.Second,
	}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg.(type) {
	case clearNotificationMsg:
		if time.Since(m.until) > 0 {
			m.notificationText = ""
			m.notificationType = NotificationNone
		}
	}

	return nil
}

func (m *Model) View() string {
	// TODO(trivial): попробовать другие стили

	w := lipgloss.Width

	helpInfo := lipgloss.NewStyle().
		Width(8).
		PaddingLeft(1).
		Background(lipgloss.Color("243")).
		Render("h Help")

	user := lipgloss.NewStyle().
		Width(w(m.CurrentUser) + 2).
		PaddingLeft(1).
		Background(lipgloss.Color("243")).
		Render(m.CurrentUser)

	rest := lipgloss.NewStyle().
		Width(m.Width - w(helpInfo) - w(user)).
		PaddingLeft(1).
		Background(notificationColors[m.notificationType]).
		Render(m.notificationText)

	return lipgloss.JoinHorizontal(lipgloss.Top, helpInfo, rest, user)
}

func (m *Model) NotifyOk(text string) tea.Cmd {
	return m.notify(text, NotificationOk)
}

func (m *Model) NotifyError(text string) tea.Cmd {
	return m.notify(text, NotificationError)
}

func (m *Model) notify(text string, t NotificationType) tea.Cmd {
	m.notificationText = text
	m.notificationType = t
	m.until = time.Now().Add(m.ttl)
	return tea.Tick(m.ttl, func(t time.Time) tea.Msg {
		return clearNotificationMsg{}
	})
}
