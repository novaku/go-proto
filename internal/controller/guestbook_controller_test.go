// Unit tests for GuestbookController: gRPC adapter delegates to GuestbookService mocks.
package controller

import (
	"context"
	"testing"

	pb "github.com/novaherdi/go-proto/gen/guestbook/v1"
)

// MockGuestbookService is a mock implementation of GuestbookService for testing
type MockGuestbookService struct {
	AddEntryFunc    func(ctx context.Context, req *pb.AddEntryRequest) (*pb.AddEntryResponse, error)
	ListEntriesFunc func(ctx context.Context, req *pb.ListEntriesRequest) (*pb.ListEntriesResponse, error)
}

func (m *MockGuestbookService) AddEntry(ctx context.Context, req *pb.AddEntryRequest) (*pb.AddEntryResponse, error) {
	if m.AddEntryFunc != nil {
		return m.AddEntryFunc(ctx, req)
	}
	return &pb.AddEntryResponse{Success: true}, nil
}

func (m *MockGuestbookService) ListEntries(ctx context.Context, req *pb.ListEntriesRequest) (*pb.ListEntriesResponse, error) {
	if m.ListEntriesFunc != nil {
		return m.ListEntriesFunc(ctx, req)
	}
	return &pb.ListEntriesResponse{Data: []*pb.GuestbookEntry{}}, nil
}

func TestNewGuestbookController(t *testing.T) {
	mockService := &MockGuestbookService{}
	controller := NewGuestbookController(mockService)

	if controller == nil {
		t.Fatal("expected controller to be non-nil")
	}

	if controller.guestbookService != mockService {
		t.Error("expected controller to have the provided service")
	}
}

func TestGuestbookController_AddEntry(t *testing.T) {
	testCases := []struct {
		name         string
		request      *pb.AddEntryRequest
		serviceResp  *pb.AddEntryResponse
		serviceErr   error
		expectedResp *pb.AddEntryResponse
		expectedErr  error
	}{
		{
			name: "successful addition",
			request: &pb.AddEntryRequest{
				Name:    "John Doe",
				Email:   "john@example.com",
				Message: "Hello, World!",
			},
			serviceResp:  &pb.AddEntryResponse{Success: true},
			serviceErr:   nil,
			expectedResp: &pb.AddEntryResponse{Success: true},
			expectedErr:  nil,
		},
		{
			name: "service returns error",
			request: &pb.AddEntryRequest{
				Name:    "John Doe",
				Message: "Hello, World!",
			},
			serviceResp:  &pb.AddEntryResponse{Success: false, Error: "validation failed"},
			serviceErr:   nil,
			expectedResp: &pb.AddEntryResponse{Success: false, Error: "validation failed"},
			expectedErr:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &MockGuestbookService{
				AddEntryFunc: func(ctx context.Context, req *pb.AddEntryRequest) (*pb.AddEntryResponse, error) {
					return tc.serviceResp, tc.serviceErr
				},
			}

			controller := NewGuestbookController(mockService)
			resp, err := controller.AddEntry(context.Background(), tc.request)

			if (err != nil) != (tc.expectedErr != nil) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}

			if resp.Success != tc.expectedResp.Success {
				t.Errorf("expected success %v, got %v", tc.expectedResp.Success, resp.Success)
			}

			if resp.Error != tc.expectedResp.Error {
				t.Errorf("expected error message %s, got %s", tc.expectedResp.Error, resp.Error)
			}
		})
	}
}

func TestGuestbookController_ListEntries(t *testing.T) {
	testCases := []struct {
		name         string
		request      *pb.ListEntriesRequest
		serviceResp  *pb.ListEntriesResponse
		serviceErr   error
		expectedResp *pb.ListEntriesResponse
		expectedErr  error
	}{
		{
			name:    "successful list",
			request: &pb.ListEntriesRequest{},
			serviceResp: &pb.ListEntriesResponse{
				Data: []*pb.GuestbookEntry{
					{Name: "John", Email: "john@example.com", Message: "Hello", CreatedAt: "2009-02-13T23:31:30Z"},
				},
			},
			serviceErr: nil,
			expectedResp: &pb.ListEntriesResponse{
				Data: []*pb.GuestbookEntry{
					{Name: "John", Email: "john@example.com", Message: "Hello", CreatedAt: "2009-02-13T23:31:30Z"},
				},
			},
			expectedErr: nil,
		},
		{
			name:         "empty list",
			request:      &pb.ListEntriesRequest{},
			serviceResp:  &pb.ListEntriesResponse{Data: []*pb.GuestbookEntry{}},
			serviceErr:   nil,
			expectedResp: &pb.ListEntriesResponse{Data: []*pb.GuestbookEntry{}},
			expectedErr:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &MockGuestbookService{
				ListEntriesFunc: func(ctx context.Context, req *pb.ListEntriesRequest) (*pb.ListEntriesResponse, error) {
					return tc.serviceResp, tc.serviceErr
				},
			}

			controller := NewGuestbookController(mockService)
			resp, err := controller.ListEntries(context.Background(), tc.request)

			if (err != nil) != (tc.expectedErr != nil) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}

			if resp != nil && len(resp.Data) != len(tc.expectedResp.Data) {
				t.Errorf("expected %d entries, got %d", len(tc.expectedResp.Data), len(resp.Data))
			}

			if resp != nil && len(resp.Data) > 0 && resp.Data[0].Name != tc.expectedResp.Data[0].Name {
				t.Errorf("expected first entry name %s, got %s", tc.expectedResp.Data[0].Name, resp.Data[0].Name)
			}
		})
	}
}
