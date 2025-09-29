package grpc

import (
	"context"
	"github.com/mkolibaba/gophkeeper/internal/client"
	pb "github.com/mkolibaba/gophkeeper/internal/common/grpc/proto/gen"
)

type CardService struct {
	client pb.DataServiceClient
}

func NewCardService(client pb.DataServiceClient) *CardService {
	return &CardService{
		client: client,
	}
}

func (c *CardService) Save(ctx context.Context, user string, data client.CardData) error {
	var card pb.Card
	card.SetName(data.Name)
	card.SetNumber(data.Number)
	card.SetExpDate(data.ExpDate)
	card.SetCvv(data.CVV)
	card.SetCardholder(data.Cardholder)
	card.SetMetadata(data.Metadata)

	var in pb.SaveDataRequest
	in.SetUser(user)
	in.SetCard(&card)

	_, err := c.client.Save(ctx, &in)
	return err
}

func (c *CardService) GetAll(ctx context.Context, user string) ([]client.CardData, error) {
	var in pb.GetAllDataRequest
	in.SetDataType(pb.DataType_CARD)
	in.SetUser(user)

	result, err := c.client.GetAll(ctx, &in)
	if err != nil {
		return nil, err
	}

	var cards []client.CardData
	for _, data := range result.GetData() {
		card := data.GetCard()
		cards = append(cards, client.CardData{
			Name:       card.GetName(),
			Number:     card.GetNumber(),
			ExpDate:    card.GetExpDate(),
			CVV:        card.GetCvv(),
			Cardholder: card.GetCardholder(),
			Metadata:   card.GetMetadata(),
		})
	}
	return cards, nil
}

func (c *CardService) Remove(ctx context.Context, name string, user string) error {
	var in pb.RemoveDataRequest
	in.SetUser(user)
	in.SetName(name)
	in.SetDataType(pb.DataType_CARD)

	_, err := c.client.Remove(ctx, &in)
	return err
}
