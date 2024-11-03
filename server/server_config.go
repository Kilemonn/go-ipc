package server

type IPCServerConfig struct {
	// If the requested IPC name is already created then this flag will remove and re-create it, otherwise fail if it already exists
	Override bool
}

// DefaultIPCServerConfig the config used when [nil] is provided to a new [IPCServer].
func DefaultIPCServerConfig() *IPCServerConfig {
	return &IPCServerConfig{}
}
