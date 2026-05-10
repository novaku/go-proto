package controller

import (
	"context"

	pb "github.com/novaherdi/go-proto/gen/guestbook/v1"
	"github.com/novaherdi/go-proto/internal/service"
)

// GuestbookController adapts gRPC GuestbookService calls to the application GuestbookService.
// Dependency Inversion: this type depends on service.GuestbookService, not on repositories.
type GuestbookController struct {
	pb.UnimplementedGuestbookServiceServer
	guestbookService service.GuestbookService
}

// NewGuestbookController injects the guestbook use-case implementation.
func NewGuestbookController(guestbookService service.GuestbookService) *GuestbookController {
	return &GuestbookController{
		guestbookService: guestbookService,
	}
}

// AddEntry forwards the RPC to the application service (thin adapter).
func (c *GuestbookController) AddEntry(ctx context.Context, req *pb.AddEntryRequest) (*pb.AddEntryResponse, error) {
	return c.guestbookService.AddEntry(ctx, req)
}

// ListEntries forwards the RPC to the application service (thin adapter).
func (c *GuestbookController) ListEntries(ctx context.Context, req *pb.ListEntriesRequest) (*pb.ListEntriesResponse, error) {
	return c.guestbookService.ListEntries(ctx, req)
}
