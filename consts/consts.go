package consts

import "time"

const (
	UNIX_PATH_PREFIX   = "/tmp/"
	UNIX_SOCKET_SUFFIX = ".sock"

	DefaultClientReadTimeout = time.Millisecond * 10
)
