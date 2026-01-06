# 📸 go-sane

> A modern web-based scanning solution built with Go and the SANE (Scanner Access Now Easy) library

**go-sane** is a powerful and elegant scanning service that transforms your scanner into a web-accessible device. Built with performance and simplicity in mind, it provides a clean API and web interface for digitizing documents directly from your browser.

## ✨ Features

- **🌐 Web Interface**: Modern, responsive UI for easy scanning operations
- **📡 RESTful API**: Clean HTTP API for integration with other applications  
- **🖨️ SANE Integration**: Full support for SANE-compatible scanners
- **🔄 Live Reloading**: Hot-reload during development with `make watch`
- **📱 Device Discovery**: Automatic scanner detection and device management
- **🎯 High-Quality Scans**: Configurable resolution and color modes
- **🔒 Thread-Safe**: Concurrent-safe scanning operations
- **📁 File Management**: Automatic scan organization and storage

## 🚀 Quick Start

### Prerequisites

- Go 1.25.5 or higher
- SANE installed on your system
- A compatible scanner device

### Installation

```bash
git clone https://github.com/piitschy/go-sane.git
cd go-sane
go mod download
```

### Running the Application

```bash
# Build and run with live reload
make watch

# Or build and run manually
make build
make run
```

Once running, open your browser and navigate to `http://localhost:8080` to access the scanning interface.

## 🛠️ Development

### Make Commands

| Command | Description |
|---------|-------------|
| `make all` | Build the application with tests |
| `make build` | Build the application binary |
| `make run` | Run the application |
| `make watch` | Live reload the application during development |
| `make test` | Run the test suite |
| `make itest` | Run integration tests |
| `make clean` | Clean up build artifacts |

### API Endpoints

The application provides RESTful endpoints for scanner operations:

- `GET /devices` - List available scanners
- `POST /scan` - Initiate a scan operation
- `GET /scans/{filename}` - Access scanned images

### Configuration

Default settings include:
- **Scan Resolution**: 300 DPI
- **Color Mode**: Color
- **Scan Area**: Full A4 (215×297mm)
- **Output Format**: JPEG
- **Storage Path**: `./scans/`

## 📋 Project Structure

```
go-sane/
├── cmd/
│   ├── api/           # HTTP API server
│   └── web/           # Web interface templates
├── internal/
│   ├── scan-service/  # Core scanning logic
│   └── server/        # HTTP server and routing
├── Makefile
└── go.mod
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Dependencies

- [templ](https://github.com/a-h/templ) - Modern Go templating
- [sane](https://github.com/tjgq/sane) - Go bindings for SANE
- [websocket](https://github.com/coder/websocket) - WebSocket support
- [godotenv](https://github.com/joho/godotenv) - Environment variable management

---

**Built with ❤️ using Go and SANE**
