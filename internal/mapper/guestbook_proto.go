// Package mapper converts between domain models and API (protobuf) types.
// Keeping mapping out of services preserves Single Responsibility and Open/Closed:
// new API shapes can add mappers without rewriting business rules.
package mapper

import (
	pb "github.com/novaherdi/go-proto/gen/guestbook/v1"
	"github.com/novaherdi/go-proto/internal/model"
)

// GuestbookEntryToProto maps one domain entry to its protobuf representation.
func GuestbookEntryToProto(e *model.GuestbookEntry) *pb.GuestbookEntry {
	if e == nil {
		return nil
	}
	return &pb.GuestbookEntry{
		Name:      e.Name,
		Email:     e.Email,
		Message:   e.Message,
		CreatedAt: e.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// GuestbookEntriesToProto maps a slice of domain entries to protobuf messages.
func GuestbookEntriesToProto(entries []model.GuestbookEntry) []*pb.GuestbookEntry {
	out := make([]*pb.GuestbookEntry, 0, len(entries))
	for i := range entries {
		out = append(out, GuestbookEntryToProto(&entries[i]))
	}
	return out
}

// GuestbookEntryFromAddRequest builds a domain entity from an AddEntry RPC request.
// Validation is assumed to have run earlier in the pipeline.
func GuestbookEntryFromAddRequest(req *pb.AddEntryRequest) *model.GuestbookEntry {
	if req == nil {
		return nil
	}
	return &model.GuestbookEntry{
		Name:    req.Name,
		Email:   req.Email,
		Message: req.Message,
	}
}
