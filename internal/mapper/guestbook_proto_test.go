package mapper

import (
	"testing"
	"time"

	pb "github.com/novaherdi/go-proto/gen/guestbook/v1"
	"github.com/novaherdi/go-proto/internal/model"
)

func TestGuestbookEntryToProto(t *testing.T) {
	if GuestbookEntryToProto(nil) != nil {
		t.Fatal("nil input should return nil")
	}

	created := time.Date(2024, 7, 1, 12, 30, 45, 0, time.UTC)
	entry := &model.GuestbookEntry{
		ID:        1,
		Name:      "N",
		Email:     "e@x.com",
		Message:   "M",
		CreatedAt: created,
	}
	got := GuestbookEntryToProto(entry)
	if got.Name != "N" || got.Email != "e@x.com" || got.Message != "M" {
		t.Errorf("unexpected fields: %+v", got)
	}
	wantTime := created.Format("2006-01-02 15:04:05")
	if got.CreatedAt != wantTime {
		t.Errorf("CreatedAt = %q, want %q", got.CreatedAt, wantTime)
	}
}

func TestGuestbookEntriesToProto(t *testing.T) {
	entries := []model.GuestbookEntry{
		{Name: "a", Email: "a@b.c", Message: "1", CreatedAt: time.Unix(100, 0).UTC()},
		{Name: "b", Email: "b@b.c", Message: "2", CreatedAt: time.Unix(200, 0).UTC()},
	}
	got := GuestbookEntriesToProto(entries)
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].Name != "a" || got[1].Name != "b" {
		t.Errorf("unexpected slice: %+v", got)
	}
}

func TestGuestbookEntryFromAddRequest(t *testing.T) {
	if GuestbookEntryFromAddRequest(nil) != nil {
		t.Fatal("nil request should return nil")
	}

	req := &pb.AddEntryRequest{Name: "X", Email: "x@y.z", Message: "hello"}
	got := GuestbookEntryFromAddRequest(req)
	if got == nil {
		t.Fatal("expected non-nil entry")
	}
	if got.Name != "X" || got.Email != "x@y.z" || got.Message != "hello" {
		t.Errorf("got %+v", got)
	}
}
