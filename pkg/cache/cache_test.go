package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang-starter-kit/pkg/cache"
)

type userCacheItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestMemoryCache_RememberAndInvalidateTag(t *testing.T) {
	ctx := context.Background()
	c := cache.NewMemoryCache()
	defer c.Close()

	key := "user:1"
	tag := "users"
	fetchCount := 0

	fetcher := func() (any, error) {
		fetchCount++
		return userCacheItem{ID: 1, Name: "Andi"}, nil
	}

	// 1. First call: miss, runs fetcher
	var u1 userCacheItem
	err := c.RememberWithTags(ctx, key, 5*time.Minute, []string{tag}, &u1, fetcher)
	require.NoError(t, err)
	assert.Equal(t, 1, fetchCount)
	assert.Equal(t, "Andi", u1.Name)

	// 2. Second call: hit, does not run fetcher
	var u2 userCacheItem
	err = c.RememberWithTags(ctx, key, 5*time.Minute, []string{tag}, &u2, fetcher)
	require.NoError(t, err)
	assert.Equal(t, 1, fetchCount)
	assert.Equal(t, "Andi", u2.Name)

	// 3. Invalidate tag
	err = c.InvalidateTag(ctx, tag)
	require.NoError(t, err)

	// 4. Third call: miss again because tag invalidated, runs fetcher
	var u3 userCacheItem
	err = c.RememberWithTags(ctx, key, 5*time.Minute, []string{tag}, &u3, fetcher)
	require.NoError(t, err)
	assert.Equal(t, 2, fetchCount)
	assert.Equal(t, "Andi", u3.Name)
}
