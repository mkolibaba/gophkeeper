package sqlite

import "os"

const (
	dataFolder = "data"
)

type Config struct {
	DataFolder string
}

func NewConfig() *Config {
	var cfg Config

	cfg.DataFolder = os.Getenv("SQLITE_DATA_FOLDER")
	if cfg.DataFolder == "" {
		cfg.DataFolder = dataFolder
	}

	return &cfg
}
