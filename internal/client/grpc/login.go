package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mkolibaba/gophkeeper/internal/client"
	pb "github.com/mkolibaba/gophkeeper/internal/common/grpc/proto/gen"
	"google.golang.org/grpc"
)

type LoginService struct {
	client pb.LoginServiceClient
}

func NewLoginService(conn *grpc.ClientConn) *LoginService {
	return &LoginService{
		client: pb.NewLoginServiceClient(conn),
	}
}

func (l *LoginService) Save(ctx context.Context, user string, data client.LoginData) error {
	var login pb.Login
	login.SetName(data.Name)
	login.SetLogin(data.Login)
	login.SetPassword(data.Password)
	login.SetMetadata(data.Metadata)

	_, err := l.client.Save(ctx, &login)
	return err
}

func (l *LoginService) GetAll(ctx context.Context, user string) ([]client.LoginData, error) {
	result, err := l.client.GetAll(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	var logins []client.LoginData
	for _, data := range result.GetResult() {
		logins = append(logins, client.LoginData{
			Name:     data.GetName(),
			Login:    data.GetLogin(),
			Password: data.GetPassword(),
			Metadata: data.GetMetadata(),
		})
	}

	return logins, nil
}

func (l *LoginService) Remove(ctx context.Context, name string, user string) error {
	var in pb.RemoveDataRequest
	in.SetName(name)

	_, err := l.client.Remove(ctx, &in)
	return err
}
