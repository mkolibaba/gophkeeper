package sqlite

import (
	"encoding/json"
	"fmt"
	"github.com/mkolibaba/gophkeeper/internal/common/errors"
	"github.com/mkolibaba/gophkeeper/internal/server"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// TODO: rename
func tryUnwrapSaveError(err error) error {
	if se, ok := errors.AsType[*sqlite.Error](err); ok {
		switch se.Code() {
		case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return server.ErrDataAlreadyExists
		case sqlite3.SQLITE_CONSTRAINT_FOREIGNKEY:
			return server.ErrUserNotFound
		}
	}
	return err
}

func unmarshalMetadata(data []byte) (metadata map[string]string, err error) {
	if err = json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("unmarshal metadata: %w", err)
	}
	return
}
