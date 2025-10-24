package server

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type Config struct {
	GRPC struct {
		Port string
	}
	SQLite struct {
		DataFolder string `mapstructure:"data_folder"`
		DSN        string
	}
	JWT struct {
		Secret string
		TTL    time.Duration
	}
	Development struct {
		Enabled bool
	}
}

func NewConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("new config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("new config: %w", err)
	}

	return &cfg, nil
}
