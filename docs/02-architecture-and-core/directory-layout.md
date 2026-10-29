# 01. Struktur Folder (Directory Layout)

Proyek ini menerapkan **Standard Go Project Layout** yang dimodifikasi untuk mengakomodasi kebutuhan Domain-Driven Design, pemrosesan antrean mandiri, serta Server-Side Rendering.

---

## 📂 Struktur Pohon Direktori

```text
.
├── cmd/                          # Entry point aplikasi yang dapat dieksekusi (main binaries)
│   ├── make/                     # CLI Tool: generator modul, seeder, migrator, & setup tailwind
│   │   ├── commands/             # Logika eksekusi perintah CLI
│   │   ├── stubs/                # Template kode generator (*.stub)
│   │   └── main.go               # Entry point CLI make
│   ├── server/                   # Entry point HTTP Web & REST API Server
│   │   └── main.go
│   └── worker/                   # Entry point Asynq Background Queue Worker Daemon
│       └── main.go
│
├── database/                     # Aset seeding database
│   └── seeders/                  # Runner & implementasi seeder transaksional
│
├── internal/                     # Kode internal aplikasi (private ke modul ini)
│   ├── app/                      # Bootstrap aplikasi, server lifecycle, routing & katalog domains.go
│   │   ├── app.go                # BuildApp container & Start graceful shutdown
│   │   ├── domains.go            # Katalog pendaftaran seluruh modul domain terpusat
│   │   ├── router.go             # Master router Fiber, static assets, & metrics
│   │   └── server.go             # Konfigurasi instance Fiber v3 engine
│   ├── config/                   # Deklarasi struct konfigurasi type-safe (cleanenv)
│   ├── domain/                   # Modul-modul bisnis aplikasi (Clean Architecture)
│   │   ├── example/              # Modul referensi CRUD hidup (model, dto, repo, svc, ctrl, route)
│   │   ├── health/               # Modul pemantauan kesehatan koneksi sistem
│   │   └── user/                 # Modul domain pengguna (hasil scaffolding make:module)
│   └── middleware/               # HTTP middlewares (Trace, Logger, RateLimiter, ErrorHandler)
│
├── pkg/                          # Paket utilitas reusable (dapat digunakan lintas modul/proyek)
│   ├── cache/                    # Abstraksi cache (Redis + in-memory fallback, tag invalidation)
│   ├── database/                 # Multi-DB manager (SQLite, Postgres, MySQL, Mongo v2)
│   ├── logger/                   # Structured logger (log/slog) dengan Trace Correlation ID
│   ├── mail/                     # Fluent Mailer client, Mailables, & multi-driver (Log, SMTP, Resend)
│   ├── queue/                    # Redis Asynq task dispatcher & handlers
│   ├── response/                 # Standard uniform JSON envelope & error mapper
│   └── validator/                # Wrapper go-playground/validator v10
│
├── public/                       # Aset publik statis (disajikan via app.Use("/public", static))
│   ├── css/                      # File CSS hasil build Tailwind (output.css)
│   └── images/                   # Gambar atau logo statis
│
├── resources/                    # Sumber daya mentah (views template, migrations, & CSS)
│   ├── css/                      # File input Tailwind CSS (input.css)
│   ├── migrations/               # File migrasi SQL Goose per database target (sqlite1, pg1, mysql1)
│   └── views/                    # Template HTML a-h/templ
│       ├── emails/               # Template HTML email (welcome, dll.)
│       └── web/                  # Halaman web SSR (CRUD example, layouts)
│
├── storages/                     # Penyimpanan file lokal persisten (di-ignore oleh git)
│   ├── database/                 # Berkas database SQLite (*.db)
│   ├── logs/                     # Berkas log aplikasi lokal (*.log)
│   └── uploads/                  # File upload dari pengguna
│
├── tests/                        # Pengujian integrasi & automated end-to-end
│   ├── helpers/                  # Test harness & pembantu pembuatan in-memory server
│   └── integration/              # Skenario pengujian integrasi CRUD & health checks
│
├── .env.example                  # Template variabel lingkungan
├── Makefile                      # Kumpulan resep otomasi tugas pengembangan
├── go.mod                        # Manifest modul Go & dependensi
└── tailwind.config.js            # Konfigurasi Tailwind CSS standalone
```

---

## 🎯 Perbedaan Direktori `cmd/`, `internal/`, dan `pkg/`

1. **`cmd/` (Application Binaries):**
   Setiap subfolder di dalam `cmd/` mewakili sebuah file binary mandiri. Proyek ini memiliki 3 binary:
   - `cmd/server/main.go`: Melayani request HTTP dari peramban dan klien API.
   - `cmd/worker/main.go`: Menjalankan daemon antrean job latar belakang.
   - `cmd/make/main.go`: Alat bantu pengembang di terminal.

2. **`internal/` (Private Application Code):**
   Direktori yang dilindungi oleh Go compiler sehingga kode di dalamnya tidak bisa di-import oleh repositori luar. Berisi logika domain bisnis, router, konfigurasi, dan middleware aplikasi.

3. **`pkg/` (Shared Core Libraries):**
   Kumpulan pustaka independen yang tidak terikat dengan logika bisnis spesifik (misal: logger, database connection manager, validator, email client). Paket-paket di sini dirancang agar mudah digunakan kembali di proyek lain jika diperlukan.
