package server

import (
	"net"
	"time"

	"github.com/Kilemonn/go-ipc/client"
	"github.com/Kilemonn/go-ipc/consts"
	server_config "github.com/Kilemonn/go-ipc/server/config"
)

type IPCServer struct {
	IpcChannelName string
	config         server_config.IPCServerConfig
	listener       net.Listener
}

func NewIPCServer(ipcChannelName string, config *server_config.IPCServerConfig) (*IPCServer, error) {
	if config == nil {
		config = &server_config.IPCServerConfig{}
	}

	server := &IPCServer{
		IpcChannelName: ipcChannelName,
		config:         *config,
	}
	return server, server.initialiseServer()
}

func (s *IPCServer) initialiseServer() (err error) {
	descriptor := consts.UNIX_PATH_PREFIX + s.IpcChannelName + consts.UNIX_SOCKET_SUFFIX
	s.listener, err = net.Listen("unix", descriptor)
	return err
}

func (s *IPCServer) Close() error {
	return s.listener.Close()
}

func (s *IPCServer) Accept(timeOut time.Duration) (client.IPCClient, error) {
	// TODO: Accept timeout
	conn, err := s.listener.Accept()
	if err != nil {
		return client.IPCClient{}, err
	}
	return client.NewIPCClientFromConnection(conn), nil
}
