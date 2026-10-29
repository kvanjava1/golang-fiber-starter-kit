# 🛠️ Step 04: HTTP Fiber v3 Engine, Middlewares & Master Router

Dokumen ini memuat panduan implementasi HTTP web server berbasis **Go-Fiber v3**, pipeline middleware (*RequestID, Context Tracing, Logger, CORS, RateLimiter*), dan orkestrasi master router dengan pola *Self-Registering Routes*.

---

## 🎯 Tujuan Langkah
1. Menginisialisasi instance Fiber v3 dengan custom error handler terpusat.
2. Membangun pipeline middleware global (Context Tracing, HTTP Logger, RateLimiter, Security Headers).
3. Menerapkan master router di `internal/app/router.go` yang mendelegasikan group `/` (Web) dan `/api/v1` (API).
4. Menyediakan endpoint static asset untuk melayani folder `storages/uploads/` dan `public/`.

---

## 📂 Struktur File yang Dibuat
```text
internal/
├── app/
│   ├── app.go            # Composition Root & Lifecycle Runner
│   ├── server.go         # Konfigurasi Fiber App
│   └── router.go         # Master Router Mount
└── middleware/
    ├── trace.go          # Correlation ID injector
    ├── logger.go         # HTTP request completion logger
    ├── rate_limit.go     # Fiber limiter wrapper
    └── error_handler.go  # Global Fiber error handler
```

---

## 📝 Detail Spesifikasi Kode

### 1. Trace Context Middleware (`internal/middleware/trace.go`)
```go
package middleware

import (
    "github.com/gofiber/fiber/v3"
    "github.com/google/uuid"
    "golang-starter-kit/pkg/logger"
)

const HeaderCorrelationID = "X-Correlation-ID"

func Trace() fiber.Handler {
    return func(c fiber.Ctx) error {
        traceID := c.Get(HeaderCorrelationID)
        if traceID == "" {
            traceID = uuid.New().String()
        }

        // Simpan ke Response Header agar client bisa melacak issue
        c.Set(HeaderCorrelationID, traceID)

        // Inject ke context.Context request secara type-safe via pkg/logger
        ctx := logger.WithCorrelationID(c.Context(), traceID)
        c.SetUserContext(ctx)

        return c.Next()
    }
}
```

### 2. HTTP Request Logger Middleware (`internal/middleware/logger.go`)
```go
package middleware

import (
    "log/slog"
    "time"
    "github.com/gofiber/fiber/v3"
)

func Logger() fiber.Handler {
    return func(c fiber.Ctx) error {
        start := time.Now()
        err := c.Next()
        duration := time.Since(start)

        slog.InfoContext(c.Context(), "HTTP Request",
            "method", c.Method(),
            "path", c.Path(),
            "status", c.Response().StatusCode(),
            "duration_ms", duration.Milliseconds(),
            "ip", c.IP(),
        )

        return err
    }
}
```

### 3. Master Router (`internal/app/router.go`)
```go
package app

import (
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/cors"
    "github.com/gofiber/fiber/v3/middleware/static"
)

func RegisterRoutes(app *fiber.App, container *Container) {
    // 📁 Static Assets
    app.Use("/uploads", static.New("./storages/uploads"))
    app.Use("/public", static.New("./public"))

    // 🌐 Web SSR Routes Group (/)
    web := app.Group("/")
    web.Get("/", container.HomeController.Index)

    // 🔌 REST API Routes Group (/api/v1)
    api := app.Group("/api/v1", cors.New())
    api.Get("/health", container.HealthController.Check)

    // 🚀 Delegasi Domain Modules
    container.ExampleDomain.RegisterRoutes(web, api)
}
```

---

## ✅ Kriteria Verifikasi
- [ ] Setiap HTTP request otomatis memiliki header `X-Correlation-ID`.
- [ ] Log HTTP mencatat method, path, status, latency dalam ms, dan correlation ID yang sama.
- [ ] Route statis `/uploads/avatar.png` berhasil menyajikan file dari `storages/uploads/`.
- [ ] Grouping route Web dan API terpisah dan beroperasi pada satu port yang sama.
