package sqlite

import (
	"context"
	"fmt"
	"github.com/mkolibaba/gophkeeper/server"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func unwrapInsertError(err error) error {
	if se, ok := asType[*sqlite.Error](err); ok {
		switch se.Code() {
		case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return server.ErrDataAlreadyExists
		case sqlite3.SQLITE_CONSTRAINT_FOREIGNKEY:
			return server.ErrUserNotFound
		}
	}
	return err
}

type anyData struct {
	user string
}

func (a anyData) GetUser() string {
	return a.user
}

func removeData(
	ctx context.Context,
	getUser func(ctx context.Context, id int64) (string, error),
	remove func(ctx context.Context, id int64) (int64, error),
	id int64,
) error {
	user, err := getUser(ctx, id)
	if err != nil {
		return fmt.Errorf("remove: %w", err)
	}

	if err := server.VerifyCanEditData(ctx, anyData{user: user}); err != nil {
		return server.ErrPermissionDenied
	}

	n, err := remove(ctx, id)
	if err != nil {
		return fmt.Errorf("remove: %w", err)
	}
	if n == 0 {
		return server.ErrDataNotFound
	}
	return nil
}
