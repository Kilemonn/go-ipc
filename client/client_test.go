package client

import (
	"os"
	"testing"

	"github.com/Kilemonn/go-ipc/consts"
	"github.com/stretchr/testify/require"
)

func TestNewClient_NoIPCChannelExists(t *testing.T) {
	ipcName := "TestNewClient_WithIPCChannel"
	_, err := os.Stat(consts.UNIX_PATH_PREFIX + ipcName + consts.UNIX_SOCKET_SUFFIX)
	require.Error(t, err, "ipc socket descriptor should not exist")

	_, err = NewIPCClient(ipcName)
	require.NoError(t, err)
}
