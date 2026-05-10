package cache

import (
	"context"
	"time"
)

// Cache is the application-level cache port (implemented by Redis or test doubles).
type Cache interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string, dest any) (bool, error)
	Delete(ctx context.Context, key string) error
}
