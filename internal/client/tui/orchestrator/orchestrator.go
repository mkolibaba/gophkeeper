package orchestrator

import (
	"context"
	"fmt"
	"github.com/mkolibaba/gophkeeper/internal/client"
	"github.com/mkolibaba/gophkeeper/internal/client/tui/state"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"sync"
)

type Orchestrator struct {
	loginService  client.LoginService
	noteService   client.NoteService
	binaryService client.BinaryService
	cardService   client.CardService

	// TODO: временно
	manager *state.Manager

	logger *zap.Logger
}

type OrchestratorParams struct {
	fx.In

	LoginService  client.LoginService
	NoteService   client.NoteService
	BinaryService client.BinaryService
	CardService   client.CardService

	// TODO: временно
	StateManager *state.Manager

	Logger *zap.Logger
}

func New(p OrchestratorParams) *Orchestrator {
	return &Orchestrator{
		loginService:  p.LoginService,
		noteService:   p.NoteService,
		binaryService: p.BinaryService,
		cardService:   p.CardService,
		manager:       p.StateManager,
		logger:        p.Logger,
	}
}

func (o *Orchestrator) GetAll(ctx context.Context) []client.Data {
	var result []client.Data

	ch := make(chan client.Data)
	var wg sync.WaitGroup

	// TODO: как можно уйти от дублирования?

	wg.Go(func() {
		elems, err := o.loginService.GetAll(ctx)
		if err != nil {
			o.logger.Error(err.Error())
			return
		}

		for _, el := range elems {
			ch <- el
		}
	})
	wg.Go(func() {
		elems, err := o.noteService.GetAll(ctx)
		if err != nil {
			o.logger.Error(err.Error())
			return
		}

		for _, el := range elems {
			ch <- el
		}
	})
	wg.Go(func() {
		elems, err := o.binaryService.GetAll(ctx)
		if err != nil {
			o.logger.Error(err.Error())
			return
		}

		for _, el := range elems {
			ch <- el
		}
	})
	wg.Go(func() {
		elems, err := o.cardService.GetAll(ctx)
		if err != nil {
			o.logger.Error(err.Error())
			return
		}

		for _, el := range elems {
			ch <- el
		}
	})

	go func() {
		wg.Wait()
		close(ch)
	}()

	for el := range ch {
		result = append(result, el)
	}

	return result
}

func (o *Orchestrator) Remove(ctx context.Context, data client.Data) (string, error) {
	switch data := data.(type) {
	case client.LoginData:
		return data.Name, o.loginService.Remove(ctx, data.Name)
	case client.NoteData:
		return data.Name, o.noteService.Remove(ctx, data.Name)
	case client.BinaryData:
		return data.Name, o.binaryService.Remove(ctx, data.Name)
	case client.CardData:
		return data.Name, o.cardService.Remove(ctx, data.Name)
	}

	return "", fmt.Errorf("unknown data type %T", data)
}
