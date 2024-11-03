package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewServer_DoesNotExist_NoOverride(t *testing.T) {
	ipcName := "TestNewServer_DoesNotExist_NoOverride"
	svr, err := NewIPCServer(ipcName, nil)
	require.NoError(t, err)

	err = svr.Close()
	require.NoError(t, err)
}

func TestAccept_WithNoClient(t *testing.T) {
	ipcName := "TestAccept_WithNoClient"
	server, err := NewIPCServer(ipcName, nil)
	require.NoError(t, err)
	defer server.Close()

	_, err = server.Accept(time.Millisecond * 500)
	require.Error(t, err)
}
