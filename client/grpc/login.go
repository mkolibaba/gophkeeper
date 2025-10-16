package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mkolibaba/gophkeeper/client"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeper"
	"google.golang.org/grpc"
)

type LoginService struct {
	client gophkeeperv1.LoginServiceClient
}

func NewLoginService(conn *grpc.ClientConn) *LoginService {
	return &LoginService{
		client: gophkeeperv1.NewLoginServiceClient(conn),
	}
}

func (l *LoginService) Save(ctx context.Context, data client.LoginData) error {
	var login gophkeeperv1.Login
	login.SetName(data.Name)
	login.SetLogin(data.Login)
	login.SetPassword(data.Password)
	login.SetWebsite(data.Website)
	login.SetNotes(data.Notes)

	_, err := l.client.Save(ctx, &login)
	return err
}

func (l *LoginService) GetAll(ctx context.Context) ([]client.LoginData, error) {
	result, err := l.client.GetAll(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	var logins []client.LoginData
	for _, data := range result.GetResult() {
		logins = append(logins, client.LoginData{
			ID:       data.GetId(),
			Name:     data.GetName(),
			Login:    data.GetLogin(),
			Password: data.GetPassword(),
			Website:  data.GetWebsite(),
			Notes:    data.GetNotes(),
		})
	}

	return logins, nil
}

func (l *LoginService) Remove(ctx context.Context, id int64) error {
	var in gophkeeperv1.RemoveDataRequest
	in.SetId(id)

	_, err := l.client.Remove(ctx, &in)
	return err
}
