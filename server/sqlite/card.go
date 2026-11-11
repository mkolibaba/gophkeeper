package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mkolibaba/gophkeeper/server"
	"github.com/mkolibaba/gophkeeper/server/sqlite/converter"
	sqlc "github.com/mkolibaba/gophkeeper/server/sqlite/sqlc/gen"
)

var _ server.CardService = (*CardService)(nil)

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

func (s *CardService) Create(ctx context.Context, data server.CardData) (int64, time.Time, error) {
	result, err := s.qs.InsertCard(ctx, s.converter.ConvertToInsertCard(ctx, data))
	return result.ID, time.Unix(result.UpdatedAt, 0), unwrapInsertError(err)
}

func (s *CardService) GetAll(ctx context.Context) ([]server.CardData, error) {
	return getAllData(ctx, s.qs.SelectCards, s.converter.ConvertToCardDataSlice)
}

func (s *CardService) Update(ctx context.Context, id int64, data server.CardDataUpdate) (time.Time, error) {
	card, err := s.qs.SelectCard(ctx, id, server.UserFromContext(ctx))
	if errors.Is(err, sql.ErrNoRows) {
		return time.Time{}, server.ErrDataNotFound
	}
	if err != nil {
		return time.Time{}, fmt.Errorf("update: %w", err)
	}

	params := s.converter.ConvertToUpdateCard(card)
	s.converter.ConvertToUpdateCardUpdate(data, &params)

	updatedAt, err := s.qs.UpdateCard(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return time.Time{}, server.ErrDataNotFound
		}
		return time.Time{}, fmt.Errorf("update: %w", err)
	}
	return time.Unix(updatedAt, 0), nil
}

func (s *CardService) Remove(ctx context.Context, id int64) error {
	return removeData(ctx, s.qs.DeleteCard, id)
}
