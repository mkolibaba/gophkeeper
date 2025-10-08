package grpc

import (
	"context"
	"fmt"
	pb "github.com/mkolibaba/gophkeeper/internal/common/grpc/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type AuthorizationService struct {
	client pb.AuthorizationServiceClient
}

func NewAuthorizationService(conn *grpc.ClientConn) *AuthorizationService {
	return &AuthorizationService{
		client: pb.NewAuthorizationServiceClient(conn),
	}
}

func (s *AuthorizationService) Authorize(ctx context.Context, login string, password string) (string, error) {
	var in pb.AuthorizationRequest
	in.SetLogin(login)
	in.SetPassword(password)

	out, err := s.client.Authorize(ctx, &in)

	if err != nil {
		if statusErr, ok := status.FromError(err); ok {
			return "", fmt.Errorf("%s", statusErr.Message())
		}
		return "", err
	}

	return out.GetToken(), nil
}

func (s *AuthorizationService) Register(ctx context.Context, login string, password string) (string, error) {
	var in pb.AuthorizationRequest
	in.SetLogin(login)
	in.SetPassword(password)

	out, err := s.client.Register(ctx, &in)

	if err != nil {
		if statusErr, ok := status.FromError(err); ok {
			return "", fmt.Errorf("%s", statusErr.Message())
		}
		return "", err
	}

	return out.GetToken(), nil
}
