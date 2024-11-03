package server

import (
	"net"
	"os"
	"time"

	"github.com/Kilemonn/go-ipc/client"
	"github.com/Kilemonn/go-ipc/consts"
	server_config "github.com/Kilemonn/go-ipc/server/config"
)

type IPCServer struct {
	IpcChannelName string
	config         server_config.IPCServerConfig
	listener       *net.UnixListener
}

func NewIPCServer(ipcChannelName string, config *server_config.IPCServerConfig) (*IPCServer, error) {
	if config == nil {
		config = server_config.DefaultIPCServerConfig()
	}

	server := &IPCServer{
		IpcChannelName: ipcChannelName,
		config:         *config,
	}
	return server, server.initialiseServer()
}

func (s *IPCServer) initialiseServer() (err error) {
	descriptor := consts.UNIX_PATH_PREFIX + s.IpcChannelName + consts.UNIX_SOCKET_SUFFIX
	addr, err := net.ResolveUnixAddr("unix", descriptor)
	if err != nil {
		return err
	}
	// If override is enabled we will just remove the descriptor before attempting to listen on it
	if s.config.Override {
		err = os.Remove(descriptor)
		if err != nil {
			return err
		}
	}
	s.listener, err = net.ListenUnix("unix", addr)
	return err
}

func (s *IPCServer) Close() error {
	return s.listener.Close()
}

func (s *IPCServer) Accept(timeOut time.Duration) (client.IPCClient, error) {
	s.listener.SetDeadline(time.Now().Add(timeOut))
	conn, err := s.listener.Accept()
	if err != nil {
		return client.IPCClient{}, err
	}
	return client.NewIPCClientFromConnection(conn), nil
}
