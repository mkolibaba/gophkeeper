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

var _ server.NoteService = (*NoteService)(nil)

type NoteService struct {
	qs        *sqlc.Queries
	converter converter.DataConverter
}

func NewNoteService(queries *sqlc.Queries, converter converter.DataConverter) *NoteService {
	return &NoteService{
		qs:        queries,
		converter: converter,
	}
}

func (s *NoteService) Create(ctx context.Context, data server.NoteData) (int64, time.Time, error) {
	result, err := s.qs.InsertNote(ctx, s.converter.ConvertToInsertNote(ctx, data))
	return result.ID, time.Unix(result.UpdatedAt, 0), unwrapInsertError(err)
}

func (s *NoteService) GetAll(ctx context.Context) ([]server.NoteData, error) {
	return getAllData(ctx, s.qs.SelectNotes, s.converter.ConvertToNoteDataSlice)
}

func (s *NoteService) Update(ctx context.Context, id int64, data server.NoteDataUpdate) (time.Time, error) {
	note, err := s.qs.SelectNote(ctx, id, server.UserFromContext(ctx))
	if errors.Is(err, sql.ErrNoRows) {
		return time.Time{}, server.ErrDataNotFound
	}
	if err != nil {
		return time.Time{}, fmt.Errorf("update: %w", err)
	}

	params := s.converter.ConvertToUpdateNote(note)
	s.converter.ConvertToUpdateNoteUpdate(data, &params)

	updatedAt, err := s.qs.UpdateNote(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return time.Time{}, server.ErrDataNotFound
		}
		return time.Time{}, fmt.Errorf("update: %w", err)
	}
	return time.Unix(updatedAt, 0), nil
}

func (s *NoteService) Remove(ctx context.Context, id int64) error {
	return removeData(ctx, s.qs.DeleteNote, id)
}
