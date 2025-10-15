package sqlite

import (
	"context"
	"fmt"
	"github.com/mkolibaba/gophkeeper/server"
	sqlc "github.com/mkolibaba/gophkeeper/server/sqlite/sqlc/gen"
)

type CardService struct {
	qs *sqlc.Queries
}

func NewCardService(queries *sqlc.Queries) *CardService {
	return &CardService{
		qs: queries,
	}
}

func (s *CardService) Save(ctx context.Context, data server.CardData, user string) error {
	err := s.qs.SaveCard(ctx, sqlc.SaveCardParams{
		Name:       data.Name,
		Number:     data.Number,
		ExpDate:    data.ExpDate,
		Cvv:        data.CVV,
		Cardholder: data.Cardholder,
		Notes:      stringOrNull(data.Notes),
		User:       user,
	})

	return tryUnwrapSaveError(err)
}

func (s *CardService) GetAll(ctx context.Context, user string) ([]server.CardData, error) {
	cards, err := s.qs.GetAllCards(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("get all: %w", err)
	}

	var result []server.CardData
	for _, card := range cards {
		result = append(result, server.CardData{
			Name:       card.Name,
			Number:     card.Number,
			ExpDate:    card.ExpDate,
			CVV:        card.Cvv,
			Cardholder: card.Cardholder,
			Notes:      stringOrEmpty(card.Notes),
		})
	}

	return result, nil
}

func (s *CardService) Update(ctx context.Context, data server.CardData, user string) error {
	// TODO: implement
	panic("implement me")
}

func (s *CardService) Remove(ctx context.Context, name string, user string) error {
	n, err := s.qs.RemoveCard(ctx, name)
	if err != nil {
		return fmt.Errorf("remove: %w", err)
	}
	if n == 0 {
		return server.ErrDataNotFound
	}
	return nil
}
