package inmem

import (
	"github.com/mkolibaba/gophkeeper/internal/client"
)

type UserService struct {
	user *client.User
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Set(user client.User) {
	s.user = &user
}

func (s *UserService) Get() *client.User {
	return s.user
}
