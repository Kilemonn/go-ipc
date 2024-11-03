package consts

import "time"

const (
	ChannelPathPrefix   = "/tmp/"
	ChannelSocketSuffix = ".sock"

	DefaultClientReadTimeout = time.Millisecond * 10
)
