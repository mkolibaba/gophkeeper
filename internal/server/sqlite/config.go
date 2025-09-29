package sqlite

import "os"

const (
	defaultDSN = "gophkeeper.sqlite"
)

type Config struct {
	DSN string
}

func NewConfig() *Config {
	var cfg Config

	cfg.DSN = os.Getenv("SQLITE_DSN")
	if cfg.DSN == "" {
		cfg.DSN = defaultDSN
	}

	return &cfg
}
