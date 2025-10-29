package sqlite

import (
	"github.com/charmbracelet/log"
	"github.com/mkolibaba/gophkeeper/server"
	sqlc "github.com/mkolibaba/gophkeeper/server/sqlite/sqlc/gen"
	"io"
	"os"
	"testing"
)

var (
	db      *DB
	queries *sqlc.Queries
)

func TestMain(m *testing.M) {
	// Инициализируем db
	tmpDir, err := os.MkdirTemp("", "sqlite-test-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	var config server.Config
	config.SQLite.DataFolder = tmpDir
	config.SQLite.DSN = ":memory:"

	db = NewDB(&config, log.New(io.Discard))

	// Запускаем миграции
	if err := OpenDB(db); err != nil {
		panic(err)
	}

	// Инициализируем queries
	queries = NewQueries(db)

	// Запускаем тест
	os.Exit(m.Run())
}

func initializeTestdata(t *testing.T) {
	// TODO: функция должна вызываться только 1 раз

	t.Log("initializing test data")

	mustCreateUser(t, "alice", "123")
	mustCreateUser(t, "bob", "123")
}
