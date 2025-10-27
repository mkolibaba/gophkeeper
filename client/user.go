package client

import "context"

//go:generate moq -stub -pkg mock -out mock/authorization.go . AuthorizationService
//go:generate moq -stub -pkg mock -out mock/user.go . UserService
type (
	AuthorizationService interface {
		Authorize(ctx context.Context, login string, password string) (string, error)
		Register(ctx context.Context, login string, password string) (string, error)
	}

	UserService interface {
		SetInfo(login string, token string)
		GetUserLogin() string
		GetBearerToken() string
	}
)
