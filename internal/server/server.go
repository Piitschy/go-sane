package server

import (
	"fmt"
	scanservice "go-sane/internal/scan-service"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
	sh   *scanservice.ScanHandler
}

func NewServer(sh *scanservice.ScanHandler) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
		sh:   sh,
	}
	slog.Info(fmt.Sprintf("listen on port %d", port))

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  5 * time.Minute,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 2 * time.Minute,
	}

	return server
}
