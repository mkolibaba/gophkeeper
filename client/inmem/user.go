package inmem

type UserService struct {
	login *string // TODO: доставать логин из токена
	token *string
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) SetInfo(login, token string) {
	s.login, s.token = &login, &token
}

func (s *UserService) GetUserLogin() *string {
	return s.login
}

func (s *UserService) GetBearerToken() *string {
	return s.token
}
