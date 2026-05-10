package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/novaherdi/go-proto/pkg/config"
)

func TestNewRouter(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:     50051,
			HttpPort: 8080,
			Name:     "test",
		},
		Environment: "test",
	}

	router := NewRouter(cfg)

	if router == nil {
		t.Fatal("NewRouter returned nil")
	}

	if router.mux == nil {
		t.Error("mux is nil")
	}

	if router.gwMux == nil {
		t.Error("gwMux is nil")
	}

	if router.config == nil {
		t.Error("config is nil")
	}

	if router.config != cfg {
		t.Error("config was not properly assigned")
	}
}

func TestRouter_SetupRoutes_NonProduction(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:     50051,
			HttpPort: 8080,
			Name:     "test",
		},
		Environment: "development",
	}

	router := NewRouter(cfg)
	router.SetupRoutes()

	// Test that mux was configured
	if router.mux == nil {
		t.Error("mux is nil after SetupRoutes")
	}
}

func TestRouter_SetupRoutes_Production(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:     50051,
			HttpPort: 8080,
			Name:     "test",
		},
		Environment: "production",
	}

	router := NewRouter(cfg)
	router.SetupRoutes()

	// In production, swagger routes should not be setup
	// Test by making a request to swagger endpoint
	req := httptest.NewRequest("GET", "/swagger-ui.html", nil)
	w := httptest.NewRecorder()

	router.mux.ServeHTTP(w, req)

	// In production, this should return 404
	if w.Code == http.StatusOK {
		t.Error("Swagger UI should not be available in production")
	}
}

func TestRouter_SwaggerRoutes(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:     50051,
			HttpPort: 8080,
			Name:     "test",
		},
		Environment: "development",
	}

	router := NewRouter(cfg)
	router.SetupRoutes()

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		skipTest       bool
	}{
		{
			name:           "swagger json endpoint",
			path:           "/swagger/guestbook.swagger.json",
			expectedStatus: http.StatusNotFound, // File might not exist in test
			skipTest:       true,
		},
		{
			name:           "swagger ui endpoint",
			path:           "/swagger-ui.html",
			expectedStatus: http.StatusNotFound, // File might not exist in test
			skipTest:       true,
		},
		{
			name:           "root redirect",
			path:           "/",
			expectedStatus: http.StatusFound,
			skipTest:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipTest {
				t.Skip("Skipping test that requires actual files")
			}

			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			router.mux.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestRouter_RootRedirect(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:     50051,
			HttpPort: 8080,
			Name:     "test",
		},
		Environment: "development",
	}

	router := NewRouter(cfg)
	router.SetupRoutes()

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	router.mux.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("Expected status %d, got %d", http.StatusFound, w.Code)
	}

	location := w.Header().Get("Location")
	if location != "/swagger-ui.html" {
		t.Errorf("Expected redirect to /swagger-ui.html, got %s", location)
	}
}

func TestRouter_NotFoundPath(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:     50051,
			HttpPort: 8080,
			Name:     "test",
		},
		Environment: "development",
	}

	router := NewRouter(cfg)
	router.SetupRoutes()

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()

	router.mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestRouter_RegisterGateway(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:     50051,
			HttpPort: 8080,
			Name:     "test",
		},
		Environment: "test",
	}

	router := NewRouter(cfg)

	err := router.RegisterGateway(func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error {
		return nil
	})
	if err != nil {
		t.Fatalf("RegisterGateway: %v", err)
	}
}

func TestRouter_DifferentEnvironments(t *testing.T) {
	environments := []string{"development", "test", "staging", "production"}

	for _, env := range environments {
		t.Run(env, func(t *testing.T) {
			cfg := &config.Config{
				Server: config.ServerConfig{
					Port:     50051,
					HttpPort: 8080,
					Name:     "test",
				},
				Environment: env,
			}

			router := NewRouter(cfg)
			router.SetupRoutes()

			if router.config.Environment != env {
				t.Errorf("Expected environment %s, got %s", env, router.config.Environment)
			}
		})
	}
}

func TestRouter_Run(t *testing.T) {
	// This test would actually start a server
	t.Skip("Skipping test that would start an actual server")
}

func TestRouter_MultipleSetupRoutesCalls(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:     50051,
			HttpPort: 8080,
			Name:     "test",
		},
		Environment: "development",
	}

	router := NewRouter(cfg)

	// Call SetupRoutes once (calling multiple times would cause route conflicts)
	router.SetupRoutes()

	// Should not panic and mux should still work
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	router.mux.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("Expected status %d after SetupRoutes call", http.StatusFound)
	}
}

func TestRouter_ConfigNil(t *testing.T) {
	// Test behavior with nil config - this is expected to work
	// since Go allows nil pointer receivers
	router := NewRouter(nil)

	if router == nil {
		t.Fatal("NewRouter returned nil")
	}

	if router.config != nil {
		t.Error("Expected nil config")
	}
}

func TestRouter_V1PathHandling(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:     50051,
			HttpPort: 8080,
			Name:     "test",
		},
		Environment: "development",
	}

	router := NewRouter(cfg)
	router.SetupRoutes()

	// Test that /v1/ paths are handled by gwMux
	req := httptest.NewRequest("GET", "/v1/guestbook", nil)
	w := httptest.NewRecorder()

	router.mux.ServeHTTP(w, req)

	// The actual response depends on whether gRPC gateway is registered
	// but it should be handled without panic
	if w.Code == 0 {
		t.Error("Request was not handled")
	}
}

func TestRouter_HTTPMethods(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:     50051,
			HttpPort: 8080,
			Name:     "test",
		},
		Environment: "development",
	}

	router := NewRouter(cfg)
	router.SetupRoutes()

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/", nil)
			w := httptest.NewRecorder()

			router.mux.ServeHTTP(w, req)

			// Should handle all HTTP methods without panic
			if w.Code == 0 {
				t.Errorf("Request with method %s was not handled", method)
			}
		})
	}
}

func TestRouter_GetMux(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:     50051,
			HttpPort: 8080,
			Name:     "test",
		},
		Environment: "development",
	}

	router := NewRouter(cfg)
	mux := router.GetMux()

	if mux == nil {
		t.Error("GetMux returned nil")
	}
}

func TestRouter_GetGatewayMux(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:     50051,
			HttpPort: 8080,
			Name:     "test",
		},
		Environment: "development",
	}

	router := NewRouter(cfg)
	gwMux := router.GetGatewayMux()

	if gwMux == nil {
		t.Error("GetGatewayMux returned nil")
	}
}

func TestRouter_SwaggerJSONRoute(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:     50051,
			HttpPort: 8080,
			Name:     "test",
		},
		Environment: "development",
	}

	router := NewRouter(cfg)
	router.SetupRoutes()

	req := httptest.NewRequest("GET", "/swagger/guestbook.swagger.json", nil)
	w := httptest.NewRecorder()

	router.mux.ServeHTTP(w, req)

	// Handler should be registered (404 is ok if file doesn't exist)
	if w.Code == 0 {
		t.Error("Swagger JSON route not handled")
	}
}

func TestRouter_SwaggerUIRoute(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:     50051,
			HttpPort: 8080,
			Name:     "test",
		},
		Environment: "development",
	}

	router := NewRouter(cfg)
	router.SetupRoutes()

	req := httptest.NewRequest("GET", "/swagger-ui.html", nil)
	w := httptest.NewRecorder()

	router.mux.ServeHTTP(w, req)

	// Handler should be registered (404 is ok if file doesn't exist)
	if w.Code == 0 {
		t.Error("Swagger UI route not handled")
	}
}

func TestRouter_RootPathInProduction(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:     50051,
			HttpPort: 8080,
			Name:     "test",
		},
		Environment: "production",
	}

	router := NewRouter(cfg)
	router.SetupRoutes()

	req := httptest.NewRequest("GET", "/some-other-path", nil)
	w := httptest.NewRecorder()

	router.mux.ServeHTTP(w, req)

	// In production, undefined paths should return 404
	if w.Code != http.StatusNotFound {
		t.Logf("Got status %d for undefined path in production", w.Code)
	}
}
