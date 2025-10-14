package grpc

import "os"

const (
	defaultPort = "8080"
)

type Config struct {
	Port string
}

func NewConfig() *Config {
	var cfg Config

	cfg.Port = os.Getenv("GRPC_PORT")
	if cfg.Port == "" {
		cfg.Port = defaultPort
	}

	return &cfg
}
