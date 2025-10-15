package server

import (
	"context"
)

type User struct {
	Login    string
	Password string // TODO: hash
}

type UserService interface {
	Get(ctx context.Context, login string) (User, error)
	Save(ctx context.Context, user User) error
}
