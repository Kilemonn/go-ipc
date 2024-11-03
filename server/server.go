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
	listener       net.Listener
}

func NewIPCServer(ipcChannelName string, config *server_config.IPCServerConfig) (*IPCServer, error) {
	if config == nil {
		config = &server_config.IPCServerConfig{}
	}

	sever := &IPCServer{
		IpcChannelName: ipcChannelName,
		config:         *config,
	}
	return sever, sever.initialiseServer()
}

func (s *IPCServer) initialiseServer() error {
	descriptor := consts.UNIX_PATH_PREFIX + s.IpcChannelName + consts.UNIX_SOCKET_SUFFIX
	_, err := os.Stat(descriptor)
	if err == nil {
		// If the file exists AND override is enabled we will remove the file
		if s.config.Override {
			err = os.Remove(descriptor)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	s.listener, err = net.Listen("unix", descriptor)
	return err
}

func (s *IPCServer) Close() error {
	return s.listener.Close()
}

func (s *IPCServer) Accept(timeOut time.Duration) (client.IPCClient, error) {
	// Accept timeout
	conn, err := s.listener.Accept()
	if err != nil {
		return client.IPCClient{}, err
	}
	return client.NewIPCClientFromConnection(conn), nil
}
