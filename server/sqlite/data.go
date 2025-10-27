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
		case sqlite3.SQLITE_CONSTRAINT_FOREIGNKEY:
			return server.ErrUserNotFound
		}
	}
	return err
}

func getAllData[S any, R any](
	ctx context.Context,
	getter func(context.Context, string) ([]S, error),
	mapper func([]S) []R,
) ([]R, error) {
	sources, err := getter(ctx, server.UserFromContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("get all: %w", err)
	}
	return mapper(sources), nil
}

func removeData(
	ctx context.Context,
	remove func(ctx context.Context, id int64, user string) (int64, error),
	id int64,
) error {
	n, err := remove(ctx, id, server.UserFromContext(ctx))
	if n == 0 {
		return server.ErrDataNotFound
	}
	if err != nil {
		return fmt.Errorf("remove: %w", err)
	}
	return nil
}
