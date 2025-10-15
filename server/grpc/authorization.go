package grpc

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeper"
	"github.com/mkolibaba/gophkeeper/server"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrInvalidCredentials = status.Error(codes.InvalidArgument, "invalid login or password")

type AuthorizationServiceServer struct {
	gophkeeperv1.UnimplementedAuthorizationServiceServer
	userService          server.UserService
	authorizationService server.AuthorizationService
	validate             *validator.Validate
}

func NewAuthorizationServiceServer(
	userService server.UserService,
	authorizationService server.AuthorizationService,
	validate *validator.Validate,
) *AuthorizationServiceServer {
	return &AuthorizationServiceServer{
		userService:          userService,
		authorizationService: authorizationService,
		validate:             validate,
	}
}

func (s *AuthorizationServiceServer) Authorize(
	ctx context.Context,
	in *gophkeeperv1.UserCredentials,
) (*gophkeeperv1.TokenResponse, error) {
	if err := s.validate.StructCtx(ctx, in); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := s.userService.Get(ctx, in.GetLogin())
	if errors.Is(err, server.ErrUserNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// TODO: заменить на хеши
	if user.Password != in.GetPassword() {
		return nil, ErrInvalidCredentials
	}

	token, err := s.authorizationService.Authorize(ctx, in.GetLogin())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var out gophkeeperv1.TokenResponse
	out.SetToken(token)
	return &out, nil
}

func (s *AuthorizationServiceServer) Register(
	ctx context.Context,
	in *gophkeeperv1.UserCredentials,
) (*gophkeeperv1.TokenResponse, error) {
	if err := s.validate.StructCtx(ctx, in); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.userService.Save(ctx, server.User{
		Login:    in.GetLogin(),
		Password: in.GetPassword(),
	})
	if errors.Is(err, server.ErrUserAlreadyExists) {
		// Здесь можно было бы возвращать ошибку "Такой пользователь уже существует",
		// но в рамках безопасности лучше не сообщать какие пользователи есть в системе.
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	token, err := s.authorizationService.Authorize(ctx, in.GetLogin())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var out gophkeeperv1.TokenResponse
	out.SetToken(token)
	return &out, nil
}
