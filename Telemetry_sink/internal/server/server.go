package server

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sensor-sink/internal/config"
	"sensor-sink/internal/handler"
	"sensor-sink/internal/service"
	"syscall"

	"google.golang.org/grpc"

	"sensor-sink/sensor-sink/pkg/pb"
)

// Server represents the GRPC server
type Server struct {
	config *config.Config
	logger *slog.Logger
	server *grpc.Server
}

// New creates a new server
func New(cfg *config.Config, logger *slog.Logger) *Server {
	return &Server{
		config: cfg,
		logger: logger,
	}
}

// Start starts the server
func (s *Server) Start() error {
	// Create TCP listener
	listener, err := net.Listen("tcp", s.config.BindAddress)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.config.BindAddress, err)
	}

	// Create GRPC server
	s.server = grpc.NewServer()

	// Create service and handler
	telemetryService := service.NewTelemetryService(s.logger, s.config.LogFilePath)
	telemetryHandler := handler.NewTelemetryHandler(telemetryService, s.logger)

	// Register service
	pb.RegisterSensorServiceServer(s.server, telemetryHandler)

	// Start server
	fmt.Printf("Telemetry Sink Starting...\n")
	fmt.Printf("Bind Address: %s\n", s.config.BindAddress)
	fmt.Printf("Starting gRPC server...\n")
	fmt.Printf("Waiting for sensor readings...\n\n")

	// Handle graceful shutdown
	go s.handleShutdown()

	// Start serving
	return s.server.Serve(listener)
}

// handleShutdown handles graceful shutdown
func (s *Server) handleShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	fmt.Printf("\nReceived signal: %v\n", sig)
	fmt.Println("Shutting down telemetry sink...")

	s.server.GracefulStop()
	fmt.Println("Server stopped")
}
