package sqlite

import (
	"context"
	"database/sql"
	stderrors "errors"
	"fmt"
	"github.com/mkolibaba/gophkeeper/internal/common/errors"
	"github.com/mkolibaba/gophkeeper/internal/server"
	sqlc "github.com/mkolibaba/gophkeeper/internal/server/sqlite/sqlc/gen"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type UserService struct {
	qs *sqlc.Queries
}

func NewUserService(queries *sqlc.Queries) *UserService {
	return &UserService{
		qs: queries,
	}
}

func (s *UserService) Get(ctx context.Context, login string) (server.User, error) {
	u, err := s.qs.GetUserForLogin(ctx, login)
	if stderrors.Is(err, sql.ErrNoRows) {
		return server.User{}, server.ErrUserNotFound
	}
	if err != nil {
		return server.User{}, fmt.Errorf("get: %w", err)
	}

	return server.User{
		Login:    u.Login,
		Password: u.Password,
	}, nil
}

func (s *UserService) Save(ctx context.Context, user server.User) error {
	err := s.qs.SaveUser(ctx, user.Login, user.Password)

	if se, ok := errors.AsType[*sqlite.Error](err); ok && se.Code() == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY {
		return server.ErrUserAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	return nil
}
