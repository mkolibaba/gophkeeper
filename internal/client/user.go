package client

import "context"

type (
	User struct {
		Login    string
		Password string
	}

	AuthorizationService interface {
		Authorize(ctx context.Context, login string, password string) error
		Register(ctx context.Context, login string, password string) error
	}

	UserService interface {
		Set(User)
		Get() *User
	}
)
