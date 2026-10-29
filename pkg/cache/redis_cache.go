package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (c *RedisCache) Get(ctx context.Context, key string, dest any) error {
	val, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrCacheMiss
		}
		return err
	}
	return json.Unmarshal(val, dest)
}

func (c *RedisCache) Set(ctx context.Context, key string, val any, ttl time.Duration, tags ...string) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}

	if len(tags) == 0 {
		return c.client.Set(ctx, key, bytes, ttl).Err()
	}

	pipe := c.client.Pipeline()
	pipe.Set(ctx, key, bytes, ttl)
	for _, tag := range tags {
		tagKey := fmt.Sprintf("tag:%s", tag)
		pipe.SAdd(ctx, tagKey, key)
		pipe.Expire(ctx, tagKey, ttl*2)
	}
	_, err = pipe.Exec(ctx)
	return err
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) Remember(ctx context.Context, key string, ttl time.Duration, dest any, fn func() (any, error)) error {
	return c.RememberWithTags(ctx, key, ttl, nil, dest, fn)
}

func (c *RedisCache) RememberWithTags(ctx context.Context, key string, ttl time.Duration, tags []string, dest any, fn func() (any, error)) error {
	err := c.Get(ctx, key, dest)
	if err == nil {
		return nil
	}

	freshData, err := fn()
	if err != nil {
		return err
	}

	if err := c.Set(ctx, key, freshData, ttl, tags...); err != nil {
		return err
	}

	bytes, err := json.Marshal(freshData)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dest)
}

func (c *RedisCache) InvalidateTag(ctx context.Context, tag string) error {
	tagKey := fmt.Sprintf("tag:%s", tag)
	keys, err := c.client.SMembers(ctx, tagKey).Result()
	if err != nil || len(keys) == 0 {
		return nil
	}

	pipe := c.client.Pipeline()
	for _, k := range keys {
		pipe.Del(ctx, k)
	}
	pipe.Del(ctx, tagKey)
	_, err = pipe.Exec(ctx)
	return err
}

func (c *RedisCache) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
