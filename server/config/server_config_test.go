package server_config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	config := IPCServerConfig{}
	require.False(t, config.Override)
}
