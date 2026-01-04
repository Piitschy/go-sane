package scanservice

import (
	"image"

	"github.com/tjgq/sane"
)

type Scanner interface {
	ListDevices() ([]*sane.Device, error)
	ReadImage(dn string) (image.Image, error)
	ScanImage(dn string, fileName string) (string, error)
	ScanAndReturnURL(dn string) (string, error)
	ValidateDeviceName(dn string) error
}

type ScanOptions struct {
	Resolution int
	Mode       string
}
