package cache

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

type memoryItem struct {
	data      []byte
	expiresAt time.Time
}

type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]memoryItem
	tags  map[string]map[string]struct{}
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		items: make(map[string]memoryItem),
		tags:  make(map[string]map[string]struct{}),
	}
}

func (m *MemoryCache) Get(ctx context.Context, key string, dest any) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.items[key]
	if !exists {
		return ErrCacheMiss
	}

	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		return ErrCacheMiss
	}

	return json.Unmarshal(item.data, dest)
}

func (m *MemoryCache) Set(ctx context.Context, key string, val any, ttl time.Duration, tags ...string) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	var exp time.Time
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}

	m.items[key] = memoryItem{
		data:      bytes,
		expiresAt: exp,
	}

	for _, t := range tags {
		if _, exists := m.tags[t]; !exists {
			m.tags[t] = make(map[string]struct{})
		}
		m.tags[t][key] = struct{}{}
	}

	return nil
}

func (m *MemoryCache) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.items, key)
	return nil
}

func (m *MemoryCache) Remember(ctx context.Context, key string, ttl time.Duration, dest any, fn func() (any, error)) error {
	return m.RememberWithTags(ctx, key, ttl, nil, dest, fn)
}

func (m *MemoryCache) RememberWithTags(ctx context.Context, key string, ttl time.Duration, tags []string, dest any, fn func() (any, error)) error {
	err := m.Get(ctx, key, dest)
	if err == nil {
		return nil
	}

	fresh, err := fn()
	if err != nil {
		return err
	}

	if err := m.Set(ctx, key, fresh, ttl, tags...); err != nil {
		return err
	}

	bytes, err := json.Marshal(fresh)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dest)
}

func (m *MemoryCache) InvalidateTag(ctx context.Context, tag string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	keys, exists := m.tags[tag]
	if !exists {
		return nil
	}

	for k := range keys {
		delete(m.items, k)
	}
	delete(m.tags, tag)
	return nil
}

func (m *MemoryCache) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items = make(map[string]memoryItem)
	m.tags = make(map[string]map[string]struct{})
	return nil
}
