package service

import (
	"context"

	pb "github.com/novaherdi/go-proto/gen/guestbook/v1"
)

// GuestbookService is the application API for guestbook operations (use cases).
type GuestbookService interface {
	AddEntry(ctx context.Context, req *pb.AddEntryRequest) (*pb.AddEntryResponse, error)
	ListEntries(ctx context.Context, req *pb.ListEntriesRequest) (*pb.ListEntriesResponse, error)
}

// AddEntryRequestValidator validates add-entry RPC payloads (Interface Segregation:
// a component may depend only on this facet without listing/pagination rules).
type AddEntryRequestValidator interface {
	ValidateAddEntryRequest(req *pb.AddEntryRequest) error
}

// ListEntriesRequestValidator validates list RPC payloads and normalizes pagination.
type ListEntriesRequestValidator interface {
	ValidateListEntriesRequest(req *pb.ListEntriesRequest) (*PaginationParams, error)
}

// GuestbookRequestValidator composes both validation facets for the full guestbook flow.
// Types like defaultGuestbookValidator implement all methods at once (Open/Closed via substitution).
type GuestbookRequestValidator interface {
	AddEntryRequestValidator
	ListEntriesRequestValidator
}

// PaginationParams holds validated pagination parameters for list queries.
type PaginationParams struct {
	Limit  int
	Offset int
}
