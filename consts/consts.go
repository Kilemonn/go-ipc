package consts

import "time"

const (
	ChannelPathPrefix   = "/tmp/"
	ChannelSocketSuffix = ".sock"

	// 10 Millisecond default timeout
	DefaultClientReadTimeout = time.Millisecond * 10
)
