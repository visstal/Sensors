package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sensor-sink/sensor-sink/pkg/pb"
	"time"
)

type sensorLogEntry struct {
	SensorName string `json:"sensor_name"`
	Value      int    `json:"value"`
	Timestamp  string `json:"timestamp"`
	ReceivedAt string `json:"received_at"`
}

type TelemetryService struct {
	logger      *slog.Logger
	logFilePath string
}

func NewTelemetryService(logger *slog.Logger, logFilePath string) *TelemetryService {
	return &TelemetryService{
		logger:      logger,
		logFilePath: logFilePath,
	}
}

func (s *TelemetryService) ProcessReading(ctx context.Context, reading *pb.SensorReading) (*pb.SendReadingResponse, error) {
	if err := s.validateReading(reading); err != nil {
		return s.handleValidationError(reading, err)
	}

	timestamp := time.Unix(reading.Timestamp, 0)
	s.logReading(reading, timestamp)

	return s.createSuccessResponse(), nil
}

func (s *TelemetryService) handleValidationError(reading *pb.SensorReading, err error) (*pb.SendReadingResponse, error) {
	s.logger.Error("Invalid sensor reading",
		"error", err,
		"sensor", reading.GetSensorName(),
	)
	return &pb.SendReadingResponse{
		Success: false,
		Message: fmt.Sprintf("Invalid reading: %v", err),
	}, nil
}

func (s *TelemetryService) logReading(reading *pb.SensorReading, timestamp time.Time) {
	if err := s.writeToLogFile(reading, timestamp); err != nil {
		s.logger.Error("Failed to write to log file", "error", err)
	}

	s.logger.Info("Received sensor reading",
		"sensor_name", reading.SensorName,
		"value", reading.Value,
		"timestamp", timestamp.Format(time.RFC3339),
	)
}

func (s *TelemetryService) createSuccessResponse() *pb.SendReadingResponse {
	return &pb.SendReadingResponse{
		Success: true,
		Message: "Reading processed successfully",
	}
}

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

	return nil
}

func (s *TelemetryService) writeToLogFile(reading *pb.SensorReading, timestamp time.Time) error {
	jsonLine, err := s.createLogEntry(reading, timestamp)
	if err != nil {
		return err
	}

	return s.appendToFile(jsonLine)
}

func (s *TelemetryService) createLogEntry(reading *pb.SensorReading, timestamp time.Time) (string, error) {
	entry := sensorLogEntry{
		SensorName: reading.SensorName,
		Value:      int(reading.Value),
		Timestamp:  timestamp.Format(time.RFC3339),
		ReceivedAt: time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		return "", fmt.Errorf("failed to marshal log entry to JSON: %w", err)
	}

	return string(jsonData) + "\n", nil
}

func (s *TelemetryService) appendToFile(logLine string) error {
	file, err := os.OpenFile(s.logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	if _, err = file.WriteString(logLine); err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	return file.Sync()
}
