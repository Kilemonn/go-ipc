package server

import (
	"io"
	"testing"
	"time"

	"github.com/Kilemonn/go-ipc/client"
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

// Ensure we return an error if the channel name contains only whitespace.
func TestNewServer_EmptyStringAsChannelName(t *testing.T) {
	ipcChannel := "  	"
	_, err := NewIPCServer(ipcChannel, nil)
	require.Error(t, err)
}

// End to end connection and data read and write.
func TestReadAndWrite(t *testing.T) {
	ipcChannelName := "TestReadAndWrite"
	svr, err := NewIPCServer(ipcChannelName, nil)
	require.NoError(t, err)
	defer svr.Close()

	client, err := client.NewIPCClient(ipcChannelName)
	require.NoError(t, err)
	defer client.Close()

	accepted, err := svr.Accept(time.Millisecond * 1000)
	require.NoError(t, err)
	defer accepted.Close()

	content := "some-data"
	n, err := client.Write([]byte(content))
	require.NoError(t, err)
	require.Equal(t, len(content), n)

	b := make([]byte, len(content))
	n, err = accepted.Read(b)
	require.NoError(t, err)
	require.Equal(t, len(content), n)

	require.Equal(t, content, string(b))
}

// Understand behaviour of read even when the connection has not been accepted by the server.
func TestRead_WithoutBeingAccepted(t *testing.T) {
	ipcChannelName := "TestRead_WithoutBeingAccepted"
	server, err := NewIPCServer(ipcChannelName, nil)
	require.NoError(t, err)
	defer server.Close()

	client, err := client.NewIPCClient(ipcChannelName)
	require.NoError(t, err)
	defer client.Close()

	b := make([]byte, 10)
	n, err := client.Read(b)
	require.Equal(t, err, io.EOF)
	require.Equal(t, 0, n)
}

// Ensure that [io.Reader.Read] attempts will timeout and return [io.EOF] if there no data to read.
func TestRead_AcceptedWithNoData(t *testing.T) {
	ipcChannelName := "TestRead_AcceptedWithNoData"
	server, err := NewIPCServer(ipcChannelName, nil)
	require.NoError(t, err)
	defer server.Close()

	client, err := client.NewIPCClient(ipcChannelName)
	require.NoError(t, err)
	defer client.Close()

	accepted, err := server.Accept(time.Millisecond * 100)
	require.NoError(t, err)
	defer accepted.Close()

	b := make([]byte, 10)
	n, err := client.Read(b)
	require.Equal(t, err, io.EOF)
	require.Equal(t, 0, n)

	n, err = accepted.Read(b)
	require.Equal(t, err, io.EOF)
	require.Equal(t, 0, n)
}

// Ensure that [io.Writer.Write] will be able to be sent from a client, even when not accepted yet by
// the server. And once the server accepts the data can be read.
func TestWrite_BeforeAccept(t *testing.T) {
	ipcChannelName := "TestWrite_BeforeAccept"
	server, err := NewIPCServer(ipcChannelName, nil)
	require.NoError(t, err)
	defer server.Close()

	client, err := client.NewIPCClient(ipcChannelName)
	require.NoError(t, err)
	defer client.Close()

	// Writing before is successful
	n, err := client.Write([]byte(ipcChannelName))
	require.NoError(t, err)
	require.Equal(t, len(ipcChannelName), n)

	accepted, err := server.Accept(time.Millisecond * 100)
	require.NoError(t, err)
	defer accepted.Close()

	b := make([]byte, len(ipcChannelName))
	n, err = accepted.Read(b)
	require.NoError(t, err)
	require.Equal(t, len(ipcChannelName), n)

	require.Equal(t, ipcChannelName, string(b))
}

// Testing how multiple connections work.
func TestReadWrite_MultipleClients(t *testing.T) {
	channel := "TestReadWrite_MultipleClients"
	server, err := NewIPCServer(channel, nil)
	require.NoError(t, err)
	defer server.Close()

	cli1, err := client.NewIPCClient(channel)
	require.NoError(t, err)
	defer cli1.Close()

	accept1, err := server.Accept(time.Millisecond * 100)
	require.NoError(t, err)
	defer accept1.Close()

	cli2, err := client.NewIPCClient(channel)
	require.NoError(t, err)
	defer cli2.Close()

	accept2, err := server.Accept(time.Millisecond * 100)
	require.NoError(t, err)
	defer accept2.Close()

	// Sending from cli1, making sure that accept2 cannot read it
	data := "TestReadWrite_MultipleClients"
	n, err := cli1.Write([]byte(data))
	require.NoError(t, err)
	require.Equal(t, len(data), n)

	b := make([]byte, len(data))
	n, err = accept2.Read(b)
	require.Error(t, err)
	require.Equal(t, io.EOF, err)
	require.Equal(t, 0, n)

	// Reading into accept1 to check it can be recieved
	n, err = accept1.Read(b)
	require.NoError(t, err)
	require.Equal(t, len(data), n)
	require.Equal(t, data, string(b))
}
