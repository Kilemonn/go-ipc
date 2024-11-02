package errors

import "errors"

var (
	ErrExistingIPCName = errors.New("ipc channel with the same name already exists and override is set to false")
)
