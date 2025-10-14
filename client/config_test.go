package client

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Setenv("GRPC_SERVER_ADDRESS", "some_address")
	t.Setenv("LOG_TRUNCATE", "false")

	config, err := NewConfig()
	require.NoError(t, err)
	require.Equal(t, "some_address", config.GRPC.ServerAddress)
	require.Equal(t, "bin/client.log", config.Log.Output)
	require.False(t, config.Log.Truncate)
}
