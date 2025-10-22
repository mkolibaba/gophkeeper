package server

import (
	"context"
)

type User struct {
	Login    string
	Password string
}

type UserService interface {
	Get(ctx context.Context, login string) (*User, error)
	Save(ctx context.Context, user User) error
}

func IsCurrentUser(ctx context.Context, user string) bool {
	return UserFromContext(ctx) == user
}
