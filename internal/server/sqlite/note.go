package sqlite

import (
	"context"
	"fmt"
	"github.com/mkolibaba/gophkeeper/internal/server"
	sqlc "github.com/mkolibaba/gophkeeper/internal/server/sqlite/sqlc/gen"
)

type NoteService struct {
	qs *sqlc.Queries
}

func NewNoteService(queries *sqlc.Queries) *NoteService {
	return &NoteService{
		qs: queries,
	}
}

func (n *NoteService) Save(ctx context.Context, data server.NoteData, user string) error {
	err := n.qs.SaveNote(ctx, sqlc.SaveNoteParams{
		Name: data.Name,
		Text: stringOrNull(data.Text),
		User: user,
	})

	return tryUnwrapSaveError(err)
}

func (n *NoteService) GetAll(ctx context.Context, user string) ([]server.NoteData, error) {
	notes, err := n.qs.GetAllNotes(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("get all: %w", err)
	}

	var result []server.NoteData
	for _, note := range notes {
		result = append(result, server.NoteData{
			Name: note.Name,
			Text: stringOrEmpty(note.Text),
		})
	}

	return result, nil
}

func (n *NoteService) Remove(ctx context.Context, name string, user string) error {
	cnt, err := n.qs.RemoveNote(ctx, name)
	if err != nil {
		return fmt.Errorf("remove: %w", err)
	}
	if cnt == 0 {
		return server.ErrDataNotFound
	}
	return nil
}
