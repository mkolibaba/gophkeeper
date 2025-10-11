package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mkolibaba/gophkeeper/internal/client"
	pb "github.com/mkolibaba/gophkeeper/internal/common/grpc/proto/gen"
	"google.golang.org/grpc"
)

type NoteService struct {
	client pb.NoteServiceClient
}

func NewNoteService(conn *grpc.ClientConn) *NoteService {
	return &NoteService{
		client: pb.NewNoteServiceClient(conn),
	}
}

func (n *NoteService) Save(ctx context.Context, data client.NoteData) error {
	var note pb.Note
	note.SetName(data.Name)
	note.SetText(data.Text)

	_, err := n.client.Save(ctx, &note)
	return err
}

func (n *NoteService) GetAll(ctx context.Context) ([]client.NoteData, error) {
	result, err := n.client.GetAll(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	var notes []client.NoteData
	for _, data := range result.GetResult() {
		notes = append(notes, client.NoteData{
			Name: data.GetName(),
			Text: data.GetText(),
		})
	}
	return notes, nil
}

func (n *NoteService) Remove(ctx context.Context, name string) error {
	var in pb.RemoveDataRequest
	in.SetName(name)

	_, err := n.client.Remove(ctx, &in)
	return err
}
