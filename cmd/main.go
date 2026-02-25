package main

import (
	"log/slog"
	"os"

	server "github.com/tdatIT/backend-go/internal"
)

func main() {
	service, err := server.InitServer()
	if err != nil {
		slog.Error("failed to initialize server", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := service.StartHTTP(); err != nil {
		slog.Error("failed to start HTTP server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
