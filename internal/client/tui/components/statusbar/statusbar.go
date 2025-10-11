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

type NotificationMsg struct {
	Text string
	T    NotificationType
}

type Model struct {
	Width            int
	notificationText string
	notificationType NotificationType
	CurrentUser      string
	ttl              time.Duration
}

func New() *Model {
	return &Model{
		ttl: 3 * time.Second,
	}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case NotificationMsg:
		m.notificationText = msg.Text
		m.notificationType = msg.T
		return m.ResetNotification()
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

	// TODO: вынести в компонент
	var restColor lipgloss.Color

	switch m.notificationType {
	case NotificationOk:
		restColor = lipgloss.Color("115")
	case NotificationError:
		restColor = lipgloss.Color("169")
	case NotificationNone:
		restColor = lipgloss.Color("105")
	}

	rest := lipgloss.NewStyle().
		Width(m.Width - w(helpInfo) - w(user)).
		PaddingLeft(1).
		Background(restColor).
		Render(m.notificationText)

	return lipgloss.JoinHorizontal(lipgloss.Top, helpInfo, rest, user)
}

// TODO: если быстро 2 уведомления прокинуть, то обнуление от первого сработает быстрее, и второе покажется малое время
func (m *Model) ResetNotification() tea.Cmd {
	return tea.Tick(m.ttl, func(t time.Time) tea.Msg {
		return NotificationMsg{
			Text: "",
			T:    NotificationNone,
		}
	})
}

func NotifyOk(text string) tea.Cmd {
	return notify(text, NotificationOk)
}

func NotifyError(text string) tea.Cmd {
	return notify(text, NotificationError)
}

func notify(text string, t NotificationType) tea.Cmd {
	return func() tea.Msg {
		return NotificationMsg{
			Text: text,
			T:    t,
		}
	}
}
