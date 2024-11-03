package server

import (
	"os"
	"testing"

	"github.com/Kilemonn/go-ipc/consts"
	"github.com/stretchr/testify/require"
)

func TestNewServer_DoesNotExist_NoOverride(t *testing.T) {
	ipcName := "TestNewServer_DoesNotExist_NoOverride"
	_, err := os.Stat(consts.UNIX_PATH_PREFIX + ipcName + consts.UNIX_SOCKET_SUFFIX)
	require.Error(t, err, "ipc socket descriptor should not exist")

	svr, err := NewIPCServer(ipcName, nil)
	require.NoError(t, err)

	_, err = os.Stat(consts.UNIX_PATH_PREFIX + ipcName + consts.UNIX_SOCKET_SUFFIX)
	require.NoError(t, err)

	svr.Close()
	_, err = os.Stat(consts.UNIX_PATH_PREFIX + ipcName + consts.UNIX_SOCKET_SUFFIX)
	require.Error(t, err)
}
