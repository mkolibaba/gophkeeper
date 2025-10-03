package grpc

import (
	"context"
	pb "github.com/mkolibaba/gophkeeper/internal/common/grpc/proto/gen"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthorizationService struct {
	pb.UnimplementedAuthorizationServiceServer
	logger *zap.Logger
}

func NewAuthorizationService(logger *zap.Logger) *AuthorizationService {
	return &AuthorizationService{
		logger: logger,
	}
}

func (s *AuthorizationService) Authorize(ctx context.Context, in *pb.AuthorizationRequest) (*pb.AuthorizationResponse, error) {
	// TODO: сделать нормальную реализацию
	if in.GetLogin() == "demo" && in.GetPassword() == "demo" {
		var out pb.AuthorizationResponse
		out.SetToken("cool token")
		return &out, nil
	}

	return nil, status.Error(codes.Unauthenticated, "invalid login or password")
}
