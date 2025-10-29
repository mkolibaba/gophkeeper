package sqlite

import (
	"context"
	"database/sql"
	stderrors "errors"
	"fmt"
	"github.com/mkolibaba/gophkeeper/server"
	sqlc "github.com/mkolibaba/gophkeeper/server/sqlite/sqlc/gen"
	"golang.org/x/crypto/bcrypt"
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

func (s *UserService) Get(ctx context.Context, login string) (*server.User, error) {
	u, err := s.qs.SelectUser(ctx, login)
	if stderrors.Is(err, sql.ErrNoRows) {
		return nil, server.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return &server.User{
		Login:    u.Login,
		Password: u.Password,
	}, nil
}

func (s *UserService) Save(ctx context.Context, user server.User) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	err = s.qs.InsertUser(ctx, user.Login, string(passwordHash))

	if se, ok := asType[*sqlite.Error](err); ok && se.Code() == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY {
		return server.ErrUserAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}

	return nil
}
