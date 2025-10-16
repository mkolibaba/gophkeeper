package sqlite

import (
	"context"
	"fmt"
	"github.com/mkolibaba/gophkeeper/server"
	sqlc "github.com/mkolibaba/gophkeeper/server/sqlite/sqlc/gen"
)

type NoteService struct {
	qs *sqlc.Queries
}

func NewNoteService(queries *sqlc.Queries) *NoteService {
	return &NoteService{
		qs: queries,
	}
}

func (s *NoteService) Create(ctx context.Context, data server.NoteData) error {
	err := s.qs.InsertNote(ctx, sqlc.InsertNoteParams{
		Name: data.Name,
		Text: stringOrNull(data.Text),
		User: server.UserFromContext(ctx),
	})
	return unwrapInsertError(err)
}

func (s *NoteService) GetAll(ctx context.Context) ([]server.NoteData, error) {
	user := server.UserFromContext(ctx)

	notes, err := s.qs.SelectNotes(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("get all: %w", err)
	}

	var result []server.NoteData
	for _, note := range notes {
		result = append(result, server.NoteData{
			Name: note.Name,
			Text: stringOrEmpty(note.Text),
			ID:   note.ID,
		})
	}

	return result, nil
}

func (s *NoteService) Update(ctx context.Context, id int64, data server.NoteDataUpdate) error {
	note, err := s.qs.SelectNote(ctx, id)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	if err := server.VerifyCanEditData(ctx, note); err != nil {
		return err
	}

	params := sqlc.UpdateNoteParams{
		Name: note.Name,
		Text: note.Text,
		ID:   note.ID,
	}

	if data.Text != nil {
		params.Text = data.Text
	}

	n, err := s.qs.UpdateNote(ctx, params)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("update: no rows")
	}
	return nil
}

func (s *NoteService) Remove(ctx context.Context, id int64) error {
	return removeData(ctx, s.qs.SelectNoteUser, s.qs.DeleteNote, id)
}
