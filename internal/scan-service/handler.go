package scanservice

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"go-sane/cmd/web"

	"github.com/a-h/templ"
)

type ScanHandler struct {
	scanner Scanner
}

func NewScanHandler(scanner Scanner) *ScanHandler {
	return &ScanHandler{scanner: scanner}
}

func (h *ScanHandler) Devices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	devices, err := h.scanner.ListDevices()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	names := make([]string, len(devices))
	for i, device := range devices {
		names[i] = device.Name
	}

	templ.Handler(web.Devices(names)).ServeHTTP(w, r)
}

func (h *ScanHandler) DevicesOptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	devices, err := h.scanner.ListDevices()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	names := make([]string, len(devices))
	for i, device := range devices {
		names[i] = device.Name
	}

	templ.Handler(web.DeviceOptions(names)).ServeHTTP(w, r)
}

func (h *ScanHandler) DeviceStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	deviceName := r.URL.Query().Get("device")
	if deviceName == "" {
		templ.Handler(web.SelectedDeviceStatus("", "Please select a scanner.")).ServeHTTP(w, r)
		return
	}

	templ.Handler(web.SelectedDeviceStatus(deviceName, "Ready to scan.")).ServeHTTP(w, r)
}

func (h *ScanHandler) Scan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	deviceName := r.FormValue("device")
	if deviceName == "" {
		http.Error(w, "device parameter required", http.StatusBadRequest)
		return
	}

	err := h.scanner.ValidateDeviceName(deviceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Scan image and get URL
	url, err := h.scanner.ScanAndReturnURL(deviceName)
	if err != nil {
		if errors.Is(err, ErrIsScanning) {
			http.Error(w, err.Error(), http.StatusTooManyRequests)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return HTML fragment with the image
	templ.Handler(web.ScanResult(url)).ServeHTTP(w, r)
}

func (h *ScanHandler) ScanSSE(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	deviceName := r.FormValue("device")
	if deviceName == "" {
		http.Error(w, "device parameter required", http.StatusBadRequest)
		return
	}

	err := h.scanner.ValidateDeviceName(deviceName)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Perform scan synchronously to avoid connection issues
	url, err := h.scanner.ScanAndReturnURL(deviceName)
	if err != nil {
		slog.Error(err.Error())
		fmt.Fprintf(w, "event: scan-error\ndata: %s\n\n", err.Error())
		flusher.Flush()
		return
	}

	web.ScanResult(url).Render(r.Context(), w)
	flusher.Flush()
}
