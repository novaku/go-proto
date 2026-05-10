package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/novaherdi/go-proto/pkg/config"
)

// Router bundles the stdlib mux, grpc-gateway mux, and route wiring.
type Router struct {
	mux            *http.ServeMux
	gwMux          *runtime.ServeMux
	config         *config.Config
	authHandlers   AuthRouteHandlers
}

// NewRouter creates an HTTP router with an empty grpc-gateway mux.
func NewRouter(cfg *config.Config) *Router {
	return &Router{
		mux:    http.NewServeMux(),
		gwMux:  runtime.NewServeMux(),
		config: cfg,
	}
}

// SetAuthHandlers registers HTTP auth handlers (login/register).
func (r *Router) SetAuthHandlers(h AuthRouteHandlers) {
	r.authHandlers = h
}

// SetAuthController is an alias for SetAuthHandlers for backward compatibility at call sites.
func (r *Router) SetAuthController(c AuthRouteHandlers) {
	r.SetAuthHandlers(c)
}

// SetupRoutes mounts auth, API, and optional Swagger routes.
func (r *Router) SetupRoutes() {
	if r.authHandlers != nil {
		r.mux.HandleFunc("/auth/login", r.authHandlers.LoginHandler)
		r.mux.HandleFunc("/auth/register", r.authHandlers.RegisterHandler)
	}

	r.mux.Handle("/v1/", http.StripPrefix("/v1", r.gwMux))

	if r.config != nil && r.config.Environment != "production" {
		r.setupSwaggerRoutes()
	}
}

// Run listens on the given HTTP port and blocks.
func (r *Router) Run(port int) error {
	addr := fmt.Sprintf(":%d", port)
	log.Printf("HTTP Gateway listening on %s", addr)
	log.Printf("Swagger UI available at: http://localhost:%d/swagger-ui.html", port)
	log.Printf("\nAuth routes:")
	log.Printf("  POST /auth/login    - Login and get JWT token")
	log.Printf("  POST /auth/register - Register new user")
	log.Printf("\nProto-defined routes:")
	log.Printf("  POST /v1/guestbook - AddEntry (requires JWT auth)")
	log.Printf("  GET  /v1/guestbook - ListEntries")

	return http.ListenAndServe(addr, r.mux)
}

// GetMux returns the underlying http.ServeMux.
func (r *Router) GetMux() *http.ServeMux {
	return r.mux
}

// GetGatewayMux returns the grpc-gateway ServeMux.
func (r *Router) GetGatewayMux() *runtime.ServeMux {
	return r.gwMux
}
