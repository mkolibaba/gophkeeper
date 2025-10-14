package sqlite

import (
	"github.com/mkolibaba/gophkeeper/server"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// TODO: rename
func tryUnwrapSaveError(err error) error {
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

// TODO: может быть лучше в дате сразу это обозначить? почитать на сайте sqlc как лучше поступать с нуллабл
func stringOrNull(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
