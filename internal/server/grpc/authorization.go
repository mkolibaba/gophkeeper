package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/mkolibaba/gophkeeper/internal/common/grpc/proto/gen"
	"github.com/mkolibaba/gophkeeper/internal/server"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthorizationServiceServer struct {
	pb.UnimplementedAuthorizationServiceServer
	userService   server.UserService
	authService   *server.AuthService
	dataValidator *validator.Validate
	logger        *zap.Logger
}

func NewAuthorizationServiceServer(
	userService server.UserService,
	authService *server.AuthService,
	dataValidator *validator.Validate,
	logger *zap.Logger,
) *AuthorizationServiceServer {
	// TODO: стоит разделить
	rules := map[string]string{
		"login":    "required",
		"password": "required",
	}

	dataValidator.RegisterStructValidationMapRules(rules, pb.AuthorizationRequest{})

	return &AuthorizationServiceServer{
		userService:   userService,
		authService:   authService,
		dataValidator: dataValidator,
		logger:        logger,
	}
}

// TODO: можно убрать токен и просто отправлять статус ок/не ок
func (s *AuthorizationServiceServer) Authorize(ctx context.Context, in *pb.AuthorizationRequest) (*empty.Empty, error) {
	if err := s.dataValidator.StructCtx(ctx, in); err != nil {
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("invalid request: %v", err.Error()))
	}

	if err := s.authService.Authorize(ctx, in.GetLogin(), in.GetPassword()); err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &empty.Empty{}, nil
}

func (s *AuthorizationServiceServer) Register(ctx context.Context, in *pb.AuthorizationRequest) (*empty.Empty, error) {
	if err := s.dataValidator.StructCtx(ctx, in); err != nil {
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("invalid request: %v", err.Error()))
	}

	err := s.userService.Save(ctx, server.User{
		Login:    in.GetLogin(),
		Password: in.GetPassword(),
	})
	if errors.Is(err, server.ErrUserAlreadyExists) {
		return nil, status.Error(codes.Unauthenticated, "invalid login or password")
	}
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error()) // TODO: человеческая ошибка
	}

	return &empty.Empty{}, nil
}
