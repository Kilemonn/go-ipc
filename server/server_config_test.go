package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Ensure default values for the [IPCServerConfig].
func TestDefaultConfig(t *testing.T) {
	config := DefaultIPCServerConfig()
	require.NotNil(t, config)
	require.False(t, config.Override)
}
