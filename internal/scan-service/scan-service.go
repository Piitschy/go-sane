package scanservice

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/tjgq/sane"
)

func ListSaneDevices() ([]*sane.Device, error) {
	devs, err := sane.Devices()
	if err != nil {
		return nil, fmt.Errorf("Error querying devices: %v", err)
	}

	d := make([]*sane.Device, 0, len(devs))
	if len(devs) == 0 {
		return d, fmt.Errorf("no SANE devices found.")
	}
	for _, dev := range devs {
		if strings.Contains(dev.Name, "video") {
			continue
		}
		d = append(d, &dev)
	}
	return d, nil
}

func ListDeviceNames() ([]string, error) {
	devs, err := ListSaneDevices()
	names := make([]string, len(devs))
	if err != nil {
		return names, err
	}
	for i, d := range devs {
		names[i] = d.Name
	}
	return names, nil
}

type ScanService struct {
	devs            []*sane.Device
	lastDeviceCheck time.Time
	TargetPath      string
	mu              sync.Mutex
}

func NewScanService() (*ScanService, error) {
	err := sane.Init()
	if err != nil {
		return nil, fmt.Errorf("could not initialize SANE: %v", err)
	}

	devs, err := ListSaneDevices()
	if err != nil {
		return nil, err
	}

	return &ScanService{
		TargetPath:      "scans",
		devs:            devs,
		lastDeviceCheck: time.Now(),
		mu:              sync.Mutex{},
	}, nil
}

func (s *ScanService) Close() {
	sane.Exit()
}

func (s *ScanService) ListDevices() ([]*sane.Device, error) {
	if time.Now().After(s.lastDeviceCheck.Add(1 * time.Minute)) {
		devs, err := ListSaneDevices()
		if err != nil {
			return nil, err
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		s.devs = devs
		s.lastDeviceCheck = time.Now()
	}
	return s.devs, nil
}

func (s *ScanService) ValidateDeviceName(dn string) error {
	devs, err := s.ListDevices()
	if err != nil {
		return err
	}
	for _, d := range devs {
		if d.Name == dn {
			return nil
		}
	}
	return ErrDeviceNotFound
}

func (s *ScanService) ReadImage(dn string) (image.Image, error) {
	if err := s.ValidateDeviceName(dn); err != nil {
		return nil, err
	}

	if isScaning := !s.mu.TryLock(); isScaning {
		return nil, ErrIsScanning
	}
	defer s.mu.Unlock()

	c, err := sane.Open(dn)
	if err != nil {
		return nil, fmt.Errorf("error while open connection to '%s': %v", dn, err)
	}
	defer c.Close()
	c.SetOption("mode", "Color")
	c.SetOption("resolution", 300)

	// Set scan area to maximum to capture full document
	c.SetOption("tl-x", 0.0)
	c.SetOption("tl-y", 0.0)
	c.SetOption("br-x", 215.0) // Maximum width for A4 scanner
	c.SetOption("br-y", 297.0) // Maximum height for A4 scanner

	img, err := c.ReadImage()
	if err != nil {
		return img, fmt.Errorf("Scan error: %v", err)
	}
	return img, nil
}

func (s *ScanService) ScanImage(dn string, fileName string) (string, error) {
	img, err := s.ReadImage(dn)
	if err != nil {
		return "", err
	}
	if fileName == "" {
		fileName = fmt.Sprintf("scan_%d.jpg", time.Now().Unix())
	}
	filePath := path.Join(s.TargetPath, fileName)

	// Ensure directory exists
	if err := os.MkdirAll(s.TargetPath, 0755); err != nil {
		return "", fmt.Errorf("error creating directory: %v", err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return filePath, fmt.Errorf("error while open target file: %v", err)
	}
	defer f.Close()

	err = jpeg.Encode(f, img, nil)
	if err != nil {
		return filePath, fmt.Errorf("encoding error: %v", err)
	}
	return filePath, nil
}

func (s *ScanService) ScanAndReturnURL(dn string) (string, error) {
	fileName := fmt.Sprintf("scan_%d.jpg", time.Now().Unix())
	_, err := s.ScanImage(dn, fileName)
	if err != nil {
		return "", err
	}

	// Return URL path instead of file path
	url := fmt.Sprintf("/scans/%s", fileName)
	return url, nil
}
