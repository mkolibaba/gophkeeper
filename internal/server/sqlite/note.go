package sqlite

import (
	"context"
	"encoding/json"
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

func (n *NoteService) Save(ctx context.Context, data server.NoteData) error {
	metadata, err := json.Marshal(data.Metadata)
	if err != nil {
		return fmt.Errorf("save: invalid metadata: %w", err)
	}

	err = n.qs.SaveNote(ctx, sqlc.SaveNoteParams{
		Name:     data.Name,
		Text:     &data.Text,
		Metadata: metadata,
		User:     data.User,
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
		metadata, err := unmarshalMetadata(note.Metadata)
		if err != nil {
			return nil, fmt.Errorf("get all: %w", err)
		}

		result = append(result, server.NoteData{
			User:     note.User,
			Name:     note.Name,
			Text:     *note.Text,
			Metadata: metadata,
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
