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
	server, err := server.NewIPCServer(ipcChannelName, nil)
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
	server, err := server.NewIPCServer(ipcChannelName, nil)
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
	server, err := server.NewIPCServer(ipcChannelName, nil)
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
	server, err := server.NewIPCServer(channel, nil)
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
