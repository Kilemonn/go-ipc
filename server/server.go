package server

import (
	"errors"
	"net"
	"os"
	"strings"
	"time"

	"github.com/Kilemonn/go-ipc/client"
	"github.com/Kilemonn/go-ipc/consts"
)

var (
	// Invalid channel name, most likely when it only contains white space
	ErrInvalidChannelName = errors.New("channel name must contain other characters besides whitespace")
)

type IPCServer struct {
	IpcChannelName string
	config         IPCServerConfig
	listener       *net.UnixListener
}

// NewIPCServer create a new [IPCServer] at the provided channel name.
func NewIPCServer(ipcChannelName string, config *IPCServerConfig) (*IPCServer, error) {
	if len(strings.TrimSpace(ipcChannelName)) == 0 {
		return &IPCServer{}, ErrInvalidChannelName
	}

	if config == nil {
		config = DefaultIPCServerConfig()
	}
	server := &IPCServer{
		IpcChannelName: ipcChannelName,
		config:         *config,
	}
	return server, server.initialiseServer()
}

// initialiseServer will setup and listen to the provided ipc channel.
func (s *IPCServer) initialiseServer() (err error) {
	descriptor := consts.ChannelPathPrefix + s.IpcChannelName + consts.ChannelSocketSuffix
	addr, err := net.ResolveUnixAddr("unix", descriptor)
	if err != nil {
		return err
	}
	// If override is enabled we will just remove the descriptor before attempting to listen on it
	if s.config.Override {
		// Ignoring the error, since if we fail to remove it, its most likely to not exist
		_ = os.Remove(descriptor)
	}
	s.listener, err = net.ListenUnix("unix", addr)
	return err
}

// Close wraps [net.UnixListener.Close].
func (s *IPCServer) Close() error {
	return s.listener.Close()
}

// Accept acccepts the next incoming connection, the provided
// timeout can be set to 0 to make this a blocking call.
func (s *IPCServer) Accept(timeOut time.Duration) (client.IPCClient, error) {
	s.listener.SetDeadline(time.Now().Add(timeOut))
	conn, err := s.listener.Accept()
	if err != nil {
		return client.IPCClient{}, err
	}
	return client.NewIPCClientFromConnection(conn), nil
}
