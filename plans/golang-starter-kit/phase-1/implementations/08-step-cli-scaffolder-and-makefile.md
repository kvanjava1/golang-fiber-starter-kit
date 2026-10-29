# 🛠️ Step 08: CLI Scaffolder (`cmd/make`) & Makefile DX

Dokumen ini memuat panduan implementasi tool generator kode berbasis CLI (`cmd/make/main.go`) yang setara dengan `bun run make` di kit Hono, manajemen stub template, serta task runner Makefile.

---

## 🎯 Tujuan Langkah
1. Membangun executable tool `cmd/make/main.go` untuk scaffolding instan.
2. Menyediakan perintah pembuat boilerplate: `make:module`, `make:controller`, `make:service`, `make:repository`, `make:mail`, `make:job`, `make:migration`, dan `make:seeder`.
3. Menyediakan `Makefile` dengan shortcut ergonomis (`make dev`, `make worker`, `make test`, `make migrate`, `make seed`, `make templ`).

---

## 📂 Struktur File yang Dibuat
```text
cmd/
└── make/
    ├── main.go               # Entry point CLI make
    ├── commands/             # Handler per sub-command
    └── stubs/                # Template mentah (*.stub)
        ├── controller.stub
        ├── service.stub
        ├── repository.stub
        ├── model.stub
        ├── route.stub
        ├── mail.stub
        ├── job.stub
        ├── migration.stub
        └── seeder.stub
Makefile
```

---

## 📝 Detail Spesifikasi Kode

### 1. Template Stub Contoh (`cmd/make/stubs/service.stub`)
```text
package {{domain}}

import (
    "context"
    "log/slog"
)

type {{PascalCase}}Service struct {
    repo   {{PascalCase}}Repository
    logger *slog.Logger
}

func New{{PascalCase}}Service(repo {{PascalCase}}Repository, logger *slog.Logger) *{{PascalCase}}Service {
    return &{{PascalCase}}Service{
        repo:   repo,
        logger: logger,
    }
}
```

### 2. Makefile Shortcut (`Makefile`)
```makefile
.PHONY: dev worker test migrate seed templ build clean

# Menjalankan HTTP Web & API Server
dev:
	templ generate
	go run cmd/server/main.go

# Menjalankan Background Job Worker
worker:
	go run cmd/worker/main.go

# Menjalankan Server + Worker dalam 1 proses (untuk kemudahan local dev)
dev-all:
	templ generate
	go run cmd/server/main.go --with-worker

# Menjalankan test suite
test:
	go test -v -race ./...

# Kompilasi template HTML a-h/templ
templ:
	templ generate

# Migrasi Database
migrate:
	go run cmd/make/main.go migrate:up --db=sqlite1

# Seeder Database
seed:
	go run cmd/make/main.go seed

# Build single binary statis
build:
	templ generate
	CGO_ENABLED=1 go build -o bin/server cmd/server/main.go
	CGO_ENABLED=1 go build -o bin/worker cmd/worker/main.go
```

---

## ✅ Kriteria Verifikasi
- [ ] Perintah `go run cmd/make/main.go make:module product` otomatis membuat folder `internal/domain/product/` berisi 6 file lengkap.
- [ ] Perintah `make dev` berhasil meng-generate template `templ` dan menyalakan server.
- [ ] Perintah `make build` menghasilkan file binary mandiri `bin/server` dan `bin/worker`.
