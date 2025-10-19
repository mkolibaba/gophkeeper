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

func (c *CardService) Save(ctx context.Context, data client.CardData) error {
	var card gophkeeperv1.Card
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

func (c *CardService) Remove(ctx context.Context, id int64) error {
	var in gophkeeperv1.RemoveDataRequest
	in.SetId(id)

	_, err := c.client.Remove(ctx, &in)
	return err
}
