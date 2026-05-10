package auth

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type stubValidator struct {
	claims *JWTClaims
	err    error
}

func (s *stubValidator) ValidateToken(string) (*JWTClaims, error) {
	return s.claims, s.err
}

func TestJWTAuthInterceptor_UnprotectedMethod(t *testing.T) {
	v := &stubValidator{err: errors.New("should not validate")}
	interceptor := JWTAuthInterceptor(v, map[string]bool{})

	info := &grpc.UnaryServerInfo{FullMethod: "/x/Y"}
	handler := func(ctx context.Context, req any) (any, error) { return "ok", nil }

	out, err := interceptor(context.Background(), "req", info, handler)
	if err != nil || out != "ok" {
		t.Fatalf("out=%v err=%v", out, err)
	}
}

func TestJWTAuthInterceptor_Success(t *testing.T) {
	claims := &JWTClaims{UserID: 1, Username: "a", Email: "a@b.c"}
	v := &stubValidator{claims: claims}
	interceptor := JWTAuthInterceptor(v, map[string]bool{"/svc/M": true})

	md := metadata.Pairs(AuthorizationHeader, "Bearer valid.token")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	info := &grpc.UnaryServerInfo{FullMethod: "/svc/M"}
	var gotCtx context.Context
	handler := func(ctx context.Context, req any) (any, error) {
		gotCtx = ctx
		return "done", nil
	}

	out, err := interceptor(ctx, nil, info, handler)
	if err != nil || out != "done" {
		t.Fatalf("out=%v err=%v", out, err)
	}
	c, ok := GetClaimsFromContext(gotCtx)
	if !ok || c.UserID != 1 {
		t.Fatalf("claims not in context: ok=%v c=%v", ok, c)
	}
}

func TestJWTAuthInterceptor_ExpiredToken(t *testing.T) {
	v := &stubValidator{err: ErrExpiredToken}
	interceptor := JWTAuthInterceptor(v, map[string]bool{"/svc/M": true})

	md := metadata.Pairs(AuthorizationHeader, "Bearer x")
	ctx := metadata.NewIncomingContext(context.Background(), md)
	info := &grpc.UnaryServerInfo{FullMethod: "/svc/M"}

	_, err := interceptor(ctx, nil, info, func(context.Context, any) (any, error) { return nil, nil })
	st := status.Convert(err)
	if st.Code() != codes.Unauthenticated {
		t.Fatalf("code %v", st.Code())
	}
}

func TestJWTAuthInterceptor_NoMetadata(t *testing.T) {
	interceptor := JWTAuthInterceptor(&stubValidator{}, map[string]bool{"/svc/M": true})
	info := &grpc.UnaryServerInfo{FullMethod: "/svc/M"}

	_, err := interceptor(context.Background(), nil, info, nil)
	if st := status.Convert(err); st.Code() != codes.Unauthenticated {
		t.Fatalf("got %v", err)
	}
}
