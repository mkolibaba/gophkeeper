package state

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mkolibaba/gophkeeper/internal/client"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"sync"
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

type FetchDataMsg struct {
	Logins   []client.LoginData
	Notes    []client.NoteData
	Binaries []client.BinaryData
	Cards    []client.CardData
}

func (m *Manager) FetchData() tea.Cmd {
	return func() tea.Msg {
		var msg FetchDataMsg

		ctx := context.Background()

		var wg sync.WaitGroup
		wg.Go(func() {
			logins, err := m.loginService.GetAll(ctx)
			if err != nil {
				m.logger.Error(err.Error())
			}
			msg.Logins = logins
		})
		wg.Go(func() {
			notes, err := m.noteService.GetAll(ctx)
			if err != nil {
				m.logger.Error(err.Error())
			}
			msg.Notes = notes
		})
		wg.Go(func() {
			binaries, err := m.binaryService.GetAll(ctx)
			if err != nil {
				m.logger.Error(err.Error())
			}
			msg.Binaries = binaries
		})
		wg.Go(func() {
			cards, err := m.cardService.GetAll(ctx)
			if err != nil {
				m.logger.Error(err.Error())
			}
			msg.Cards = cards
		})

		wg.Wait()
		return msg
	}
}

func (m *Manager) StartDownloadBinary(data client.BinaryData) {
	go func() {
		m.binaryService.Download(context.Background(), data.Name)
	}()
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

func (m *Manager) SetInSession(login, password string) {
	m.session.SetCurrentUser(client.User{
		Login:    login,
		Password: password,
	})
}

func (m *Manager) GetCurrentUserLogin() string {
	user := m.session.GetCurrentUser()
	if user == nil {
		return ""
	}
	return user.Login
}
