package config

import (
	"flag"
	"fmt"
)

// Config holds all configuration for the sensor node
type Config struct {
	Rate        int    // messages per second
	SensorName  string // name of the sensor
	SinkAddress string // address of telemetry sink
	BindAddress string // bind address for sensor node (if needed later)
}

// LoadFromCLI parses command line flags and returns configuration
func LoadFromCLI() (*Config, error) {
	cfg := &Config{}

	flag.IntVar(&cfg.Rate, "rate", 0, "Number of messages per second to send (REQUIRED)")
	flag.StringVar(&cfg.SensorName, "sensor-name", "", "Name of the sensor (REQUIRED)")
	flag.StringVar(&cfg.SinkAddress, "sink-address", "", "Address of the telemetry sink (REQUIRED)")
	flag.StringVar(&cfg.BindAddress, "bind-address", ":0", "Bind address for the sensor node (optional, default random port)")

	flag.Parse()

	// Validate all required fields
	if cfg.Rate <= 0 {
		return nil, fmt.Errorf("REQUIRED: -rate must be provided and positive, got %d", cfg.Rate)
	}

	if cfg.SensorName == "" {
		return nil, fmt.Errorf("REQUIRED: -sensor-name must be provided")
	}

	if cfg.SinkAddress == "" {
		return nil, fmt.Errorf("REQUIRED: -sink-address must be provided")
	}

	return cfg, nil
}
