package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultIPCServerConfig()
	require.NotNil(t, config)
	require.False(t, config.Override)
}
