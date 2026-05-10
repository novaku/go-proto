package framework

import (
	"testing"

	"google.golang.org/grpc"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name string
		port int
		opts []grpc.ServerOption
	}{
		{
			name: "basic server creation",
			port: 50051,
			opts: nil,
		},
		{
			name: "server with custom port",
			port: 9000,
			opts: nil,
		},
		{
			name: "server with options",
			port: 50051,
			opts: []grpc.ServerOption{
				grpc.MaxRecvMsgSize(1024),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewServer(tt.port, tt.opts...)

			if srv == nil {
				t.Fatal("NewServer returned nil")
			}

			if srv.grpcServer == nil {
				t.Error("grpcServer is nil")
			}

			if srv.port != tt.port {
				t.Errorf("port = %d; want %d", srv.port, tt.port)
			}
		})
	}
}

func TestServer_RegisterService(t *testing.T) {
	srv := NewServer(50051)

	registered := false
	srv.RegisterService(func(s *grpc.Server) {
		registered = true
		if s == nil {
			t.Error("Received nil grpc.Server in registration function")
		}
	})

	if !registered {
		t.Error("Service registration function was not called")
	}
}

func TestServer_RegisterMultipleServices(t *testing.T) {
	srv := NewServer(50051)

	count := 0
	srv.RegisterService(func(s *grpc.Server) {
		count++
	})
	srv.RegisterService(func(s *grpc.Server) {
		count++
	})
	srv.RegisterService(func(s *grpc.Server) {
		count++
	})

	if count != 3 {
		t.Errorf("Expected 3 services registered, got %d", count)
	}
}

func TestServer_PortConfiguration(t *testing.T) {
	tests := []struct {
		name string
		port int
	}{
		{"standard port", 50051},
		{"high port", 60000},
		{"low port", 1024},
		{"http alt port", 8080},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewServer(tt.port)
			if srv.port != tt.port {
				t.Errorf("port = %d; want %d", srv.port, tt.port)
			}
		})
	}
}

func TestServer_WithInterceptor(t *testing.T) {
	interceptorCalled := false

	// Create a test interceptor
	srv := NewServer(50051, grpc.ChainUnaryInterceptor(
		UnaryLoggingInterceptor(),
		UnaryRecoveryInterceptor(),
	))

	if srv == nil {
		t.Fatal("NewServer returned nil")
	}

	// Verify server was created successfully with interceptor
	if srv.grpcServer == nil {
		t.Error("grpcServer is nil")
	}

	// The interceptor will be tested when the server actually handles requests
	_ = interceptorCalled
}

func TestServer_NilOptions(t *testing.T) {
	// grpc.NewServer doesn't accept nil options directly,
	// so this test should pass empty options instead
	srv := NewServer(50051)

	if srv == nil {
		t.Fatal("NewServer returned nil")
	}

	if srv.grpcServer == nil {
		t.Error("grpcServer is nil")
	}
}

func TestServer_MultipleOptions(t *testing.T) {
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(1024 * 1024),
		grpc.MaxSendMsgSize(1024 * 1024),
	}

	srv := NewServer(50051, opts...)

	if srv == nil {
		t.Fatal("NewServer returned nil with multiple options")
	}

	if srv.grpcServer == nil {
		t.Error("grpcServer is nil")
	}
}

// TestServer_Run is skipped as it would actually start a server
// and require proper cleanup. Integration tests should cover this.
func TestServer_Run(t *testing.T) {
	t.Skip("Skipping test that would start an actual server")
}

func TestServer_StructureValidation(t *testing.T) {
	srv := NewServer(50051)

	// Validate the server structure
	if srv.listener != nil {
		t.Error("listener should be nil before Run is called")
	}

	if srv.port == 0 {
		t.Error("port should not be 0")
	}

	if srv.grpcServer == nil {
		t.Error("grpcServer should not be nil")
	}
}

func TestServer_ZeroPort(t *testing.T) {
	// Testing with port 0 (OS will assign a random port)
	srv := NewServer(0)

	if srv == nil {
		t.Fatal("NewServer returned nil with port 0")
	}

	if srv.port != 0 {
		t.Errorf("port = %d; want 0", srv.port)
	}
}

func TestServer_NegativePort(t *testing.T) {
	// The server allows negative ports but they will fail at Run time
	srv := NewServer(-1)

	if srv == nil {
		t.Fatal("NewServer returned nil with negative port")
	}

	if srv.port != -1 {
		t.Errorf("port = %d; want -1", srv.port)
	}
}
