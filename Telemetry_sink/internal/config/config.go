package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the telemetry sink
type Config struct {
	BindAddress string
	LogFilePath string
}

// Load loads configuration from .env file
func Load() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	// Get required BindAddress
	bindAddress := os.Getenv("BIND_ADDRESS")
	if bindAddress == "" {
		return nil, fmt.Errorf("BIND_ADDRESS is required in .env file")
	}

	// Get required LogFilePath
	logFilePath := os.Getenv("LOG_FILE_PATH")
	if logFilePath == "" {
		return nil, fmt.Errorf("LOG_FILE_PATH is required in .env file")
	}

	return &Config{
		BindAddress: bindAddress,
		LogFilePath: logFilePath,
	}, nil
}
