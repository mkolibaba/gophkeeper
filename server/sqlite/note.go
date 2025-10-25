package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mkolibaba/gophkeeper/server"
	"github.com/mkolibaba/gophkeeper/server/sqlite/converter"
	sqlc "github.com/mkolibaba/gophkeeper/server/sqlite/sqlc/gen"
)

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

func (s *NoteService) Create(ctx context.Context, data server.NoteData) error {
	_, err := s.qs.InsertNote(ctx, s.converter.ConvertToInsertNote(ctx, data))
	return unwrapInsertError(err)
}

func (s *NoteService) GetAll(ctx context.Context) ([]server.NoteData, error) {
	return getAllData(ctx, s.qs.SelectNotes, s.converter.ConvertToNoteDataSlice)
}

func (s *NoteService) Update(ctx context.Context, id int64, data server.NoteDataUpdate) error {
	note, err := s.qs.SelectNote(ctx, id, server.UserFromContext(ctx))
	if errors.Is(err, sql.ErrNoRows) {
		return server.ErrDataNotFound
	}
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	params := s.converter.ConvertToUpdateNote(note)
	s.converter.ConvertToUpdateNoteUpdate(data, &params)

	n, err := s.qs.UpdateNote(ctx, params)
	if n == 0 {
		return server.ErrDataNotFound
	}
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}

func (s *NoteService) Remove(ctx context.Context, id int64) error {
	return removeData(ctx, s.qs.DeleteNote, id)
}
