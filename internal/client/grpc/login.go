package grpc

import (
	"context"
	"github.com/mkolibaba/gophkeeper/internal/client"
	pb "github.com/mkolibaba/gophkeeper/internal/common/grpc/proto/gen"
)

type LoginService struct {
	client pb.DataServiceClient
}

func NewLoginService(client pb.DataServiceClient) *LoginService {
	return &LoginService{
		client: client,
	}
}

func (l *LoginService) Save(ctx context.Context, user string, data client.LoginData) error {
	var login pb.Login
	login.SetName(data.Name)
	login.SetLogin(data.Login)
	login.SetPassword(data.Password)
	login.SetMetadata(data.Metadata)

	var in pb.SaveDataRequest
	in.SetUser(user)
	in.SetLogin(&login)

	_, err := l.client.Save(ctx, &in)
	return err
}

func (l *LoginService) GetAll(ctx context.Context, user string) ([]client.LoginData, error) {
	var in pb.GetAllDataRequest
	in.SetDataType(pb.DataType_LOGIN)
	in.SetUser(user)

	result, err := l.client.GetAll(ctx, &in)
	if err != nil {
		return nil, err
	}

	var logins []client.LoginData
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
}

func (l *LoginService) Remove(ctx context.Context, name string, user string) error {
	var in pb.RemoveDataRequest
	in.SetUser(user)
	in.SetName(name)
	in.SetDataType(pb.DataType_LOGIN)

	_, err := l.client.Remove(ctx, &in)
	return err
}
