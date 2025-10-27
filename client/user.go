package client

import "context"

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
