# 🛠️ Step 02: Core Utility Packages (`pkg/`)

Dokumen ini memuat panduan implementasi paket utilitas dasar di folder `pkg/`: Response Envelope, Centralized Error Handling, Structured Logger dengan Context Trace Handler, dan Validator Wrapper.

---

## 🎯 Tujuan Langkah
1. Mengimplementasikan `pkg/response/` untuk response JSON seragam (*Uniform Envelope*) dan pemetaan error domain ke status HTTP.
2. Mengimplementasikan `pkg/logger/` menggunakan `log/slog` bawaan Go 1.21+ yang otomatis meng-ekstrak `correlation_id` dari `context.Context`.
3. Mengimplementasikan `pkg/validator/` menggunakan `go-playground/validator/v10` untuk validasi struct tags otomatis.

---

## 📂 Struktur File yang Dibuat
```text
pkg/
├── response/
│   ├── response.go       # Success & Error JSON builders
│   ├── envelope.go       # Struct Envelope { Success, Message, Data, Errors, Meta }
│   └── error_mapper.go   # Pemetaan domain error -> HTTP status code
├── logger/
│   ├── logger.go         # Inisialisasi slog (JSON di prod, Text di dev)
│   └── trace_handler.go  # slog.Handler decorator untuk correlation_id
└── validator/
    ├── validator.go      # Wrapper go-playground/validator
    └── translator.go     # Mengubah validation error ke slice of field errors
```

---

## 📝 Detail Spesifikasi Kode

### 1. Response Envelope (`pkg/response/envelope.go`)
```go
package response

type FieldError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

type Meta struct {
    RequestID string `json:"requestId,omitempty"`
    Timestamp string `json:"timestamp"`
    Page      int    `json:"page,omitempty"`
    PerPage   int    `json:"perPage,omitempty"`
    Total     int64  `json:"total,omitempty"`
}

type Envelope struct {
    Success bool         `json:"success"`
    Message string       `json:"message"`
    Data    any          `json:"data,omitempty"`
    Errors  []FieldError `json:"errors,omitempty"`
    Meta    Meta         `json:"meta"`
}
```

### 2. Error Mapper (`pkg/response/error_mapper.go`)
```go
package response

import (
    "errors"
    "net/http"
    "github.com/gofiber/fiber/v3"
)

// Sentinel domain errors
var (
    ErrNotFound     = errors.New("resource not found")
    ErrConflict     = errors.New("resource conflict or duplicate entry")
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")
    ErrValidation   = errors.New("validation error")
)

func MapError(err error) (int, string) {
    switch {
    case errors.Is(err, ErrNotFound):
        return http.StatusNotFound, "Resource not found"
    case errors.Is(err, ErrConflict):
        return http.StatusConflict, "Resource already exists"
    case errors.Is(err, ErrUnauthorized):
        return http.StatusUnauthorized, "Unauthorized access"
    case errors.Is(err, ErrForbidden):
        return http.StatusForbidden, "Access forbidden"
    case errors.Is(err, ErrValidation):
        return http.StatusBadRequest, "Validation failed"
    default:
        return http.StatusInternalServerError, "Internal server error"
    }
}
```

### 3. Trace Context Handler & Typed Key (`pkg/logger/trace_handler.go`)
```go
package logger

import (
    "context"
    "log/slog"
)

type correlationIDKey struct{}

// WithCorrelationID menyuntikkan correlation_id ke context secara type-safe
func WithCorrelationID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, correlationIDKey{}, id)
}

// CorrelationIDFromContext mengekstrak correlation_id dari context
func CorrelationIDFromContext(ctx context.Context) string {
    if id, ok := ctx.Value(correlationIDKey{}).(string); ok {
        return id
    }
    return ""
}

type TraceHandler struct {
    slog.Handler
}

func NewTraceHandler(handler slog.Handler) *TraceHandler {
    return &TraceHandler{Handler: handler}
}

func (h *TraceHandler) Handle(ctx context.Context, r slog.Record) error {
    if traceID := CorrelationIDFromContext(ctx); traceID != "" {
        r.AddAttrs(slog.String("correlation_id", traceID))
    }
    return h.Handler.Handle(ctx, r)
}

// WithAttrs memastikan wrapper TraceHandler tidak hilang saat chaining logger.With(...)
func (h *TraceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
    return &TraceHandler{Handler: h.Handler.WithAttrs(attrs)}
}

// WithGroup memastikan wrapper TraceHandler tidak hilang saat chaining logger.WithGroup(...)
func (h *TraceHandler) WithGroup(name string) slog.Handler {
    return &TraceHandler{Handler: h.Handler.WithGroup(name)}
}
```

---

## ✅ Kriteria Verifikasi
- [ ] Response sukses menghasilkan JSON seragam dengan field `success: true`.
- [ ] Error domain (`ErrNotFound`, `ErrConflict`) dipetakan secara akurat ke HTTP status code 404 dan 409.
- [ ] Pemanggilan `slog.InfoContext(ctx, "test")` otomatis menyematkan `correlation_id` pada output terminal/JSON.
- [ ] Validator berhasil menangkap field kosong yang memiliki tag `validate:"required"`.
