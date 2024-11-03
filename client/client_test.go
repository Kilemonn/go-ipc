package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Ensure that if we attempt to connect to a channel that does not exist
// that we get an error
func TestNewClient_NoIPCChannelExists(t *testing.T) {
	ipcChannelName := "TestNewClient_WithIPCChannel"
	_, err := NewIPCClient(ipcChannelName)
	require.Error(t, err)
}
