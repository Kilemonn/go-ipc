package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient_NoIPCChannelExists(t *testing.T) {
	ipcName := "TestNewClient_WithIPCChannel"
	_, err := NewIPCClient(ipcName)
	require.Error(t, err)
}
