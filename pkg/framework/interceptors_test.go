package framework

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnaryLoggingInterceptor(t *testing.T) {
	interceptor := UnaryLoggingInterceptor()

	if interceptor == nil {
		t.Fatal("UnaryLoggingInterceptor returned nil")
	}

	tests := []struct {
		name         string
		handler      grpc.UnaryHandler
		wantErr      bool
		expectedCode codes.Code
	}{
		{
			name: "successful request",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return "success", nil
			},
			wantErr:      false,
			expectedCode: codes.OK,
		},
		{
			name: "request with error",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, status.Error(codes.InvalidArgument, "invalid argument")
			},
			wantErr:      true,
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "request with internal error",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, status.Error(codes.Internal, "internal error")
			},
			wantErr:      true,
			expectedCode: codes.Internal,
		},
		{
			name: "request with not found error",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, status.Error(codes.NotFound, "not found")
			},
			wantErr:      true,
			expectedCode: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			req := "test request"
			info := &grpc.UnaryServerInfo{
				FullMethod: "/test.Service/TestMethod",
			}

			resp, err := interceptor(ctx, req, info, tt.handler)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnaryLoggingInterceptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				st := status.Convert(err)
				if st.Code() != tt.expectedCode {
					t.Errorf("Expected code %v, got %v", tt.expectedCode, st.Code())
				}
			}

			if !tt.wantErr && resp == nil {
				t.Error("Expected non-nil response for successful request")
			}
		})
	}
}

func TestUnaryRecoveryInterceptor_Panic(t *testing.T) {
	interceptor := UnaryRecoveryInterceptor()

	ctx := context.Background()
	req := "test request"
	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/PanicMethod",
	}

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("test panic")
	}

	resp, err := interceptor(ctx, req, info, handler)

	if err == nil {
		t.Error("Expected error after panic, got nil")
	}

	if resp != nil {
		t.Error("Expected nil response after panic")
	}

	st := status.Convert(err)
	if st.Code() != codes.Internal {
		t.Errorf("Expected Internal code after panic, got %v", st.Code())
	}

	if st.Message() != "internal server error" {
		t.Errorf("Expected 'internal server error' message, got '%s'", st.Message())
	}
}

func TestUnaryRecoveryInterceptor_PanicWithString(t *testing.T) {
	interceptor := UnaryRecoveryInterceptor()

	ctx := context.Background()
	req := "test request"
	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/PanicMethod",
	}

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("string panic message")
	}

	_, err := interceptor(ctx, req, info, handler)

	if err == nil {
		t.Error("Expected error after panic, got nil")
	}

	st := status.Convert(err)
	if st.Code() != codes.Internal {
		t.Errorf("Expected Internal code, got %v", st.Code())
	}
}

func TestUnaryRecoveryInterceptor_PanicWithError(t *testing.T) {
	interceptor := UnaryRecoveryInterceptor()

	ctx := context.Background()
	req := "test request"
	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/PanicMethod",
	}

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		panic(errors.New("error panic"))
	}

	_, err := interceptor(ctx, req, info, handler)

	if err == nil {
		t.Error("Expected error after panic, got nil")
	}
}

func TestUnaryLoggingInterceptor_Timing(t *testing.T) {
	interceptor := UnaryLoggingInterceptor()

	ctx := context.Background()
	req := "test request"
	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/SlowMethod",
	}

	sleepDuration := 10 * time.Millisecond
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		time.Sleep(sleepDuration)
		return "success", nil
	}

	start := time.Now()
	_, err := interceptor(ctx, req, info, handler)
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if elapsed < sleepDuration {
		t.Errorf("Expected elapsed time >= %v, got %v", sleepDuration, elapsed)
	}
}

func TestUnaryLoggingInterceptor_WithContext(t *testing.T) {
	interceptor := UnaryLoggingInterceptor()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "test-key", "test-value")

	req := "test request"
	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}

	contextReceived := false
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		if ctx.Value("test-key") == "test-value" {
			contextReceived = true
		}
		return "success", nil
	}

	_, err := interceptor(ctx, req, info, handler)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !contextReceived {
		t.Error("Context was not properly passed to handler")
	}
}

func TestUnaryLoggingInterceptor_DifferentMethods(t *testing.T) {
	interceptor := UnaryLoggingInterceptor()

	methods := []string{
		"/guestbook.v1.GuestbookService/AddEntry",
		"/guestbook.v1.GuestbookService/ListEntries",
		"/health.v1.Health/Check",
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			ctx := context.Background()
			req := "test request"
			info := &grpc.UnaryServerInfo{
				FullMethod: method,
			}

			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return "success", nil
			}

			_, err := interceptor(ctx, req, info, handler)

			if err != nil {
				t.Errorf("Unexpected error for method %s: %v", method, err)
			}
		})
	}
}

func TestUnaryLoggingInterceptor_NilRequest(t *testing.T) {
	interceptor := UnaryLoggingInterceptor()

	ctx := context.Background()
	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		if req == nil {
			return nil, status.Error(codes.InvalidArgument, "nil request")
		}
		return "success", nil
	}

	_, err := interceptor(ctx, nil, info, handler)

	if err == nil {
		t.Error("Expected error for nil request")
	}
}

func TestUnaryLoggingInterceptor_StatusCodes(t *testing.T) {
	interceptor := UnaryLoggingInterceptor()

	testCodes := []codes.Code{
		codes.OK,
		codes.Canceled,
		codes.Unknown,
		codes.InvalidArgument,
		codes.DeadlineExceeded,
		codes.NotFound,
		codes.AlreadyExists,
		codes.PermissionDenied,
		codes.ResourceExhausted,
		codes.FailedPrecondition,
		codes.Aborted,
		codes.OutOfRange,
		codes.Unimplemented,
		codes.Internal,
		codes.Unavailable,
		codes.DataLoss,
		codes.Unauthenticated,
	}

	for _, code := range testCodes {
		t.Run(code.String(), func(t *testing.T) {
			ctx := context.Background()
			req := "test request"
			info := &grpc.UnaryServerInfo{
				FullMethod: "/test.Service/TestMethod",
			}

			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				if code == codes.OK {
					return "success", nil
				}
				return nil, status.Error(code, "test error")
			}

			_, err := interceptor(ctx, req, info, handler)

			if code == codes.OK {
				if err != nil {
					t.Errorf("Expected no error for OK code, got %v", err)
				}
			} else {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				st := status.Convert(err)
				if st.Code() != code {
					t.Errorf("Expected code %v, got %v", code, st.Code())
				}
			}
		})
	}
}
