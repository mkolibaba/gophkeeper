package jwt

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mkolibaba/gophkeeper/server"
	"time"
)

type AuthorizationService struct {
	secret string
	ttl    time.Duration
}

func NewAuthorizationService(config *server.Config) *AuthorizationService {
	return &AuthorizationService{
		secret: config.JWT.Secret,
		ttl:    config.JWT.TTL,
	}
}

func (s *AuthorizationService) Authorize(ctx context.Context, login string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   login,
		Issuer:    "gophkeeper",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.ttl)),
	}).SignedString([]byte(s.secret))
}
