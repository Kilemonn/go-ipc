package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Ensure when a server isn't already listening on the requested channel that we can successfully
// create the server.
func TestNewServer_DoesNotExist_NoOverride(t *testing.T) {
	ipcChannelName := "TestNewServer_DoesNotExist"
	svr, err := NewIPCServer(ipcChannelName, nil)
	require.NoError(t, err)

	err = svr.Close()
	require.NoError(t, err)
}

// Ensure that when we attempt to accept a new connection and none are incoming,
// that we timeout and receive an error.
func TestAccept_WithNoClientAndTimeout(t *testing.T) {
	ipcChannelName := "TestAccept_WithNoClient"
	server, err := NewIPCServer(ipcChannelName, nil)
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

func TestNewServer_EmptyStringAsChannelName(t *testing.T) {
	ipcChannel := "  	"
	_, err := NewIPCServer(ipcChannel, nil)
	require.Error(t, err)
}
