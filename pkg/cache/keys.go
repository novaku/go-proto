package cache

import "fmt"

// ListEntriesCacheKey builds a deterministic key for paginated guestbook list caching.
func ListEntriesCacheKey(limit, offset int32) string {
	return "listentries:" + fmt.Sprintf("limit=%d&offset=%d", limit, offset)
}
