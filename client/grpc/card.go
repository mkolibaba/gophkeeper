package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mkolibaba/gophkeeper/client"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"
	"google.golang.org/grpc"
)

type CardService struct {
	client gophkeeperv1.CardServiceClient
}

func NewCardService(conn *grpc.ClientConn) *CardService {
	return &CardService{
		client: gophkeeperv1.NewCardServiceClient(conn),
	}
}

func (s *CardService) Save(ctx context.Context, data client.CardData) error {
	var card gophkeeperv1.Card
	card.SetName(data.Name)
	card.SetNumber(data.Number)
	card.SetExpDate(data.ExpDate)
	card.SetCvv(data.CVV)
	card.SetCardholder(data.Cardholder)
	card.SetNotes(data.Notes)

	_, err := s.client.Save(ctx, &card)
	return err
}

func (s *CardService) GetAll(ctx context.Context) ([]client.CardData, error) {
	result, err := s.client.GetAll(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	var cards []client.CardData
	for _, data := range result.GetResult() {
		cards = append(cards, client.CardData{
			ID:         data.GetId(),
			Name:       data.GetName(),
			Number:     data.GetNumber(),
			ExpDate:    data.GetExpDate(),
			CVV:        data.GetCvv(),
			Cardholder: data.GetCardholder(),
			Notes:      data.GetNotes(),
		})
	}
	return cards, nil
}

func (s *CardService) Update(ctx context.Context, data client.CardDataUpdate) error {
	var in gophkeeperv1.Card
	in.SetId(data.ID)
	if data.Name != nil {
		in.SetName(*data.Name)
	}
	if data.Number != nil {
		in.SetNumber(*data.Number)
	}
	if data.ExpDate != nil {
		in.SetExpDate(*data.ExpDate)
	}
	if data.CVV != nil {
		in.SetCvv(*data.CVV)
	}
	if data.Cardholder != nil {
		in.SetCardholder(*data.Cardholder)
	}
	if data.Notes != nil {
		in.SetNotes(*data.Notes)
	}

	_, err := s.client.Update(ctx, &in)
	return err
}

func (s *CardService) Remove(ctx context.Context, id int64) error {
	var in gophkeeperv1.RemoveDataRequest
	in.SetId(id)

	_, err := s.client.Remove(ctx, &in)
	return err
}
