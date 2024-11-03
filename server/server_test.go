package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewServer_DoesNotExist_NoOverride(t *testing.T) {
	ipcName := "TestNewServer_DoesNotExist_NoOverride"
	svr, err := NewIPCServer(ipcName, nil)
	require.NoError(t, err)

	svr.Close()
	require.Error(t, err)
}
