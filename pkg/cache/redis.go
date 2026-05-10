package cache

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements Cache using Redis and a pluggable ValueCodec.
type RedisCache struct {
	client *redis.Client
	codec  ValueCodec
}

// NewRedisCache builds a Redis-backed cache using JSON encoding.
func NewRedisCache(client *redis.Client) *RedisCache {
	return NewRedisCacheWithCodec(client, JSONCodec{})
}

// NewRedisCacheWithCodec allows injecting a custom codec (Open/Closed: extend without modifying Redis logic).
func NewRedisCacheWithCodec(client *redis.Client, codec ValueCodec) *RedisCache {
	if codec == nil {
		codec = JSONCodec{}
	}
	return &RedisCache{
		client: client,
		codec:  codec,
	}
}

// Set stores a value with TTL.
func (rc *RedisCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	data, err := rc.codec.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := rc.client.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// Get loads a value if present.
func (rc *RedisCache) Get(ctx context.Context, key string, dest any) (bool, error) {
	val, err := rc.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get cache: %w", err)
	}

	if err := rc.codec.Unmarshal([]byte(val), dest); err != nil {
		return false, fmt.Errorf("failed to unmarshal cache: %w", err)
	}

	return true, nil
}

// Delete removes a key.
func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	if isPatternKey(key) {
		var cursor uint64
		var keys []string

		for {
			matchedKeys, nextCursor, err := rc.client.Scan(ctx, cursor, key, 100).Result()
			if err != nil {
				return fmt.Errorf("failed to scan cache keys: %w", err)
			}

			if len(matchedKeys) > 0 {
				keys = append(keys, matchedKeys...)
			}

			cursor = nextCursor
			if cursor == 0 {
				break
			}
		}

		if len(keys) == 0 {
			return nil
		}

		if err := rc.client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("failed to delete cache by pattern: %w", err)
		}
		return nil
	}

	if err := rc.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete cache: %w", err)
	}
	return nil
}

func isPatternKey(key string) bool {
	return strings.ContainsAny(key, "*?[]")
}

// Close closes the Redis client.
func (rc *RedisCache) Close() error {
	return rc.client.Close()
}
