package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/mkolibaba/gophkeeper/server"
	"io/fs"
	_ "modernc.org/sqlite"
	"os"
	"sort"
)

//go:embed migration/*.sql
var migrationFS embed.FS

type DB struct {
	db             *sql.DB
	dsn            string
	binariesFolder string
	logger         *log.Logger
}

func NewDB(config *server.Config, logger *log.Logger) *DB {
	return &DB{
		dsn:            config.SQLite.DSN,
		binariesFolder: fmt.Sprintf("%s/assets/binary", config.SQLite.DataFolder),
		logger:         logger,
	}
}

func (d *DB) Open() (err error) {
	d.db, err = sql.Open("sqlite", d.dsn)
	if err != nil {
		return err
	}

	err = os.MkdirAll(d.binariesFolder, 0755)
	if err != nil {
		return err
	}

	// Включаем проверки foreign key, которые в SQLite по умолчанию выключены.
	if _, err := d.db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return fmt.Errorf("migrate: foreign keys pragma: %w", err)
	}

	if err := d.migrate(); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}

func (d *DB) migrate() error {
	if _, err := d.db.Exec(`CREATE TABLE IF NOT EXISTS migration (name TEXT PRIMARY KEY);`); err != nil {
		return fmt.Errorf("cannot create migrations table: %w", err)
	}

	names, err := fs.Glob(migrationFS, "migration/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(names)

	for _, name := range names {
		if err := d.migrateFile(name); err != nil {
			return fmt.Errorf("migrate script %q: %w", name, err)
		}
	}
	return nil
}

func (d *DB) migrateFile(name string) error {
	d.logger.Info("running migration script", "name", name)

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var n int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM migration WHERE name = ?`, name).Scan(&n); err != nil {
		return err
	} else if n != 0 {
		d.logger.Info("migration script already run", "name", name)
		return nil
	}

	if buf, err := fs.ReadFile(migrationFS, name); err != nil {
		return err
	} else if _, err := tx.Exec(string(buf)); err != nil {
		return err
	}

	if _, err := tx.Exec(`INSERT INTO migration (name) VALUES (?)`, name); err != nil {
		return err
	}

	d.logger.Info("migration script run done", "name", name)
	return tx.Commit()
}

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
