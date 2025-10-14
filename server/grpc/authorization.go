package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mkolibaba/gophkeeper/proto/gen/go/gophkeeper"
	"github.com/mkolibaba/gophkeeper/server"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type AuthorizationServiceServer struct {
	gophkeeperv1.UnimplementedAuthorizationServiceServer
	userService   server.UserService
	authService   *server.AuthService
	dataValidator *validator.Validate

	jwtSecret string
	jwtTTL    time.Duration
}

func NewAuthorizationServiceServer(
	userService server.UserService,
	authService *server.AuthService,
	dataValidator *validator.Validate,
	config *server.Config,
) *AuthorizationServiceServer {
	// TODO: стоит разделить
	rules := map[string]string{
		"login":    "required",
		"password": "required",
	}

	dataValidator.RegisterStructValidationMapRules(rules, gophkeeperv1.AuthorizationRequest{})

	return &AuthorizationServiceServer{
		userService:   userService,
		authService:   authService,
		dataValidator: dataValidator,
		jwtSecret:     config.JWT.Secret,
		jwtTTL:        config.JWT.TTL,
	}
}

// TODO: можно убрать токен и просто отправлять статус ок/не ок
func (s *AuthorizationServiceServer) Authorize(ctx context.Context, in *gophkeeperv1.AuthorizationRequest) (*gophkeeperv1.AuthorizationResponse, error) {
	if err := s.dataValidator.StructCtx(ctx, in); err != nil {
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("invalid request: %v", err.Error()))
	}

	if err := s.authService.Authorize(ctx, in.GetLogin(), in.GetPassword()); err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	token, err := s.newJWT(in.GetLogin())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	var out gophkeeperv1.AuthorizationResponse
	out.SetToken(token)
	return &out, nil
}

func (s *AuthorizationServiceServer) Register(ctx context.Context, in *gophkeeperv1.AuthorizationRequest) (*empty.Empty, error) {
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

func (s *AuthorizationServiceServer) newJWT(login string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   login,
		Issuer:    "gophkeeper",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtTTL)),
	}).SignedString([]byte(s.jwtSecret))
}
