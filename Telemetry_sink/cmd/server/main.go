package main

import (
	"log/slog"
	"os"

	"sensor-sink/internal/config"
	"sensor-sink/internal/server"
)

func main() {
	// Create logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Create and start server
	srv := server.New(cfg, logger)
	if err := srv.Start(); err != nil {
		logger.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
