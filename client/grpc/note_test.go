package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mkolibaba/gophkeeper/client"
	"github.com/mkolibaba/gophkeeper/client/grpc/mock"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"testing"
)

func TestNoteSave(t *testing.T) {
	clientMock := &mock.NoteServiceClientMock{}
	srv := NewNoteService(clientMock)

	err := srv.Save(t.Context(), client.NoteData{
		Name: "new name",
		Text: "new text",
	})
	require.NoError(t, err)

	cc := clientMock.SaveCalls()
	require.Len(t, cc, 1)
	c := cc[0]
	require.Equal(t, c.In.GetName(), "new name")
	require.Equal(t, c.In.GetText(), "new text")
}

func TestNoteGetAll(t *testing.T) {
	clientMock := &mock.NoteServiceClientMock{
		GetAllFunc: func(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*gophkeeperv1.GetAllNotesResponse, error) {
			var note1 gophkeeperv1.Note
			note1.SetId(1)
			note1.SetName("name 1")
			note1.SetText("note 1")
			var out gophkeeperv1.GetAllNotesResponse
			out.SetResult([]*gophkeeperv1.Note{&note1})
			return &out, nil
		},
	}
	srv := NewNoteService(clientMock)

	all, err := srv.GetAll(t.Context())
	require.NoError(t, err)

	require.Len(t, all, 1)
	c := all[0]
	require.Equal(t, c.Name, "name 1")
	require.Equal(t, c.Text, "note 1")
}

func TestNoteUpdate(t *testing.T) {
	clientMock := &mock.NoteServiceClientMock{}
	srv := NewNoteService(clientMock)

	name := "new name"
	err := srv.Update(t.Context(), client.NoteDataUpdate{
		Name: &name,
	})
	require.NoError(t, err)

	cc := clientMock.UpdateCalls()
	require.Len(t, cc, 1)
	c := cc[0]
	require.Equal(t, c.In.GetName(), "new name")
	require.False(t, c.In.HasText())
}
