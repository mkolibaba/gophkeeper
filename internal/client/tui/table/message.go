package table

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mkolibaba/gophkeeper/internal/client"
	"go.uber.org/zap"
	"sync"
)

type FetchDataMsg struct {
	Logins   []client.LoginData
	Notes    []client.NoteData
	Binaries []client.BinaryData
	Cards    []client.CardData
}

func FetchData(
	loginService client.LoginService,
	noteService client.NoteService,
	binaryService client.BinaryService,
	cardService client.CardService,
	logger *zap.Logger,
) tea.Cmd {
	return func() tea.Msg {
		var msg FetchDataMsg

		ctx := context.Background()

		var wg sync.WaitGroup
		wg.Go(func() {
			logins, err := loginService.GetAll(ctx, "demo")
			if err != nil {
				logger.Error(err.Error())
			}
			msg.Logins = logins
		})
		wg.Go(func() {
			notes, err := noteService.GetAll(ctx, "demo")
			if err != nil {
				logger.Error(err.Error())
			}
			msg.Notes = notes
		})
		wg.Go(func() {
			binaries, err := binaryService.GetAll(ctx, "demo")
			if err != nil {
				logger.Error(err.Error())
			}
			msg.Binaries = binaries
		})
		wg.Go(func() {
			cards, err := cardService.GetAll(ctx, "demo")
			if err != nil {
				logger.Error(err.Error())
			}
			msg.Cards = cards
		})

		wg.Wait()
		return msg
	}
}
