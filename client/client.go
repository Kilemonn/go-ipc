package client

import (
	"io"
	"net"
	"time"

	"github.com/Kilemonn/go-ipc/consts"
)

const (
	defaultReadTimeout = 10
)

// IPCClient which implements the [io.ReadWriteCloser] interface. Effectively wrapping the [net.Conn] object.
type IPCClient struct {
	Conn net.Conn
	// Time in microseconds to wait during read before timing out. Default is 10.
	readTimeout int
}

func NewIPCClient(ipcChannelName string) (IPCClient, error) {
	descriptor := consts.UNIX_PATH_PREFIX + ipcChannelName + consts.UNIX_SOCKET_SUFFIX
	conn, err := net.Dial("unix", descriptor)
	if err != nil {
		return IPCClient{}, nil
	}
	return NewIPCClientFromConnection(conn), nil
}

func NewIPCClientFromConnection(conn net.Conn) IPCClient {
	return IPCClient{
		Conn:        conn,
		readTimeout: defaultReadTimeout,
	}
}

func (c *IPCClient) SetReadTimeout(timeout int) {
	if timeout < 0 {
		timeout = defaultReadTimeout
	}
	c.readTimeout = timeout
}

func (c *IPCClient) GetReadTimeout() int {
	return c.readTimeout
}

func (c *IPCClient) Close() error {
	return c.Conn.Close()
}

// If [IPCClient.SetReadTimeout] is greater than 0, than a read deadline will be set. Allowing this call to block only for the provided
// period of time before returning 0, and [io.EOF] if there is no result read in that time.
func (c *IPCClient) Read(b []byte) (n int, err error) {
	if c.readTimeout > 0 {
		c.Conn.SetDeadline(time.Now().Add(time.Duration(c.readTimeout) * time.Microsecond))
	}

	n, err = c.Conn.Read(b)
	if err != nil {
		if e, ok := err.(net.Error); ok && e.Timeout() {
			// During timeout, we should return EOF
			return 0, io.EOF
		} else {
			// On other errors, make sure we return the error back to the caller
			return
		}
	}
	return
}

func (c *IPCClient) Write(p []byte) (n int, err error) {
	return c.Conn.Write(p)
}
