package sqlite

import (
	"context"
	"fmt"
	"github.com/mkolibaba/gophkeeper/server"
	"github.com/mkolibaba/gophkeeper/server/sqlite/converter"
	sqlc "github.com/mkolibaba/gophkeeper/server/sqlite/sqlc/gen"
)

type CardService struct {
	qs        *sqlc.Queries
	converter converter.DataConverter
}

func NewCardService(queries *sqlc.Queries, converter converter.DataConverter) *CardService {
	return &CardService{
		qs:        queries,
		converter: converter,
	}
}

func (s *CardService) Create(ctx context.Context, data server.CardData) error {
	err := s.qs.InsertCard(ctx, s.converter.ConvertToInsertCard(ctx, data))
	return unwrapInsertError(err)
}

func (s *CardService) GetAll(ctx context.Context) ([]server.CardData, error) {
	return getAllData(ctx, s.qs.SelectCards, s.converter.ConvertToCardDataSlice)
}

func (s *CardService) Update(ctx context.Context, id int64, data server.CardDataUpdate) error {
	card, err := s.qs.SelectCard(ctx, id)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	if err := server.VerifyCanEditData(ctx, card); err != nil {
		return err
	}

	params := s.converter.ConvertToUpdateCard(card)
	s.converter.ConvertToUpdateCardUpdate(data, &params)

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
