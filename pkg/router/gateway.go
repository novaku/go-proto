package router

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RegisterGateway wires the grpc-gateway mux to the gRPC server using the provided registration function.
func (r *Router) RegisterGateway(register GatewayRegisterFunc) error {
	if r.config == nil {
		return fmt.Errorf("router config is nil")
	}
	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	endpoint := fmt.Sprintf("localhost:%d", r.config.Server.Port)
	return register(ctx, r.gwMux, endpoint, opts)
}
