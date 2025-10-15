package sqlite

import (
	"context"
	"fmt"
	"github.com/mkolibaba/gophkeeper/server"
	sqlc "github.com/mkolibaba/gophkeeper/server/sqlite/sqlc/gen"
)

type LoginService struct {
	qs *sqlc.Queries
}

func NewLoginService(queries *sqlc.Queries) *LoginService {
	return &LoginService{
		qs: queries,
	}
}

func (s *LoginService) Save(ctx context.Context, data server.LoginData, user string) error {
	err := s.qs.SaveLogin(ctx, sqlc.SaveLoginParams{
		Name:     data.Name,
		Login:    data.Login,
		Password: stringOrNull(data.Password),
		Website:  stringOrNull(data.Website),
		Notes:    stringOrNull(data.Notes),
		User:     user,
	})

	return tryUnwrapSaveError(err)
}

func (s *LoginService) GetAll(ctx context.Context, user string) ([]server.LoginData, error) {
	logins, err := s.qs.GetAllLogins(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("get all: %w", err)
	}

	var result []server.LoginData
	for _, login := range logins {
		result = append(result, server.LoginData{
			Name:     login.Name,
			Login:    login.Login,
			Password: stringOrEmpty(login.Password),
			Website:  stringOrEmpty(login.Website),
			Notes:    stringOrEmpty(login.Notes),
		})
	}

	return result, nil
}

func (s *LoginService) Update(ctx context.Context, data server.LoginDataUpdate, user string) error {
	// TODO: implement
	panic("implement me")
}

func (s *LoginService) Remove(ctx context.Context, name string, user string) error {
	n, err := s.qs.RemoveLogin(ctx, name)
	if err != nil {
		return fmt.Errorf("remove: %w", err)
	}
	if n == 0 {
		return server.ErrDataNotFound
	}
	return nil
}
