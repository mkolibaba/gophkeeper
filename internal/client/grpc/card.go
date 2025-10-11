package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mkolibaba/gophkeeper/internal/client"
	pb "github.com/mkolibaba/gophkeeper/internal/common/grpc/proto/gen"
	"google.golang.org/grpc"
)

type CardService struct {
	client pb.CardServiceClient
}

func NewCardService(conn *grpc.ClientConn) *CardService {
	return &CardService{
		client: pb.NewCardServiceClient(conn),
	}
}

func (c *CardService) Save(ctx context.Context, data client.CardData) error {
	var card pb.Card
	card.SetName(data.Name)
	card.SetNumber(data.Number)
	card.SetExpDate(data.ExpDate)
	card.SetCvv(data.CVV)
	card.SetCardholder(data.Cardholder)
	card.SetNotes(data.Notes)

	_, err := c.client.Save(ctx, &card)
	return err
}

func (c *CardService) GetAll(ctx context.Context) ([]client.CardData, error) {
	result, err := c.client.GetAll(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	var cards []client.CardData
	for _, data := range result.GetResult() {
		cards = append(cards, client.CardData{
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

func (c *CardService) Remove(ctx context.Context, name string) error {
	var in pb.RemoveDataRequest
	in.SetName(name)

	_, err := c.client.Remove(ctx, &in)
	return err
}
