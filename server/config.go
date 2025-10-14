package server

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	v *viper.Viper
}

func NewConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("new config: %w", err)
	}

	return &Config{
		v: v,
	}, nil
}

func (c *Config) Viper() *viper.Viper {
	return c.v
}

func (c *Config) GetGRPCPort() string {
	return c.v.GetString("grpc.port")
}

func (c *Config) GetJWTSecret() string {
	return c.v.GetString("jwt.secret")
}

func (c *Config) GetJWTTTL() time.Duration {
	return c.v.GetDuration("jwt.ttl")
}
