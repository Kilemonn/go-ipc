package goipc

import (
	"io"
	"testing"
	"time"

	"github.com/Kilemonn/go-ipc/client"
	"github.com/Kilemonn/go-ipc/server"
	"github.com/stretchr/testify/require"
)

// End to end connection and data read and write.
func TestReadAndWrite(t *testing.T) {
	ipcChannelName := "TestReadAndWrite"
	svr, err := server.NewIPCServer(ipcChannelName, nil)
	require.NoError(t, err)
	defer svr.Close()

	client, err := client.NewIPCClient(ipcChannelName)
	require.NoError(t, err)
	defer client.Close()

	accepted, err := svr.Accept(time.Millisecond * 1000)
	require.NoError(t, err)

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
	ipcName := "TestRead_WithoutBeingAccepted"
	server, err := server.NewIPCServer(ipcName, nil)
	require.NoError(t, err)
	defer server.Close()

	client, err := client.NewIPCClient(ipcName)
	require.NoError(t, err)
	defer client.Close()

	b := make([]byte, 10)
	n, err := client.Read(b)
	require.Equal(t, err, io.EOF)
	require.Equal(t, 0, n)
}

// Ensure that [io.Reader.Read] attempts will timeout and return [io.EOF] if there no data to read.
func TestRead_AcceptedWithNoData(t *testing.T) {
	ipcName := "TestRead_AcceptedWithNoData"
	server, err := server.NewIPCServer(ipcName, nil)
	require.NoError(t, err)
	defer server.Close()

	client, err := client.NewIPCClient(ipcName)
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
	ipcName := "TestWrite_BeforeAccept"
	server, err := server.NewIPCServer(ipcName, nil)
	require.NoError(t, err)
	defer server.Close()

	client, err := client.NewIPCClient(ipcName)
	require.NoError(t, err)
	defer client.Close()

	// Writing before is successful
	n, err := client.Write([]byte(ipcName))
	require.NoError(t, err)
	require.Equal(t, len(ipcName), n)

	accepted, err := server.Accept(time.Millisecond * 100)
	require.NoError(t, err)
	defer accepted.Close()

	b := make([]byte, len(ipcName))
	n, err = accepted.Read(b)
	require.NoError(t, err)
	require.Equal(t, len(ipcName), n)

	require.Equal(t, ipcName, string(b))
}
