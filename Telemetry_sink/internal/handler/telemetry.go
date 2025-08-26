package handler

import (
	"context"
	"log/slog"

	"sensor-sink/internal/service"
	pb "sensor-sink/sensor-sink/pkg/pb"
)

// TelemetryHandler handles GRPC requests for telemetry
type TelemetryHandler struct {
	pb.UnimplementedSensorServiceServer
	telemetryService *service.TelemetryService
	logger           *slog.Logger
}

// NewTelemetryHandler creates a new telemetry handler
func NewTelemetryHandler(telemetryService *service.TelemetryService, logger *slog.Logger) *TelemetryHandler {
	return &TelemetryHandler{
		telemetryService: telemetryService,
		logger:           logger,
	}
}

// SendReading handles incoming sensor readings
func (h *TelemetryHandler) SendReading(ctx context.Context, reading *pb.SensorReading) (*pb.SendReadingResponse, error) {
	return h.telemetryService.ProcessReading(ctx, reading)
}
