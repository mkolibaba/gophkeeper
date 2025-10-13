package inmem

import (
	"github.com/mkolibaba/gophkeeper/internal/client"
	"go.uber.org/zap"
)

type UserService struct {
	user   *client.User
	logger *zap.Logger
}

func NewUserService(logger *zap.Logger) *UserService {
	return &UserService{logger: logger}
}

func (s *UserService) Set(user client.User) {
	s.logger.Debug("setting user", zap.String("login", user.Login))
	s.user = &user
}

func (s *UserService) Get() *client.User {
	return s.user
}
