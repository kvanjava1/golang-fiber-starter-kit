# 💡 Opsi Solusi 01: HTTP Web Framework, OpenAPI & Production Error Flow

> **Status Keputusan:** ✅ **DECIDED: Opsi A (Go-Fiber v3)**  
> **Keputusan:** Dipilih **Go-Fiber v3** untuk mempertahankan ergonomi dan rasa ngoding yang mirip dengan Hono/Express, performa throughput tertinggi (`fasthttp`), serta ekosistem middleware resmi yang sangat lengkap. Dokumentasi OpenAPI menggunakan standar `swag` + `fiber-swagger`.

---

## 🎯 Standar Production-Grade yang Harus Dipenuhi
1. **Idiomatic `context.Context` Propagation:** Setiap HTTP handler mengekstrak `c.Request().Context()` (atau context Fiber) dan menyematkan `X-Request-ID` / `X-Correlation-ID`, lalu meneruskannya ke seluruh layer di bawahnya.
2. **Clean Error Mapping (Tanpa `panic`):**
   - Service mengembalikan error domain murni (misal: `domain.ErrNotFound`, `domain.ErrValidation`).
   - Layer Controller/Presenter memetakan error domain tersebut ke status HTTP dan format JSON seragam (*Uniform Envelope*).
   - Detail error internal (seperti koneksi database putus) dicatat di log privat dengan stack trace, sedangkan client publik hanya menerima pesan aman `Internal Server Error` berserta `requestId`.
3. **Validasi Input Terstandarisasi:** Validasi DTO deklaratif menggunakan struct tags (`validate:"required,email"`), mengembalikan HTTP 400 otomatis dengan rincian field yang gagal.
4. **OpenAPI / Swagger Documentation:** Dokumentasi interaktif yang mudah diakses di browser (`/docs` atau `/api/docs`).

---

## ⚖️ Pilihan Framework HTTP

### 🔹 Opsi A: Go-Fiber v3 (Paling Mirip Hono, Performa Tertinggi)
- **Karakteristik:** Berbasis `fasthttp`, routing super cepat, alokasi memori mendekati nol.
- **Kemiripan dengan Hono:** Sintaks middleware dan handlernya hampir identik dengan Hono:
  ```go
  app.Post("/api/example/cruds", handler.Create)
  ```
- **Kelebihan:** Ekosistem middleware resmi sangat kaya (CORS, RateLimiter, Helmet, Logger, Prometheus).
- **Kekurangan:** Tidak mengimplementasikan `http.Handler` standar secara native (meskipun Fiber v3 menyediakan adapter `adaptor.FiberHandler`).

### 🔹 Opsi B: Echo v4 (Standar Industri Enterprise `net/http`)
- **Karakteristik:** 100% berbasis standar library Go `net/http`, sangat stabil, mature, dan menjadi pilihan utama perusahaan skala besar.
- **Kelebihan:** Kompatibilitas sempurna dengan seluruh ekosistem middleware standar Go, context handling sangat rapi (`echo.Context`), dan penanganan error terpusat yang sangat matang (`echo.HTTPErrorHandler`).
- **Kekurangan:** Performa sedikit di bawah fasthttp pada benchmark sintetis ekstrem (namun lebih dari cukup untuk 99.9% kebutuhan riil: puluhan ribu req/detik).

### 🔹 Opsi C: Huma Framework (`huma.rocks`) di atas Echo / Chi (Otomatis OpenAPI ala `@hono/zod-openapi`)
- **Karakteristik:** Framework modern berbasis Go generics dan struct typing.
- **Kelebihan:** Sama persis seperti `@hono/zod-openapi`: Struct Go langsung menghasilkan validasi input, parser query/body/headers, dan spesifikasi OpenAPI 3.1 tanpa perlu menulis komentar anotasi `swag`.
- **Kekurangan:** Struktur handlernya murni fungsional (`func(ctx, input) (*output, error)`), sedikit berbeda dari gaya controller konvensional.

---

## 🛡️ Standar Pola Penanganan Error Idiomatik Go (*The Go Way*)

Untuk mencegah error handling yang kotor atau memaksakan exception, starter kit akan menerapkan **Pola Sentinel & Domain Error Mapping**:

```mermaid
flowchart LR
    Service["Service Layer"] -->|Returns (nil, domain.ErrNotFound)| Controller["Controller / Handler"]
    Controller -->|Calls| ErrorMapper["pkg/response/error.go"]
    ErrorMapper -->|Checks errors.Is(err, domain.ErrNotFound)| MapperLogic{"Mapping Rule"}
    MapperLogic -->|Matches ErrNotFound| HTTP404["HTTP 404 Not Found"]
    MapperLogic -->|Matches ErrValidation| HTTP400["HTTP 400 Bad Request"]
    MapperLogic -->|Matches ErrConflict| HTTP409["HTTP 409 Conflict"]
    MapperLogic -->|Unknown Error| HTTP500["HTTP 500 Internal Error + Log Trace"]
    HTTP404 --> ClientResponse["JSON Envelope Response"]
```

### 1. Definisi Error Domain di `internal/domain/errors.go`:
```go
package domain

import "errors"

var (
    ErrNotFound     = errors.New("resource not found")
    ErrConflict     = errors.New("resource already exists")
    ErrUnauthorized = errors.New("unauthorized access")
    ErrForbidden    = errors.New("forbidden action")
)

type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}
```

### 2. Standard JSON Response Envelope (`pkg/response/`):
```json
// Sukses
{
  "success": true,
  "message": "Resource created successfully",
  "data": { ... },
  "meta": {
    "requestId": "req_01hzabc123",
    "timestamp": "2026-10-28T16:00:00Z"
  }
}

// Error
{
  "success": false,
  "message": "Validation failed",
  "errors": [
    { "field": "email", "message": "Email must be a valid email address" }
  ],
  "meta": {
    "requestId": "req_01hzabc123",
    "timestamp": "2026-10-28T16:00:00Z"
  }
}
```

---

## 📋 Rekomendasi Solusi & Bahan Diskusi:

| Kriteria | Opsi A (Fiber v3) | Opsi B (Echo v4) | Opsi C (Huma + Echo/Chi) |
|---|---|---|---|
| **Kemiripan Rasa dengan Hono** | ⭐⭐⭐⭐⭐ (Identik) | ⭐⭐⭐⭐ (Mirip) | ⭐⭐⭐ (Pendekatan Baru) |
| **Generasi OpenAPI** | Semi-otomatis (`swag`) | Semi-otomatis (`swag`) | 100% Otomatis (via Struct) |
| **Kepatuhan Standar Go `net/http`** | Cukup (fasthttp) | Sempurna (100% stdlib) | Sempurna (100% stdlib) |
| **Kematangan di Ekosistem** | Sangat Populer | Standar Enterprise | Modern / Naik Daun |

> 💬 **Bahan Diskusi Topik 1:**
> - Apakah Anda lebih memilih **Opsi A (Fiber v3)** untuk menjaga kemiripan rasa maksimal dengan Hono, atau **Opsi B (Echo v4)** untuk kepatuhan penuh standar `net/http` Go?
> - Untuk OpenAPI docs, apakah nyaman menggunakan anotasi `swag` standar Go, atau tertarik menggunakan **Opsi C (Huma)** agar otomatis tanpa anotasi?
