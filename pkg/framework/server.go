package framework

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server is a wrapper around grpc.Server with graceful shutdown capabilities.
type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
	port       int
}

// NewServer creates a new Server instance.
func NewServer(port int, opts ...grpc.ServerOption) *Server {
	// Add default interceptors here if needed, or allow passing them via opts
	s := grpc.NewServer(opts...)
	reflection.Register(s) // Enable reflection for tools like grpcurl

	return &Server{
		grpcServer: s,
		port:       port,
	}
}

// RegisterService registers a service implementation with the gRPC server.
func (s *Server) RegisterService(reg func(*grpc.Server)) {
	reg(s.grpcServer)
}

// Run starts the gRPC server and waits for a shutdown signal.
func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	s.listener = lis

	// Handle graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		fmt.Println("\nShutting down gRPC server...")
		s.grpcServer.GracefulStop()
	}()

	fmt.Printf("gRPC server listening on %s\n", addr)
	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
