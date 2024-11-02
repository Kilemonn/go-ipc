package client

import (
	"io"
	"net"
	"time"

	"github.com/Kilemonn/golang-ipc/consts"
)

// IPCClient which implements the [io.ReadWriteCloser] interface.
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
		readTimeout: 10,
	}
}

func (c *IPCClient) SetReadTimeout(timeout int) {
	if timeout < 0 {
		timeout = 10
	}
	c.readTimeout = timeout
}

func (c *IPCClient) GetReadTimeout() int {
	return c.readTimeout
}

func (c *IPCClient) Close() error {
	return c.Conn.Close()
}

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
