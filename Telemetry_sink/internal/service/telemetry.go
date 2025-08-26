package service

import (
	"context"
	"fmt"
	"log/slog"
	"sensor-sink/sensor-sink/pkg/pb"
	"time"
)

// TelemetryService handles the business logic for telemetry operations
type TelemetryService struct {
	logger *slog.Logger
}

// NewTelemetryService creates a new telemetry service
func NewTelemetryService(logger *slog.Logger) *TelemetryService {
	return &TelemetryService{
		logger: logger,
	}
}

// ProcessReading processes an incoming sensor reading
func (s *TelemetryService) ProcessReading(ctx context.Context, reading *pb.SensorReading) (*pb.SendReadingResponse, error) {
	// Validate the reading
	if err := s.validateReading(reading); err != nil {
		s.logger.Error("Invalid sensor reading",
			"error", err,
			"sensor", reading.GetSensorName(),
		)
		return &pb.SendReadingResponse{
			Success: false,
			Message: fmt.Sprintf("Invalid reading: %v", err),
		}, nil
	}

	// Convert Unix timestamp to time.Time for processing
	timestamp := time.Unix(reading.Timestamp, 0)

	// Log the received reading with structured logging
	s.logger.Info("Received sensor reading",
		"sensor_name", reading.SensorName,
		"value", reading.Value,
		"timestamp", timestamp.Format(time.RFC3339),
	)

	// Print to console for immediate visibility (as per requirements)
	fmt.Printf("[%s] Received from Sensor: %s | Value: %d | Timestamp: %s\n",
		timestamp.Format("2006-01-02 15:04:05"),
		reading.SensorName,
		reading.Value,
		timestamp.Format(time.RFC3339),
	)

	// Here you could add additional processing:
	// - Store to database
	// - Forward to other systems
	// - Apply business rules
	// - Trigger alerts

	return &pb.SendReadingResponse{
		Success: true,
		Message: "Reading processed successfully",
	}, nil
}

// validateReading validates the sensor reading data
func (s *TelemetryService) validateReading(reading *pb.SensorReading) error {
	if reading == nil {
		return fmt.Errorf("reading is nil")
	}

	if reading.SensorName == "" {
		return fmt.Errorf("sensor name is required")
	}

	if reading.Timestamp <= 0 {
		return fmt.Errorf("invalid timestamp")
	}

	// Add more validation rules as needed
	return nil
}
