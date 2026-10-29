# 🛠️ Step 09: Pengujian Otomatis, Test Suite & CI/CD Pipeline

Dokumen ini memuat panduan implementasi strategi pengujian menyeluruh (Unit Test & Integration Test), pembuatan test helper, isolasi state memori, serta konfigurasi pipeline GitHub Actions CI.

---

## 🎯 Tujuan Langkah
1. Membangun struktur pengujian di `tests/unit/` dan `tests/integration/`.
2. Menyediakan test helper (`tests/helpers/`) untuk inisialisasi Fiber test client (`app.Test()`), mock repository, dan in-memory SQLite database (`:memory:`).
3. Menyusun workflow GitHub Actions CI (`.github/workflows/ci.yml`) untuk verifikasi linter (`golangci-lint`) dan test otomatis pada branch utama dan pull request.

---

## 📂 Struktur File yang Dibuat
```text
tests/
├── helpers/
│   ├── app.go            # Helper inisialisasi Fiber app untuk integration test
│   ├── db.go             # Helper database :memory: terisolasi
│   └── mock_repo.go      # Mock repository generator
├── unit/
│   ├── response_test.go
│   ├── mail_driver_test.go
│   └── service_test.go
└── integration/
    ├── crud_test.go
    ├── health_test.go
    └── mail_queue_test.go
.github/
└── workflows/
    └── ci.yml
```

---

## 📝 Detail Spesifikasi Kode

### 1. Integration Test Helper (`tests/helpers/app.go`)
```go
package helpers

import (
    "net/http"
    "github.com/gofiber/fiber/v3"
    "golang-starter-kit/internal/app"
    "golang-starter-kit/internal/config"
    "golang-starter-kit/pkg/database"
)

func SetupTestApp() (*fiber.App, func()) {
    // 1. Buat koneksi in-memory SQLite terisolasi
    db, err := database.ConnectSQLite(":memory:")
    if err != nil {
        panic(err)
    }

    // 2. Jalankan migrasi schema ke database in-memory
    RunTestMigrations(db)

    // 3. Build container dengan mock config
    cfg := &config.Config{
        AppEnv: "testing",
        Port:   3001,
    }
    container, _ := app.BuildContainer(cfg, db)

    // 4. Inisialisasi Fiber
    fiberApp := fiber.New()
    app.RegisterRoutes(fiberApp, container)

    teardown := func() {
        _ = db.Close()
    }

    return fiberApp, teardown
}
```

### 2. Contoh Integration Test (`tests/integration/crud_test.go`)
```go
package integration

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/stretchr/testify/assert"
    "golang-starter-kit/tests/helpers"
)

func TestCRUD_CreateAndGet(t *testing.T) {
    app, teardown := helpers.SetupTestApp()
    defer teardown()

    // 1. Test POST /api/v1/example/cruds
    payload := map[string]string{
        "title":       "Test CRUD",
        "description": "Deskripsi test",
    }
    body, _ := json.Marshal(payload)

    req := httptest.NewRequest(http.MethodPost, "/api/v1/example/cruds", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusCreated, resp.StatusCode)

    // 2. Test GET /api/v1/example/cruds
    reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/example/cruds", nil)
    respGet, err := app.Test(reqGet)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, respGet.StatusCode)
}
```

### 3. GitHub Actions CI (`.github/workflows/ci.yml`)
```yaml
name: CI Test Suite

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true

      - name: Install templ
        run: go install github.com/a-h/templ/cmd/templ@latest

      - name: Generate Templates
        run: templ generate

      - name: Run Tests with Race Detector
        run: go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

      - name: Verify Build
        run: |
          go build -v -o /dev/null cmd/server/main.go
          go build -v -o /dev/null cmd/worker/main.go
          go build -v -o /dev/null cmd/make/main.go
```

---

## ✅ Kriteria Verifikasi
- [ ] Menjalankan `go test -v -race ./...` menghasilkan 100% PASS tanpa race condition.
- [ ] Database test menggunakan `:memory:` sehingga tidak meninggalkan file sisa di disk.
- [ ] CI workflow di GitHub Actions selesai hijau pada setiap push.
