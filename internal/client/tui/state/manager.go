package state

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mkolibaba/gophkeeper/internal/client"
	"go.uber.org/zap"
	"sync"
)

type Manager struct {
	loginService  client.LoginService
	noteService   client.NoteService
	binaryService client.BinaryService
	cardService   client.CardService
	logger        *zap.Logger
}

func NewManager(
	loginService client.LoginService,
	noteService client.NoteService,
	binaryService client.BinaryService,
	cardService client.CardService,
	logger *zap.Logger,
) *Manager {
	return &Manager{
		loginService:  loginService,
		noteService:   noteService,
		binaryService: binaryService,
		cardService:   cardService,
		logger:        logger,
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
