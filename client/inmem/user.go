package inmem

import (
	"github.com/charmbracelet/log"
	"github.com/mkolibaba/gophkeeper/client"
)

type UserService struct {
	user   *client.User
	logger *log.Logger
}

func NewUserService(logger *log.Logger) *UserService {
	return &UserService{logger: logger}
}

func (s *UserService) Set(user client.User) {
	s.logger.Debug("setting user", "login", user.Login)
	s.user = &user
}

func (s *UserService) Get() *client.User {
	return s.user
}
