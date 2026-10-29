# 03. Format Response Standar & Penanganan Error

Paket [`pkg/response`](file:///home/melodelavic/Documents/golang/golang-fiber-starter-kit/pkg/response/response.go) dan middleware [`internal/middleware/error.go`](file:///home/melodelavic/Documents/golang/golang-fiber-starter-kit/internal/middleware/error.go) menjamin setiap respons REST API memiliki format seragam (*uniform envelope*).

---

## 📦 Struktur JSON Envelope Standar

Seluruh endpoint REST API di bawah prefix `/api/v1/*` menghasilkan format konsisten:

### 1. Respons Berhasil (Success)
```json
{
  "success": true,
  "message": "Data berhasil diambil",
  "data": {
    "id": "e98e4d3a-23b9-4c5c-9c71-3c66f5bb392c",
    "title": "Belajar Go Fiber v3",
    "status": "active"
  },
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 45
  }
}
```

### 2. Respons Gagal Validasi (Unprocessable Entity - 422)
```json
{
  "success": false,
  "message": "Validasi input gagal",
  "errors": {
    "title": "Field title wajib diisi",
    "email": "Format email tidak valid"
  }
}
```

### 3. Respons Error Sistem (Internal Server Error - 500)
```json
{
  "success": false,
  "message": "Terjadi kesalahan internal pada server",
  "error": "Internal Server Error"
}
```

---

## 💻 Penggunaan di Controller API

Gunakan fungsi pembantu dari paket `pkg/response`:

```go
package example

import (
    "github.com/gofiber/fiber/v3"
    "golang-starter-kit/pkg/response"
)

func (c *APIController) GetByID(ctx fiber.Ctx) error {
    id := ctx.Params("id")
    item, err := c.service.FindByID(ctx.Context(), id)
    if err != nil {
        return response.NotFound(ctx, "Data tidak ditemukan")
    }

    return response.OK(ctx, "Detail item berhasil diambil", item)
}

func (c *APIController) Create(ctx fiber.Ctx) error {
    var req CreateRequest
    if err := ctx.Bind().Body(&req); err != nil {
        return response.BadRequest(ctx, "Payload JSON tidak valid", err.Error())
    }

    // Validasi field
    if errMap := c.validator.Validate(req); len(errMap) > 0 {
        return response.UnprocessableEntity(ctx, "Validasi gagal", errMap)
    }

    created, err := c.service.Create(ctx.Context(), req)
    if err != nil {
        return response.InternalServerError(ctx, "Gagal menyimpan item", err.Error())
    }

    return response.Created(ctx, "Item berhasil dibuat", created)
}
```

---

## 🛡️ Global Error Handler Middleware

Middleware [`GlobalErrorHandler`](file:///home/melodelavic/Documents/golang/golang-fiber-starter-kit/internal/middleware/error.go) menangani:
1. **Panic Recovery:** Menangkap kepanikan (*panic*) runtime tanpa membuat server crash, mencatat stack trace ke logger, dan mengembalikan status code 500.
2. **Context Distinction:**
   - Jika request berasal dari API (`/api/*`), middleware merespons dengan JSON envelope standar.
   - Jika request berasal dari peramban Web SSR, middleware menyajikan halaman HTML ramah pengguna.
