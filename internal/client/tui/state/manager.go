package state

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mkolibaba/gophkeeper/internal/client"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// TODO: нужен ли этот компонент?
type Manager struct {
	authorizationService client.AuthorizationService
	loginService         client.LoginService
	noteService          client.NoteService
	binaryService        client.BinaryService
	cardService          client.CardService
	session              *client.Session
	logger               *zap.Logger
}

type ManagerParams struct {
	fx.In

	AuthorizationService client.AuthorizationService
	LoginService         client.LoginService
	NoteService          client.NoteService
	BinaryService        client.BinaryService
	CardService          client.CardService
	Session              *client.Session
	Logger               *zap.Logger
}

func NewManager(p ManagerParams) *Manager {
	return &Manager{
		authorizationService: p.AuthorizationService,
		loginService:         p.LoginService,
		noteService:          p.NoteService,
		binaryService:        p.BinaryService,
		cardService:          p.CardService,
		session:              p.Session,
		logger:               p.Logger,
	}
}

type AuthorizationResultMsg struct {
	Login    string
	Password string
	Err      error
}

func (m *Manager) Authorize(login, password string) tea.Cmd {
	return func() tea.Msg {
		_, err := m.authorizationService.Authorize(context.Background(), login, password)
		return AuthorizationResultMsg{
			Login:    login,
			Password: password,
			Err:      err,
		}
	}
}
