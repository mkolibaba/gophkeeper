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

func (s *CardService) Create(ctx context.Context, data server.CardData) error {
	err := s.qs.InsertCard(ctx, sqlc.InsertCardParams{
		Name:       data.Name,
		Number:     data.Number,
		ExpDate:    data.ExpDate,
		Cvv:        data.CVV,
		Cardholder: data.Cardholder,
		Notes:      stringOrNull(data.Notes),
		User:       server.UserFromContext(ctx),
	})
	return unwrapInsertError(err)
}

func (s *CardService) GetAll(ctx context.Context) ([]server.CardData, error) {
	user := server.UserFromContext(ctx)

	cards, err := s.qs.SelectCards(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("get all: %w", err)
	}

	var result []server.CardData
	for _, card := range cards {
		result = append(result, server.CardData{
			ID:         card.ID,
			Name:       card.Name,
			Number:     card.Number,
			ExpDate:    card.Cvv,
			Cardholder: card.Cardholder,
			Notes:      stringOrEmpty(card.Notes),
			User:       card.User,
		})
	}

	return result, nil
}

func (s *CardService) Update(ctx context.Context, id int64, data server.CardDataUpdate) error {
	card, err := s.qs.SelectCard(ctx, id)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	if err := server.VerifyCanEditData(ctx, card); err != nil {
		return err
	}

	params := sqlc.UpdateCardParams{
		Name:       card.Name,
		Number:     card.Number,
		ExpDate:    card.ExpDate,
		Cvv:        card.Cvv,
		Cardholder: card.Cardholder,
		Notes:      card.Notes,
		ID:         card.ID,
	}

	if data.Name != nil {
		params.Name = *data.Name
	}
	if data.Number != nil {
		params.Number = *data.Number
	}
	if data.ExpDate != nil {
		params.ExpDate = *data.ExpDate
	}
	if data.CVV != nil {
		params.Cvv = *data.CVV
	}
	if data.Cardholder != nil {
		params.Cardholder = *data.Cardholder
	}
	if data.Notes != nil {
		params.Notes = data.Notes
	}

	n, err := s.qs.UpdateCard(ctx, params)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("update: no rows")
	}
	return nil
}

func (s *CardService) Remove(ctx context.Context, id int64) error {
	return removeData(ctx, s.qs.SelectCardUser, s.qs.DeleteCard, id)
}
