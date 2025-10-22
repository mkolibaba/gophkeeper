package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mkolibaba/gophkeeper/client"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"
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

func (s *LoginService) Save(ctx context.Context, data client.LoginData) error {
	var login gophkeeperv1.Login
	login.SetName(data.Name)
	login.SetLogin(data.Login)
	login.SetPassword(data.Password)
	login.SetWebsite(data.Website)
	login.SetNotes(data.Notes)

	_, err := s.client.Save(ctx, &login)
	return err
}

func (s *LoginService) GetAll(ctx context.Context) ([]client.LoginData, error) {
	result, err := s.client.GetAll(ctx, &empty.Empty{})
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

func (s *LoginService) Update(ctx context.Context, data client.LoginDataUpdate) error {
	var in gophkeeperv1.Login
	in.SetId(data.ID)
	if data.Name != nil {
		in.SetName(*data.Name)
	}
	if data.Login != nil {
		in.SetLogin(*data.Login)
	}
	if data.Password != nil {
		in.SetPassword(*data.Password)
	}
	if data.Website != nil {
		in.SetWebsite(*data.Website)
	}
	if data.Notes != nil {
		in.SetNotes(*data.Notes)
	}

	_, err := s.client.Update(ctx, &in)
	return err
}

func (s *LoginService) Remove(ctx context.Context, id int64) error {
	var in gophkeeperv1.RemoveDataRequest
	in.SetId(id)

	_, err := s.client.Remove(ctx, &in)
	return err
}
