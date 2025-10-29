package grpc

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"
	"github.com/mkolibaba/gophkeeper/server"
	"github.com/mkolibaba/gophkeeper/server/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestNoteSave(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := createNoteServiceServer(t, &mock.NoteServiceMock{})

		var in gophkeeperv1.Note
		in.SetName("new note")
		in.SetText("some text")

		_, err := srv.Save(t.Context(), &in)
		require.NoError(t, err)
	})
	t.Run("validation_error", func(t *testing.T) {
		srv := createNoteServiceServer(t, &mock.NoteServiceMock{})

		var in gophkeeperv1.Note

		_, err := srv.Save(t.Context(), &in)
		requireGrpcError(t, err, codes.InvalidArgument)
	})
	t.Run("db_error", func(t *testing.T) {
		noteServiceMock := &mock.NoteServiceMock{
			CreateFunc: func(_ context.Context, _ server.NoteData) error {
				return fmt.Errorf("some error")
			},
		}

		srv := createNoteServiceServer(t, noteServiceMock)

		var in gophkeeperv1.Note
		in.SetName("new note")
		in.SetText("some text")

		_, err := srv.Save(t.Context(), &in)
		requireGrpcError(t, err, codes.Internal)
	})
}

func TestNoteUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := createNoteServiceServer(t, &mock.NoteServiceMock{})

		var in gophkeeperv1.Note
		in.SetId(1)
		in.SetName("new note name")
		in.SetText("updated text")

		_, err := srv.Update(t.Context(), &in)
		require.NoError(t, err)
	})
	t.Run("validation_error", func(t *testing.T) {
		srv := createNoteServiceServer(t, &mock.NoteServiceMock{})

		var in gophkeeperv1.Note
		in.SetName("new note name")

		_, err := srv.Update(t.Context(), &in)
		requireGrpcError(t, err, codes.InvalidArgument)
	})
}

func TestNoteRemove(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := createNoteServiceServer(t, &mock.NoteServiceMock{})

		var in gophkeeperv1.RemoveDataRequest
		in.SetId(1)

		_, err := srv.Remove(t.Context(), &in)
		require.NoError(t, err)
	})
	t.Run("validation_error", func(t *testing.T) {
		srv := createNoteServiceServer(t, &mock.NoteServiceMock{})

		var in gophkeeperv1.RemoveDataRequest

		_, err := srv.Remove(t.Context(), &in)
		requireGrpcError(t, err, codes.InvalidArgument)
	})
	t.Run("not_found", func(t *testing.T) {
		service := &mock.NoteServiceMock{
			RemoveFunc: func(ctx context.Context, id int64) error {
				return server.ErrDataNotFound
			},
		}

		srv := createNoteServiceServer(t, service)

		var in gophkeeperv1.RemoveDataRequest
		in.SetId(1)

		_, err := srv.Remove(t.Context(), &in)
		requireGrpcError(t, err, codes.NotFound)
	})
}

func TestNoteGetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := &mock.NoteServiceMock{
			GetAllFunc: func(ctx context.Context) ([]server.NoteData, error) {
				return []server.NoteData{
					{ID: 1, Name: "note1"},
					{ID: 2, Name: "note2"},
				}, nil
			},
		}
		srv := createNoteServiceServer(t, service)
		resp, err := srv.GetAll(t.Context(), nil)
		require.NoError(t, err)
		require.Len(t, resp.GetResult(), 2)
	})
	t.Run("db_error", func(t *testing.T) {
		service := &mock.NoteServiceMock{
			GetAllFunc: func(ctx context.Context) ([]server.NoteData, error) {
				return nil, fmt.Errorf("db error")
			},
		}
		srv := createNoteServiceServer(t, service)
		_, err := srv.GetAll(t.Context(), nil)
		requireGrpcError(t, err, codes.Internal)
	})
}

func createNoteServiceServer(t *testing.T, noteService server.NoteService) *NoteServiceServer {
	return NewNoteServiceServer(noteService, newTestValidator(t), log.New(io.Discard))
}
