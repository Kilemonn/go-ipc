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

func TestAccept_WithNoClientAndTimeout(t *testing.T) {
	ipcName := "TestAccept_WithNoClient"
	server, err := NewIPCServer(ipcName, nil)
	require.NoError(t, err)
	defer server.Close()

	_, err = server.Accept(time.Millisecond * 500)
	require.Error(t, err)
}

// Ensure that when a server already exists and we attempt to create a new server using the same channel name
// that it fails, but once [IPCServerConfig.Override]=true is provided it will succeed.
func TestNewServer_AlreadyExists(t *testing.T) {
	ipcChannel := "TestNewServer_AlreadyExists"
	server, err := NewIPCServer(ipcChannel, nil)
	require.NoError(t, err)
	defer server.Close()

	_, err = NewIPCServer(ipcChannel, nil)
	require.Error(t, err)

	server2, err := NewIPCServer(ipcChannel, &IPCServerConfig{Override: true})
	require.NoError(t, err)
	server2.Close()
}

// Ensure that when an IPCServer is created with
// [IPCServerConfig.Override]=true and the socket file
// doesn't exist that we get no error when attempting to remove a
// non-existent file.
func TestNewServer_WithOverrideButDoesntExist(t *testing.T) {
	ipcChannel := "TestNewServer_WithOverrideButDoesntExist"
	server, err := NewIPCServer(ipcChannel, nil)
	require.NoError(t, err)
	server.Close()

	server, err = NewIPCServer(ipcChannel, &IPCServerConfig{Override: true})
	require.NoError(t, err)
	server.Close()
}
