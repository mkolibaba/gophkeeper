package server

import (
	"context"
)

// AuthorizationService представляет сервис авторизации.
type AuthorizationService interface {
	// Authorize авторизует пользователя по его login, возвращая токен.
	Authorize(ctx context.Context, login string) (string, error)
}
