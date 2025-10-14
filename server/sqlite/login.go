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

func (l *LoginService) Save(ctx context.Context, data server.LoginData, user string) error {
	err := l.qs.SaveLogin(ctx, sqlc.SaveLoginParams{
		Name:     data.Name,
		Login:    data.Login,
		Password: stringOrNull(data.Password),
		Website:  stringOrNull(data.Website),
		Notes:    stringOrNull(data.Notes),
		User:     user,
	})

	return tryUnwrapSaveError(err)
}

func (l *LoginService) GetAll(ctx context.Context, user string) ([]server.LoginData, error) {
	logins, err := l.qs.GetAllLogins(ctx, user)
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

func (l *LoginService) Remove(ctx context.Context, name string, user string) error {
	n, err := l.qs.RemoveLogin(ctx, name)
	if err != nil {
		return fmt.Errorf("remove: %w", err)
	}
	if n == 0 {
		return server.ErrDataNotFound
	}
	return nil
}
