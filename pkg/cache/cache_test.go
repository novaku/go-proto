package cache

import (
	"testing"
)

func TestListEntriesCacheKey(t *testing.T) {
	tests := []struct {
		name     string
		limit    int32
		offset   int32
		expected string
	}{
		{
			name:     "basic key generation",
			limit:    10,
			offset:   0,
			expected: "listentries:limit=10&offset=0",
		},
		{
			name:     "with offset",
			limit:    20,
			offset:   40,
			expected: "listentries:limit=20&offset=40",
		},
		{
			name:     "zero values",
			limit:    0,
			offset:   0,
			expected: "listentries:limit=0&offset=0",
		},
		{
			name:     "large values",
			limit:    1000,
			offset:   5000,
			expected: "listentries:limit=1000&offset=5000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ListEntriesCacheKey(tt.limit, tt.offset)
			if result != tt.expected {
				t.Errorf("ListEntriesCacheKey(%d, %d) = %s; want %s",
					tt.limit, tt.offset, result, tt.expected)
			}
		})
	}
}
