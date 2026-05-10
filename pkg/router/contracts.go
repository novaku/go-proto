package router

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// AuthRouteHandlers exposes HTTP login/register endpoints implemented by the application
// (e.g. internal/controller.AuthController). The router depends on this abstraction, not on internal packages.
type AuthRouteHandlers interface {
	LoginHandler(w http.ResponseWriter, r *http.Request)
	RegisterHandler(w http.ResponseWriter, r *http.Request)
}

// GatewayRegisterFunc registers grpc-gateway handlers against the runtime mux for a gRPC endpoint.
// Typically the generated pb.RegisterGuestbookServiceHandlerFromEndpoint is passed from cmd.
type GatewayRegisterFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
