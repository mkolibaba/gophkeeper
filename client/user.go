package client

import "context"

//go:generate moq -stub -pkg mock -out mock/authorization.go . AuthorizationService
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
