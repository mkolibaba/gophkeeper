package grpc

import (
	"context"
	"fmt"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type AuthorizationService struct {
	client gophkeeperv1.AuthorizationServiceClient
}

func NewAuthorizationService(conn *grpc.ClientConn) *AuthorizationService {
	return &AuthorizationService{
		client: gophkeeperv1.NewAuthorizationServiceClient(conn),
	}
}

func (s *AuthorizationService) Authorize(ctx context.Context, login string, password string) error {
	var in gophkeeperv1.AuthorizationRequest
	in.SetLogin(login)
	in.SetPassword(password)

	_, err := s.client.Authorize(ctx, &in)

	if err != nil {
		if statusErr, ok := status.FromError(err); ok {
			return fmt.Errorf("%s", statusErr.Message())
		}
		return err
	}

	return nil
}

func (s *AuthorizationService) Register(ctx context.Context, login string, password string) error {
	var in gophkeeperv1.AuthorizationRequest
	in.SetLogin(login)
	in.SetPassword(password)

	_, err := s.client.Register(ctx, &in)

	if err != nil {
		if statusErr, ok := status.FromError(err); ok {
			return fmt.Errorf("%s", statusErr.Message())
		}
		return err
	}

	return nil
}
