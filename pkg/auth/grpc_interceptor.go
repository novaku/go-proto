package auth

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// JWTAuthInterceptor returns a unary interceptor that enforces bearer JWT on selected methods.
// It depends on TokenValidator (Dependency Inversion), not on concrete *JWTService.
func JWTAuthInterceptor(validator TokenValidator, protectedMethods map[string]bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if !protectedMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		authHeader := md.Get(AuthorizationHeader)
		if len(authHeader) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		token := authHeader[0]
		if strings.HasPrefix(strings.ToLower(token), "bearer ") {
			token = token[7:]
		}

		claims, err := validator.ValidateToken(token)
		if err != nil {
			if errors.Is(err, ErrExpiredToken) {
				return nil, status.Errorf(codes.Unauthenticated, "token has expired")
			}
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		ctx = context.WithValue(ctx, UserContextKey, claims)
		return handler(ctx, req)
	}
}
