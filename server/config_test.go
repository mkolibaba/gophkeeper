package server

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	t.Setenv("SQLITE_DATA_FOLDER", "some_path")
	t.Setenv("JWT_TTL", "20m")

	config, err := NewConfig()
	require.NoError(t, err)
	require.Equal(t, "some_path", config.SQLite.DataFolder)
	require.Equal(t, "8080", config.GRPC.Port)
	require.Equal(t, 20*time.Minute, config.JWT.TTL)
}

func TestInvalidConfig(t *testing.T) {
	t.Setenv("JWT_TTL", "heh")
	_, err := NewConfig()
	require.Error(t, err)
}
