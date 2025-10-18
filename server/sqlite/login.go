package sqlite

import (
	"context"
	"fmt"
	"github.com/mkolibaba/gophkeeper/server"
	"github.com/mkolibaba/gophkeeper/server/sqlite/converter"
	sqlc "github.com/mkolibaba/gophkeeper/server/sqlite/sqlc/gen"
)

type LoginService struct {
	qs        *sqlc.Queries
	converter converter.DataConverter
}

func NewLoginService(queries *sqlc.Queries, converter converter.DataConverter) *LoginService {
	return &LoginService{
		qs:        queries,
		converter: converter,
	}
}

func (s *LoginService) Create(ctx context.Context, data server.LoginData) error {
	err := s.qs.InsertLogin(ctx, s.converter.ConvertToInsertLogin(ctx, data))
	return unwrapInsertError(err)
}

func (s *LoginService) GetAll(ctx context.Context) ([]server.LoginData, error) {
	return getAllData(ctx, s.qs.SelectLogins, s.converter.ConvertToLoginDataSlice)
}

func (s *LoginService) Update(ctx context.Context, id int64, data server.LoginDataUpdate) error {
	login, err := s.qs.SelectLogin(ctx, id)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	if err := server.VerifyCanEditData(ctx, login); err != nil {
		return err
	}

	params := s.converter.ConvertToUpdateLogin(login)
	s.converter.ConvertToUpdateLoginUpdate(data, &params)

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
