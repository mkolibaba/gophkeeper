package sqlite

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mkolibaba/gophkeeper/internal/server"
	sqlc "github.com/mkolibaba/gophkeeper/internal/server/sqlite/sqlc/gen"
)

type CardService struct {
	qs *sqlc.Queries
}

func (c *CardService) Save(ctx context.Context, data server.CardData) error {
	metadata, err := json.Marshal(data.Metadata)
	if err != nil {
		return fmt.Errorf("save: invalid metadata: %w", err)
	}

	err = c.qs.SaveCard(ctx, sqlc.SaveCardParams{
		Name:       data.Name,
		Number:     data.Number,
		ExpDate:    data.ExpDate,
		Cvv:        data.CVV,
		Cardholder: data.Cardholder,
		Metadata:   metadata,
		User:       data.User,
	})

	return tryUnwrapSaveError(err)
}

func (c *CardService) GetAll(ctx context.Context, user string) ([]server.CardData, error) {
	cards, err := c.qs.GetAllCards(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("get all: %w", err)
	}

	var result []server.CardData
	for _, card := range cards {
		metadata, err := unmarshalMetadata(card.Metadata)
		if err != nil {
			return nil, fmt.Errorf("get all: %w", err)
		}

		result = append(result, server.CardData{
			User:       card.User,
			Name:       card.Name,
			Number:     card.Number,
			ExpDate:    card.ExpDate,
			CVV:        card.Cvv,
			Cardholder: card.Cardholder,
			Metadata:   metadata,
		})
	}

	return result, nil
}

func (c *CardService) Remove(ctx context.Context, name string, user string) error {
	n, err := c.qs.RemoveCard(ctx, name)
	if err != nil {
		return fmt.Errorf("remove: %w", err)
	}
	if n == 0 {
		return server.ErrDataNotFound
	}
	return nil
}
