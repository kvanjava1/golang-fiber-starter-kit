# 💡 Opsi Solusi 06: Project Layout, Dependency Wiring, Storages & CLI Tooling

> **Status Keputusan:** ✅ **DECIDED: Standar Layout + Storages + Cleanenv + Cobra Make**  
> **Keputusan:** Dipilih **Standard Go Project Layout** dengan direktori `storages/` khusus untuk file database SQLite, log rotasi, dan upload media/gambar. DI menggunakan *Constructor Injection* terpusat di `internal/app/app.go`. Konfigurasi menggunakan `cleanenv`, dan scaffolder CLI menggunakan `cmd/make`.

---

## 🎯 Standar Production-Grade yang Harus Dipenuhi
1. **Dedicated Local Storages:** Direktori terpusat `storages/` di root proyek untuk file dinamis runtime (SQLite database, log files, dan upload user/gambar) dengan aturan `.gitignore` yang aman.
2. **Zero Circular Dependencies:** Struktur paket dirancang sedemikian rupa sehingga tidak ada siklus impor paket (*import cycle*), yang merupakan kesalahan kompilasi fatal di Go.
3. **Clear Boundary: `internal/` vs `pkg/`:**
   - `internal/`: Berisi kode logika bisnis privat yang tidak boleh diimpor oleh proyek eksternal.
   - `pkg/`: Berisi pustaka teknis murni yang *reusable* dan agnostik domain (mail, cache, logger, response).
4. **Fail-Fast Configuration:** Aplikasi menolak menyala jika variabel lingkungan wajib (*required .env*) tidak terisi atau salah format.
5. **Ergonomic CLI Generator:** Perintah CLI untuk men-generate boilerplate kode baru secara instan dan type-safe.

---

## 🏛️ 1. Standar Struktur Direktori (Production Layout)

```text
golang-starter-kit/
├── cmd/                          # Entry point aplikasi yang dapat di-build menjadi binary
│   ├── server/
│   │   └── main.go               # Entry point HTTP Web (SSR) & REST API Server (Go-Fiber v3)
│   ├── worker/
│   │   └── main.go               # Entry point Asynq Background Worker
│   └── make/
│       └── main.go               # CLI Scaffolder tool (pengganti bun run make)
│
├── internal/                     # Private application logic (anti-bocor)
│   ├── app/                      # Composition Root (Wiring dependensi & bootstrap)
│   │   ├── app.go                # Server instance & graceful shutdown
│   │   └── router.go             # Master Router (Delegasi RegisterRoutes per domain)
│   ├── config/                   # Struct konfigurasi .env & validasi (cleanenv)
│   ├── middleware/               # Middleware Fiber (auth, trace, logger, ratelimit)
│   └── domain/                   # Modul Domain Bisnis (Gaya A: Co-location per Domain)
│       └── example/              # Contoh domain referensi
│           ├── route.go          # 🚦 Self-Registering Routes: RegisterRoutes(web, api)
│           ├── controller_web.go # 🖥️ Web Controller: render templ HTML untuk browser
│           ├── controller_api.go # 🔌 API Controller: return JSON untuk client
│           ├── service.go        # 🧠 Pure domain business logic
│           ├── repository.go     # 🗄️ Interface & implementasi query (uptrace/bun)
│           ├── model.go          # 📦 Entity database (bun model)
│           └── dto.go            # 📋 Request / Response structs & validasi
│
├── pkg/                          # Technical packages (Re-usable & agnostik domain)
│   ├── mail/                     # Mailable interface & drivers (log, smtp, resend)
│   ├── cache/                    # Redis wrapper (go-redis) & tag-based invalidation
│   ├── queue/                    # Asynq queue dispatcher helper
│   ├── response/                 # Standard uniform JSON response & error mapper
│   ├── logger/                   # slog logger helper dengan context trace
│   └── validator/                # Struct validation wrapper (go-playground/validator)
│
├── resources/                    # Non-Go assets
│   ├── views/                    # Template HTML web & email (a-h/templ)
│   └── migrations/               # File SQL migrasi database (pressly/goose)
│
├── database/                     # 📦 Database Seeding
│   └── seeders/                  # Script seeder data awal (sqlite, pg, mysql, mongo)
│
├── storages/                     # 📂 Local Storage Runtime (diabaikan oleh git kecuali .gitkeep)
│   ├── database/                 # File database lokal (SQLite: sqlite.dev.db, wal, shm)
│   ├── logs/                     # File log harian/rotasi (app/ & cli/)
│   └── uploads/                  # File statis upload gambar / media (public asset)
│
├── Makefile                      # Developer shortcut (make dev, make build, make test)
├── .env.example
├── go.mod
└── go.sum
```

---

## 💾 2. Manajemen Direktori `storages/`

Sama seperti starter kit eksisting dan filosofi Laravel `storage/`, folder `storages/` dialokasikan untuk data lokal runtime:

1. **`storages/database/`**:
   - Menampung file SQLite lokal: `sqlite.dev.db`.
   - DSN di file `.env`: `DATABASE_URL="file:storages/database/sqlite.dev.db?cache=shared"`.
2. **`storages/logs/`**:
   - Menampung log file jika mode logging file diaktifkan (`storages/logs/app/2026-10-29.log`).
   - Dapat diatur rotasi harian menggunakan pustaka seperti `natefinch/lumberjack`.
3. **`storages/uploads/`**:
   - Menampung hasil upload gambar, dokumen, atau avatar pengguna.
   - Di Fiber v3, folder ini dapat di-expose langsung sebagai static route publik:
     ```go
     app.Static("/uploads", "./storages/uploads")
     ```
4. **Git Rule (`.gitignore`):**
   ```gitignore
   storages/database/*.db*
   storages/logs/*/*.log
   storages/uploads/*
   !storages/**/.gitkeep
   ```
   Seluruh isi runtime diabaikan git, tetapi struktur folder tetap tersimpan rapi berkat file `.gitkeep`.

---

## 🔄 2. Pola Dependency Injection Bersih (Composition Root)

Untuk menghindari global variable (`var DB = ...`) yang menyulitkan testing, kita menggunakan **Constructor Injection** yang di-wire pada satu titik terpusat (*Composition Root*) di `internal/app/app.go`:

```go
package app

func BuildContainer(cfg *config.Config) (*Container, error) {
    // 1. Inisialisasi Technical Infrastructure
    db, err := database.New(cfg.DatabaseURL)
    if err != nil {
        return nil, err
    }
    rdb := redis.New(cfg.RedisURL)
    cacheInstance := cache.NewRedisCache(rdb)
    mailClient := mail.New(cfg.MailDriver, cfg.MailConfig)

    // 2. Inisialisasi Repositories
    crudRepo := repository.NewCrudRepository(db)

    // 3. Inisialisasi Services (Menginjeksi Repo, Cache, Mail)
    crudService := example.NewCrudService(crudRepo, cacheInstance, mailClient)

    // 4. Inisialisasi Controllers (Menginjeksi Service)
    crudController := example.NewCrudController(crudService)

    return &Container{
        CrudController: crudController,
        // ...
    }, nil
}
```

---

## ⚙️ 3. Konfigurasi `cleanenv` (Fail-Fast Startup)

Menggunakan pustaka `ilyakaznacheev/cleanenv`:
```go
package config

import (
    "github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
    AppEnv      string `env:"APP_ENV" env-default:"development"`
    Port        int    `env:"PORT" env-default:"3000"`
    DatabaseURL string `env:"DATABASE_URL" env-required:"true"`
    RedisURL    string `env:"REDIS_URL" env-default:"redis://localhost:6379"`
    MailDriver  string `env:"MAIL_DRIVER" env-default:"log"`
}

func Load() (*Config, error) {
    var cfg Config
    if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
        // Read dari environment system jika .env tidak ada
        _ = cleanenv.ReadEnv(&cfg)
    }
    return &cfg, nil
}
```

---

## 🛠️ 4. CLI Scaffolding Tool (`cmd/make/main.go`)

CLI generator dibangun menggunakan pustaka standar Go atau Cobra:

```bash
go run cmd/make/main.go make:module <domain>      # Scaffold 1 modul lengkap (controller, service, repo, dto)
go run cmd/make/main.go make:controller <domain>  # Scaffold Controller
go run cmd/make/main.go make:service <domain>     # Scaffold Service
go run cmd/make/main.go make:repository <domain>  # Scaffold Repository
go run cmd/make/main.go make:mail <name>          # Scaffold Mailable struct + email template
go run cmd/make/main.go make:job <name>           # Scaffold Background Job task
go run cmd/make/main.go make:migration <name>     # Membuat file migrasi SQL baru
go run cmd/make/main.go make:seeder <target> <name> # Membuat file seeder baru (sqlite/pg/mysql/mongo)
```

---

## 📋 5. Makefile Ergonomics

Menyediakan kenyamanan kerja layaknya Bun:
```makefile
.PHONY: dev build test migrate lint

dev:
	go run cmd/server/main.go

worker:
	go run cmd/worker/main.go

test:
	go test -v -race ./...

migrate:
	go run cmd/make/main.go migrate:up
```

---

## 📋 Rekomendasi Solusi & Bahan Diskusi:

> 💬 **Bahan Diskusi Topik 6:**
> Apakah tata letak folder di atas sudah sesuai dengan visi arsitektur Anda untuk starter kit Go ini?
