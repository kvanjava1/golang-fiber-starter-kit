# 02. Panduan Scaffolding Modul Baru (`make:module`)

Membuat modul domain baru dari awal secara manual bisa memakan waktu. Starter kit ini menyediakan perintah generator cerdas untuk mempercepat proses pembuatan boilerplate modul secara instan.

---

## ⚡ Langkah Cepat Membuat Fitur Baru (Contoh: `product`)

### Langkah 1: Jalankan CLI Generator
Jalankan perintah berikut di terminal Anda:
```bash
go run cmd/make/main.go make:module product
```

Output di terminal:
```text
🚀 Generating domain module 'product'...
   ✅ Created: internal/domain/product/model.go
   ✅ Created: internal/domain/product/dto.go
   ✅ Created: internal/domain/product/repository.go
   ✅ Created: internal/domain/product/service.go
   ✅ Created: internal/domain/product/controller_api.go
   ✅ Created: internal/domain/product/controller_web.go
   ✅ Created: internal/domain/product/route.go
🎉 Module 'product' generated successfully!
```

---

### Langkah 2: Buat File Migrasi Database
Buat file migrasi untuk tabel `products`:
```bash
go run cmd/make/main.go make:migration create_products_table
```

Buka file migrasi yang dihasilkan di folder `resources/migrations/sqlite1/` dan tambahkan skema SQL:
```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
```

Lalu jalankan migrasinya:
```bash
make migrate
```

---

### Langkah 3: Daftarkan Modul ke Katalog Domain (`internal/app/domains.go`)
Buka file [`internal/app/domains.go`](file:///home/melodelavic/Documents/golang/golang-fiber-starter-kit/internal/app/domains.go) dan tambahkan 1 baris modul baru Anda menggunakan fungsi provider `NewModule()`:

```go
import (
    "golang-starter-kit/internal/domain/product" // Import modul baru
)

func RegisterAppDomains(c *Container) {
    // ...
    sqlDB := c.DB.DefaultSQL()

    // 📦 Katalog Modul (Cukup 1 Baris per Modul!):
    c.RegisterDomain(example.NewModule(sqlDB, c.Cache))
    c.RegisterDomain(user.NewModule(sqlDB, c.Cache))
    c.RegisterDomain(product.NewModule(sqlDB, c.Cache)) // <-- Tambahkan di sini
}
```

---

### Langkah 4: Uji Coba Endpoint
Jalankan server:
```bash
make dev
```

Endpoint Anda kini sudah otomatis aktif:
- **REST API:** `GET /api/v1/products`, `POST /api/v1/products`, `GET /api/v1/products/:id`
- **Web SSR:** `GET /products`, `GET /products/create`
