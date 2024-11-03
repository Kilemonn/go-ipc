package client

import (
	"io"
	"net"
	"time"

	"github.com/Kilemonn/go-ipc/consts"
)

// IPCClient which implements the [io.ReadWriteCloser] interface. Effectively wrapping the [net.Conn] object.
type IPCClient struct {
	Conn net.Conn
	// A timeout applied at each read attempt, you can set this to 0 to wait indefinitely to make the call blocking
	ReadTimeout time.Duration
}

// NewIPCClient creates a new [IPCClient] connected to the provided channel name.
func NewIPCClient(ipcChannelName string) (IPCClient, error) {
	descriptor := consts.ChannelPathPrefix + ipcChannelName + consts.ChannelSocketSuffix
	conn, err := net.Dial("unix", descriptor)
	if err != nil {
		return IPCClient{}, err
	}
	return NewIPCClientFromConnection(conn), nil
}

// NewIPCClientFromConnection wraps a [net.Conn], used by the [IPCServer] when accepting new connections.
func NewIPCClientFromConnection(conn net.Conn) IPCClient {
	return IPCClient{
		Conn:        conn,
		ReadTimeout: consts.DefaultClientReadTimeout,
	}
}

// Close wraps [net.Conn.Close].
func (c IPCClient) Close() error {
	return c.Conn.Close()
}

// If [IPCClient.SetReadTimeout] is greater than 0, than a read deadline will be set. Allowing this call to block only for the provided
// period of time before returning 0, and [io.EOF] if there is no result read in that time.
func (c IPCClient) Read(b []byte) (n int, err error) {
	c.Conn.SetDeadline(time.Now().Add(c.ReadTimeout))
	n, err = c.Conn.Read(b)
	if err != nil {
		if e, ok := err.(net.Error); ok && e.Timeout() {
			// During timeout, we should return EOF
			return 0, io.EOF
		}
	}
	return
}

// Write wraps [net.Conn.Write].
func (c IPCClient) Write(p []byte) (n int, err error) {
	return c.Conn.Write(p)
}
