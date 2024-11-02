package golangipc

import (
	"testing"
	"time"

	"github.com/Kilemonn/golang-ipc/client"
	"github.com/Kilemonn/golang-ipc/server"
	"github.com/stretchr/testify/require"
)

func TestSendAndReceive(t *testing.T) {
	ipcChannelName := "TestSendAndReceive"
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
