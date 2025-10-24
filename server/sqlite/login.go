package sqlite

import (
	"context"
	"database/sql"
	"errors"
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
	_, err := s.qs.InsertLogin(ctx, s.converter.ConvertToInsertLogin(ctx, data))
	return unwrapInsertError(err)
}

func (s *LoginService) GetAll(ctx context.Context) ([]server.LoginData, error) {
	return getAllData(ctx, s.qs.SelectLogins, s.converter.ConvertToLoginDataSlice)
}

func (s *LoginService) Update(ctx context.Context, id int64, data server.LoginDataUpdate) error {
	login, err := s.qs.SelectLogin(ctx, id, server.UserFromContext(ctx))
	if errors.Is(err, sql.ErrNoRows) {
		return server.ErrDataNotFound
	}
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	params := s.converter.ConvertToUpdateLogin(login)
	s.converter.ConvertToUpdateLoginUpdate(data, &params)

	n, err := s.qs.UpdateLogin(ctx, params)
	if n == 0 {
		return server.ErrDataNotFound
	}
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}

func (s *LoginService) Remove(ctx context.Context, id int64) error {
	return removeDataV2(ctx, s.qs.DeleteLogin, id)
}
