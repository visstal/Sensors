package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the telemetry sink
type Config struct {
	BindAddress string
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

	return &Config{
		BindAddress: bindAddress,
	}, nil
}
