# 🛠️ Step 05: Redis Caching, Prometheus Telemetry & Health Checks

Dokumen ini memuat panduan implementasi paket caching terdistribusi berbasis **Redis** (`go-redis/v9`) dengan dukungan *Tag-based Invalidation*, metrik Prometheus (`client_golang`), dan endpoint health check probes.

---

## 🎯 Tujuan Langkah
1. Mengimplementasikan `pkg/cache/` dengan operasi standar `Get`, `Set`, `Delete`, `Remember`, dan `InvalidateTag`.
2. Menyediakan endpoint Prometheus `/metrics` dengan histogram durasi HTTP dan counter request.
3. Menyediakan endpoint `/api/v1/health` dengan diagnostik mendalam (cek koneksi DB dan Redis).

---

## 📂 Struktur File yang Dibuat
```text
pkg/
└── cache/
    ├── cache.go          # Interface Cache & inisialisasi Redis client
    ├── redis_cache.go    # Implementasi Redis dengan tag-based invalidation
    └── memory_cache.go   # Fallback in-memory cache untuk testing/dev
internal/
├── middleware/
│   └── metrics.go        # Middleware pencatat latensi Prometheus
└── domain/
    └── health/
        ├── controller.go # Handler GET /api/v1/health
        └── route.go
```

---

## 📝 Detail Spesifikasi Kode

### 1. Tag-Based Invalidation di Redis (`pkg/cache/redis_cache.go`)
```go
package cache

import (
    "context"
    "encoding/json"
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

// Remember (Cache-Aside pattern standar tanpa tag)
func (c *RedisCache) Remember(ctx context.Context, key string, ttl time.Duration, dest any, fn func() (any, error)) error {
    return c.RememberWithTags(ctx, key, ttl, nil, dest, fn)
}

// RememberWithTags (Cache-Aside pattern dengan asosiasi Tag Invalidation)
func (c *RedisCache) RememberWithTags(ctx context.Context, key string, ttl time.Duration, tags []string, dest any, fn func() (any, error)) error {
    val, err := c.client.Get(ctx, key).Bytes()
    if err == nil {
        return json.Unmarshal(val, dest)
    }

    freshData, err := fn()
    if err != nil {
        return err
    }

    bytes, err := json.Marshal(freshData)
    if err != nil {
        return err
    }

    // Gunakan pipeline untuk atomic set cache + daftarkan key ke setiap tag set
    pipe := c.client.Pipeline()
    pipe.Set(ctx, key, bytes, ttl)
    for _, tag := range tags {
        tagKey := fmt.Sprintf("tag:%s", tag)
        pipe.SAdd(ctx, tagKey, key)
        // Set TTL pada tagKey agar tidak menjadi memory leak permanen
        pipe.Expire(ctx, tagKey, ttl*2)
    }
    _, _ = pipe.Exec(ctx)

    return json.Unmarshal(bytes, dest)
}

// InvalidateTag: Menghapus semua cache yang terikat dengan tag secara atomic
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
```

### 2. Prometheus Metrics Middleware (`internal/middleware/metrics.go`)
```go
package middleware

import (
    "strconv"
    "time"
    "github.com/gofiber/fiber/v3"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "http_requests_total",
        Help: "Total number of HTTP requests processed",
    }, []string{"method", "path", "status"})

    httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "http_request_duration_seconds",
        Help:    "HTTP request latencies in seconds",
        Buckets: prometheus.DefBuckets,
    }, []string{"method", "path"})
)

func Prometheus() fiber.Handler {
    return func(c fiber.Ctx) error {
        start := time.Now()
        err := c.Next()
        duration := time.Since(start).Seconds()

        status := strconv.Itoa(c.Response().StatusCode())
        path := c.Route().Path // Normalized route path

        httpRequestsTotal.WithLabelValues(c.Method(), path, status).Inc()
        httpRequestDuration.WithLabelValues(c.Method(), path).Observe(duration)

        return err
    }
}
```

---

## ✅ Kriteria Verifikasi
- [ ] Operasi `cache.Remember()` mengambil data dari Redis pada pemanggilan kedua (*cache hit*).
- [ ] `cache.InvalidateTag()` berhasil membersihkan seluruh key dalam grup tag secara atomic.
- [ ] Endpoint `GET /metrics` menyajikan format metrik standar teks Prometheus.
- [ ] Endpoint `GET /api/v1/health` mengembalikan HTTP 200 dengan status "healthy" saat DB dan Redis aktif.
