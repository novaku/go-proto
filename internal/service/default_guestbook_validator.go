package service

import (
	"errors"
	"strings"

	pb "github.com/novaherdi/go-proto/gen/guestbook/v1"
)

// defaultGuestbookValidator implements GuestbookRequestValidator with basic rules.
// Single Responsibility: input validation only; no I/O or mapping.
type defaultGuestbookValidator struct{}

// NewDefaultValidator returns the default GuestbookRequestValidator implementation.
func NewDefaultValidator() GuestbookRequestValidator {
	return &defaultGuestbookValidator{}
}

// ValidateAddEntryRequest validates the add entry request and normalizes whitespace on the request.
func (v *defaultGuestbookValidator) ValidateAddEntryRequest(req *pb.AddEntryRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Message = strings.TrimSpace(req.Message)

	if req.Name == "" {
		return errors.New("name is required")
	}

	if req.Message == "" {
		return errors.New("message is required")
	}

	return nil
}

// ValidateListEntriesRequest validates the list entries request and returns pagination parameters.
func (v *defaultGuestbookValidator) ValidateListEntriesRequest(req *pb.ListEntriesRequest) (*PaginationParams, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := int(req.Offset)
	if offset < 0 {
		offset = 0
	}

	return &PaginationParams{
		Limit:  limit,
		Offset: offset,
	}, nil
}
