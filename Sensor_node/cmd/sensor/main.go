package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sensor-node/internal/config"
	"sensor-node/pkg/datasource"
	"sensor-node/pkg/grpc"
)

func main() {
	// Load configuration from CLI
	cfg, err := config.LoadFromCLI()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Printf("Sensor Node Starting...\n")
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Rate: %d messages/second\n", cfg.Rate)
	fmt.Printf("  Sensor Name: %s\n", cfg.SensorName)
	fmt.Printf("  Sink Address: %s\n", cfg.SinkAddress)
	fmt.Printf("\n")

	// Create GRPC client
	grpcClient, err := grpc.NewTelemetryClient(cfg.SinkAddress)
	if err != nil {
		log.Fatalf("Failed to create GRPC client: %v", err)
	}
	defer grpcClient.Close()

	fmt.Printf("Connected to telemetry sink at %s\n\n", cfg.SinkAddress)

	// Create sensor data source
	dataSource := datasource.NewSensorDataSource(cfg.SensorName)

	// Calculate interval between messages
	interval := time.Second / time.Duration(cfg.Rate)

	// Create ticker for periodic data generation
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("Starting sensor data generation (Press Ctrl+C to stop)...\n\n")

	// Main application loop
	for {
		select {
		case <-ticker.C:
			// Generate sensor reading
			reading := dataSource.GetReading()

			// Output to console
			fmt.Printf("[%s] Sensor: %s | Value: %d | Timestamp: %s\n",
				reading.Timestamp.Format("2006-01-02 15:04:05"),
				reading.SensorName,
				reading.Value,
				reading.Timestamp.Format(time.RFC3339),
			)

			// Send to telemetry sink via GRPC
			if err := grpcClient.SendReading(reading); err != nil {
				log.Printf("Failed to send reading via GRPC: %v", err)
			}

		case sig := <-sigChan:
			fmt.Printf("\nReceived signal: %v\n", sig)
			fmt.Println("Shutting down sensor node...")
			return
		}
	}
}
