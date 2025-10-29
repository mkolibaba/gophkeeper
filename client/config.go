package client

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	GRPC struct {
		ServerAddress string `mapstructure:"server_address"`
	}
	Log struct {
		Output   string
		Truncate bool
	}
	Development struct {
		Enabled    bool
		SpewOutput string `mapstructure:"spew_output"`
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
