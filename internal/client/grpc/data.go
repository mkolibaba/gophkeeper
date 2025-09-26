package grpc

import (
	"context"
	"errors"
	"github.com/mkolibaba/gophkeeper/internal/client"
	pb "github.com/mkolibaba/gophkeeper/internal/common/grpc/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var clientToGRPCDataType = map[client.DataType]pb.DataType{
	client.DataTypeLogin:  pb.DataType_LOGIN,
	client.DataTypeNote:   pb.DataType_NOTE,
	client.DataTypeBinary: pb.DataType_BINARY,
	client.DataTypeCard:   pb.DataType_CARD,
}

type DataService struct {
	conn *grpc.ClientConn
	c    pb.DataServiceClient
}

func NewDataService(serverAddress string) (*DataService, error) {
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := pb.NewDataServiceClient(conn)

	return &DataService{
		conn: conn,
		c:    c,
	}, nil
}

func (d *DataService) Save(ctx context.Context, user string, data client.Data) error {
	var in pb.SaveDataRequest
	in.SetUser(user)

	// TODO: validation
	switch data := data.(type) {
	case client.LoginData:
		var login pb.Login
		login.SetName(data.Name)
		login.SetLogin(data.Login)
		login.SetPassword(data.Password)
		login.SetMetadata(data.Metadata)
		in.SetLogin(&login)
	case client.NoteData:
		var note pb.Note
		note.SetName(data.Name)
		note.SetText(data.Text)
		note.SetMetadata(data.Metadata)
		in.SetNote(&note)
	case client.BinaryData:
		var binary pb.Binary
		binary.SetName(data.Name)
		binary.SetData(data.Data)
		binary.SetMetadata(data.Metadata)
		in.SetBinary(&binary)
	case client.CardData:
		var card pb.Card
		card.SetName(data.Name)
		card.SetNumber(data.Number)
		card.SetExpDate(data.ExpDate)
		card.SetCvv(data.CVV)
		card.SetCardholder(data.Cardholder)
		card.SetMetadata(data.Metadata)
		in.SetCard(&card)
	default:
		return errors.New("invalid data type")
	}

	_, err := d.c.Save(ctx, &in)
	return err
}

func (d *DataService) GetAll(ctx context.Context, user string, dataType client.DataType) ([]client.Data, error) {
	var in pb.GetAllDataRequest
	in.SetDataType(clientToGRPCDataType[dataType]) // TODO: validation
	in.SetUser(user)

	result, err := d.c.GetAll(ctx, &in)
	if err != nil {
		return nil, err
	}

	switch dataType {
	case client.DataTypeLogin:
		var logins []client.Data
		for _, data := range result.GetData() {
			login := data.GetLogin()
			logins = append(logins, client.LoginData{
				Name:     login.GetName(),
				Login:    login.GetLogin(),
				Password: login.GetPassword(),
				Metadata: login.GetMetadata(),
			})
		}
		return logins, nil
	case client.DataTypeNote:
		var notes []client.Data
		for _, data := range result.GetData() {
			note := data.GetNote()
			notes = append(notes, client.NoteData{
				Name:     note.GetName(),
				Text:     note.GetText(),
				Metadata: note.GetMetadata(),
			})
		}
		return notes, nil
	case client.DataTypeBinary:
		var binaries []client.Data
		for _, data := range result.GetData() {
			binary := data.GetBinary()
			binaries = append(binaries, client.BinaryData{
				Name:     binary.GetName(),
				Data:     binary.GetData(),
				Metadata: binary.GetMetadata(),
			})
		}
		return binaries, nil
	case client.DataTypeCard:
		var cards []client.Data
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

	return nil, err
}

func (d *DataService) Remove(ctx context.Context, user string, name string, dataType client.DataType) error {
	var in pb.RemoveDataRequest
	in.SetUser(user)
	in.SetName(name)
	in.SetDataType(clientToGRPCDataType[dataType]) // TODO: validation

	_, err := d.c.Remove(ctx, &in)
	return err
}
