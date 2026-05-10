package service

import (
	"context"
	"fmt"
	"time"

	pb "github.com/novaherdi/go-proto/gen/guestbook/v1"
	"github.com/novaherdi/go-proto/internal/mapper"
	"github.com/novaherdi/go-proto/internal/repository"
	"github.com/novaherdi/go-proto/pkg/cache"
)

// guestbookService implements GuestbookService (business orchestration).
// It depends on repository, validation, and optional cache abstractions (Dependency Inversion).
type guestbookService struct {
	repo      repository.GuestbookRepository
	validator GuestbookRequestValidator
	cache     cache.Cache
}

// NewGuestbookService creates a guestbook service without caching.
func NewGuestbookService(repo repository.GuestbookRepository, validator GuestbookRequestValidator) GuestbookService {
	return &guestbookService{
		repo:      repo,
		validator: validator,
		cache:     nil,
	}
}

// NewGuestbookServiceWithCache creates a guestbook service with list caching.
func NewGuestbookServiceWithCache(repo repository.GuestbookRepository, validator GuestbookRequestValidator, c cache.Cache) GuestbookService {
	return &guestbookService{
		repo:      repo,
		validator: validator,
		cache:     c,
	}
}

// AddEntry validates input, maps to domain, persists, and invalidates list caches.
func (s *guestbookService) AddEntry(ctx context.Context, req *pb.AddEntryRequest) (*pb.AddEntryResponse, error) {
	if err := s.validator.ValidateAddEntryRequest(req); err != nil {
		return &pb.AddEntryResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	entry := mapper.GuestbookEntryFromAddRequest(req)

	if err := s.repo.Create(ctx, entry); err != nil {
		return &pb.AddEntryResponse{
			Success: false,
			Error:   "Failed to save entry",
		}, fmt.Errorf("failed to create entry: %w", err)
	}

	if s.cache != nil {
		_ = s.cache.Delete(ctx, "listentries:*")
	}

	return &pb.AddEntryResponse{
		Success: true,
	}, nil
}

// ListEntries returns a paginated list, optionally served from cache.
func (s *guestbookService) ListEntries(ctx context.Context, req *pb.ListEntriesRequest) (*pb.ListEntriesResponse, error) {
	params, err := s.validator.ValidateListEntriesRequest(req)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	cacheKey := cache.ListEntriesCacheKey(req.Limit, req.Offset)

	if s.cache != nil {
		var cachedResponse pb.ListEntriesResponse
		if found, getErr := s.cache.Get(ctx, cacheKey, &cachedResponse); found && getErr == nil {
			return &cachedResponse, nil
		}
	}

	dbEntries, err := s.repo.FindWithPagination(ctx, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve entries: %w", err)
	}

	response := &pb.ListEntriesResponse{
		Data: mapper.GuestbookEntriesToProto(dbEntries),
	}

	if s.cache != nil {
		_ = s.cache.Set(ctx, cacheKey, response, time.Hour)
	}

	return response, nil
}
