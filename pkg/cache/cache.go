package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrCacheMiss = errors.New("cache: key not found")
)

type Cache interface {
	Get(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, key string, val any, ttl time.Duration, tags ...string) error
	Delete(ctx context.Context, key string) error
	Remember(ctx context.Context, key string, ttl time.Duration, dest any, fn func() (any, error)) error
	RememberWithTags(ctx context.Context, key string, ttl time.Duration, tags []string, dest any, fn func() (any, error)) error
	InvalidateTag(ctx context.Context, tag string) error
	Close() error
}

// NewRedisClient menginisialisasi client Redis dari URI
func NewRedisClient(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis url '%s': %w", redisURL, err)
	}

	client := redis.NewClient(opts)
	return client, nil
}
