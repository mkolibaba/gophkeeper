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

func (s *LoginService) Create(ctx context.Context, data server.LoginData) error {
	err := s.qs.InsertLogin(ctx, sqlc.InsertLoginParams{
		Name:     data.Name,
		Login:    data.Login,
		Password: stringOrNull(data.Password),
		Website:  stringOrNull(data.Website),
		Notes:    stringOrNull(data.Notes),
		User:     server.UserFromContext(ctx),
	})
	return unwrapInsertError(err)
}

func (s *LoginService) GetAll(ctx context.Context) ([]server.LoginData, error) {
	user := server.UserFromContext(ctx)

	logins, err := s.qs.SelectLogins(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("get all: %w", err)
	}

	var result []server.LoginData
	for _, login := range logins {
		result = append(result, server.LoginData{
			ID:       login.ID,
			Name:     login.Name,
			Login:    login.Login,
			Password: stringOrEmpty(login.Password),
			Website:  stringOrEmpty(login.Website),
			Notes:    stringOrEmpty(login.Notes),
			User:     user,
		})
	}

	return result, nil
}

func (s *LoginService) Update(ctx context.Context, id int64, data server.LoginDataUpdate) error {
	login, err := s.qs.SelectLogin(ctx, id)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	if err := server.VerifyCanEditData(ctx, login); err != nil {
		return err
	}

	params := sqlc.UpdateLoginParams{
		ID:       login.ID,
		Name:     login.Name,
		Login:    login.Login,
		Password: login.Password,
		Website:  login.Website,
		Notes:    login.Notes,
	}

	if data.Name != nil {
		params.Name = *data.Name
	}
	if data.Login != nil {
		params.Login = *data.Login
	}
	if data.Password != nil {
		params.Password = data.Password
	}
	if data.Website != nil {
		params.Website = data.Website
	}
	if data.Notes != nil {
		params.Notes = data.Notes
	}

	n, err := s.qs.UpdateLogin(ctx, params)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("update: no rows")
	}
	return nil
}

func (s *LoginService) Remove(ctx context.Context, id int64) error {
	return removeData(ctx, s.qs.SelectLoginUser, s.qs.DeleteLogin, id)
}
