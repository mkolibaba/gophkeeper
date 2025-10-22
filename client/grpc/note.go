package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mkolibaba/gophkeeper/client"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"
	"google.golang.org/grpc"
)

type NoteService struct {
	client gophkeeperv1.NoteServiceClient
}

func NewNoteService(conn *grpc.ClientConn) *NoteService {
	return &NoteService{
		client: gophkeeperv1.NewNoteServiceClient(conn),
	}
}

func (s *NoteService) Save(ctx context.Context, data client.NoteData) error {
	var note gophkeeperv1.Note
	note.SetName(data.Name)
	note.SetText(data.Text)

	_, err := s.client.Save(ctx, &note)
	return err
}

func (s *NoteService) GetAll(ctx context.Context) ([]client.NoteData, error) {
	result, err := s.client.GetAll(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	var notes []client.NoteData
	for _, data := range result.GetResult() {
		notes = append(notes, client.NoteData{
			ID:   data.GetId(),
			Name: data.GetName(),
			Text: data.GetText(),
		})
	}
	return notes, nil
}

func (s *NoteService) Update(ctx context.Context, data client.NoteDataUpdate) error {
	var in gophkeeperv1.Note
	in.SetId(data.ID)
	if data.Name != nil {
		in.SetName(*data.Name)
	}
	if data.Text != nil {
		in.SetText(*data.Text)
	}

	_, err := s.client.Update(ctx, &in)
	return err
}

func (s *NoteService) Remove(ctx context.Context, id int64) error {
	var in gophkeeperv1.RemoveDataRequest
	in.SetId(id)

	_, err := s.client.Remove(ctx, &in)
	return err
}
