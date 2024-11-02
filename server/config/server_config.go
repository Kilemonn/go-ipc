package server_config

type IPCServerConfig struct {
	// If the requested IPC name is already created then this flag will remove and re-create it, otherwise fail if it already exists
	Override bool
}
