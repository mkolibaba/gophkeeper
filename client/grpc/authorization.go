package grpc

import (
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeperv1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type AuthorizationService struct {
	client gophkeeperv1.AuthorizationServiceClient
	logger *log.Logger
}

func NewAuthorizationService(conn *grpc.ClientConn, logger *log.Logger) *AuthorizationService {
	return &AuthorizationService{
		client: gophkeeperv1.NewAuthorizationServiceClient(conn),
		logger: logger,
	}
}

func (s *AuthorizationService) Authorize(ctx context.Context, login string, password string) (string, error) {
	s.logger.Debug("trying to authorize", "login", login, "password", password)
	return s.send(ctx, s.client.Authorize, login, password)
}

func (s *AuthorizationService) Register(ctx context.Context, login string, password string) (string, error) {
	return s.send(ctx, s.client.Register, login, password)
}

func (s *AuthorizationService) send(
	ctx context.Context,
	sender func(ctx context.Context, in *gophkeeperv1.UserCredentials, opts ...grpc.CallOption) (*gophkeeperv1.TokenResponse, error),
	login, password string,
) (string, error) {
	var in gophkeeperv1.UserCredentials
	in.SetLogin(login)
	in.SetPassword(password)

	out, err := sender(ctx, &in)

	if err != nil {
		if statusErr, ok := status.FromError(err); ok {
			return "", fmt.Errorf("%s", statusErr.Message())
		}
		return "", err
	}

	return out.GetToken(), nil
}
