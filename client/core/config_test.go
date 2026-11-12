package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Setenv("GRPC_SERVER_ADDRESS", "some_address")
	t.Setenv("LOG_TRUNCATE", "false")

	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "config.toml"), []byte(`
[grpc]
server_address = "127.0.0.1:8080"

[log]
output = "bin/client.log"
truncate = true
`), 0600))
	t.Chdir(dir)

	config, err := NewConfig()
	require.NoError(t, err)
	require.Equal(t, "some_address", config.GRPC.ServerAddress)
	require.Equal(t, "bin/client.log", config.Log.Output)
	require.False(t, config.Log.Truncate)
}
