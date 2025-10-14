package grpc

import (
	"fmt"
	"os"
)

type Config struct {
	ServerAddress string
}

func NewConfig() (*Config, error) {
	var cfg Config

	cfg.ServerAddress = os.Getenv("GRPC_SERVER_ADDRESS")
	if cfg.ServerAddress == "" {
		return nil, fmt.Errorf("missing environment variable GRPC_SERVER_ADDRESS")
	}

	return &cfg, nil
}
