package scanservice

import "errors"

var (
	ErrIsScanning     error = errors.New("scanner is busy")
	ErrDeviceNotFound error = errors.New("device not found")
)
