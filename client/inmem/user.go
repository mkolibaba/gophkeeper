package inmem

import "github.com/charmbracelet/log"

type UserService struct {
	login string
	token string

	logger *log.Logger
}

func NewUserService(logger *log.Logger) *UserService {
	return &UserService{logger: logger}
}

func (s *UserService) SetInfo(login, token string) {
	s.login = login
	s.token = token
}

func (s *UserService) GetUserLogin() string {
	return s.login
}

func (s *UserService) GetBearerToken() string {
	return s.token
}
