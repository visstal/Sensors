package grpc

import (
	"context"
	"fmt"
	"log"
	sensor "sensor-node/internal/models"
	pb "sensor-node/sensor-node/pkg/pb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TelemetryClient handles GRPC communication with telemetry sink
type TelemetryClient struct {
	client pb.SensorServiceClient
	conn   *grpc.ClientConn
}

// NewTelemetryClient creates a new GRPC client connection to the telemetry sink
func NewTelemetryClient(sinkAddress string) (*TelemetryClient, error) {
	// Create connection with insecure credentials for now
	conn, err := grpc.Dial(sinkAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to telemetry sink at %s: %v", sinkAddress, err)
	}

	client := pb.NewSensorServiceClient(conn)

	return &TelemetryClient{
		client: client,
		conn:   conn,
	}, nil
}

// SendReading sends a sensor reading to the telemetry sink via GRPC
func (tc *TelemetryClient) SendReading(reading *sensor.SensorReading) error {
	// Convert internal sensor reading to protobuf message
	pbReading := &pb.SensorReading{
		SensorName: reading.SensorName,
		Value:      int32(reading.Value),
		Timestamp:  reading.Timestamp.Unix(),
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Send the reading
	response, err := tc.client.SendReading(ctx, pbReading)
	if err != nil {
		return fmt.Errorf("failed to send reading: %v", err)
	}

	// Log response for debugging (optional)
	if !response.Success {
		log.Printf("Warning: Telemetry sink reported failure: %s", response.Message)
	}

	return nil
}

// Close closes the GRPC connection
func (tc *TelemetryClient) Close() error {
	if tc.conn != nil {
		return tc.conn.Close()
	}
	return nil
}
