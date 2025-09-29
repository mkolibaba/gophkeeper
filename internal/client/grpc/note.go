package grpc

import (
	"context"
	"github.com/mkolibaba/gophkeeper/internal/client"
	pb "github.com/mkolibaba/gophkeeper/internal/common/grpc/proto/gen"
)

type NoteService struct {
	client pb.DataServiceClient
}

func NewNoteService(client pb.DataServiceClient) *NoteService {
	return &NoteService{
		client: client,
	}
}

func (n *NoteService) Save(ctx context.Context, user string, data client.NoteData) error {
	var note pb.Note
	note.SetName(data.Name)
	note.SetText(data.Text)
	note.SetMetadata(data.Metadata)

	var in pb.SaveDataRequest
	in.SetUser(user)
	in.SetNote(&note)

	_, err := n.client.Save(ctx, &in)
	return err
}

func (n *NoteService) GetAll(ctx context.Context, user string) ([]client.NoteData, error) {
	var in pb.GetAllDataRequest
	in.SetDataType(pb.DataType_NOTE)
	in.SetUser(user)

	result, err := n.client.GetAll(ctx, &in)
	if err != nil {
		return nil, err
	}

	var notes []client.NoteData
	for _, data := range result.GetData() {
		note := data.GetNote()
		notes = append(notes, client.NoteData{
			Name:     note.GetName(),
			Text:     note.GetText(),
			Metadata: note.GetMetadata(),
		})
	}
	return notes, nil
}

func (n *NoteService) Remove(ctx context.Context, name string, user string) error {
	var in pb.RemoveDataRequest
	in.SetUser(user)
	in.SetName(name)
	in.SetDataType(pb.DataType_NOTE)

	_, err := n.client.Remove(ctx, &in)
	return err
}
