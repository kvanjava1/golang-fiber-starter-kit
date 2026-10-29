# 📚 Dokumentasi Resmi Go Fiber Starter Kit

Selamat datang di dokumentasi resmi **Go Fiber Starter Kit**. Starter kit ini dirancang sebagai fondasi backend Go modern berstandar enterprise yang menggabungkan performa tinggi, kemudahan pengembangan (*Developer Experience*), arsitektur modular (*Domain-Driven Clean Architecture*), serta fleksibilitas melayani **REST API murni** maupun **Web Server-Side Rendered (SSR)**.

---

## 🗺️ Peta Navigasi Dokumentasi

Dokumentasi ini dibagi ke dalam 7 topik tematik yang dapat Anda pelajari sesuai kebutuhan:

### [01. Getting Started](./01-getting-started/)
Panduan dasar untuk memulai proyek dari nol hingga berjalan di komputer lokal Anda:
- [**01. Overview & Fitur**](./01-getting-started/overview.md) — Filosofi desain, fitur unggulan, dan perbandingan arsitektur.
- [**02. Instalasi & Setup**](./01-getting-started/installation.md) — Prasyarat, setup environment (`.env`), dan inisialisasi awal.
- [**03. Perintah CLI & Makefile**](./01-getting-started/commands-and-scripts.md) — Panduan lengkap perintah `make` dan generator `cmd/make`.

### [02. Architecture & Core Utilities](./02-architecture-and-core/)
Fondasi teknis dan utilitas inti di dalam paket `internal/` dan `pkg/`:
- [**01. Struktur Folder (Directory Layout)**](./02-architecture-and-core/directory-layout.md) — Konvensi Standard Go Project Layout yang digunakan.
- [**02. Manajemen Konfigurasi (cleanenv)**](./02-architecture-and-core/configuration.md) — Type-safe configuration, validasi env, dan default values.
- [**03. Standard Response & Error Handling**](./02-architecture-and-core/response-and-error-handling.md) — Format envelope standar, custom error, dan Global Error Handler.
- [**04. Logging & Observability**](./02-architecture-and-core/logging-and-observability.md) — `log/slog`, Trace Correlation ID, Prometheus metrics, dan Health checks.

### [03. Database, Migrations & Seeders](./03-database-and-migrations/)
Pengelolaan data relasional dan NoSQL berbasis `uptrace/bun`:
- [**01. Multi-Database Setup**](./03-database-and-migrations/multi-database-setup.md) — Dukungan SQLite (pure Go), PostgreSQL, MySQL, dan MongoDB v2.
- [**02. Migrasi Database (Goose)**](./03-database-and-migrations/migrations.md) — Pembuatan dan eksekusi file migrasi SQL.
- [**03. Transaksional Database Seeders**](./03-database-and-migrations/seeders.md) — Menjalankan seeder transaksional dengan rollback otomatis jika terjadi kegagalan.

### [04. Routing & Web UI](./04-routing-and-web/)
Pengelolaan lalu lintas HTTP dan render tampilan:
- [**01. Routing & Middlewares**](./04-routing-and-web/routing-and-middleware.md) — Master Router delegation, CORS, Rate Limiter, Logger, dan Trace correlation.
- [**02. Server-Side Rendering (a-h/templ)**](./04-routing-and-web/web-ssr-and-templ.md) — Type-safe HTML templates yang dikompilasi ke native Go code.
- [**03. Tailwind CSS On-Demand**](./04-routing-and-web/tailwind-css.md) — Integrasi standalone Tailwind CLI tanpa Node.js/npm.

### [05. Background Jobs & Mail System](./05-background-jobs-and-mail/)
Pemrosesan asinkron dan komunikasi email:
- [**01. Asynq Queue & Worker**](./05-background-jobs-and-mail/queue-and-worker.md) — Antrean job berbasis Redis dan proses daemon mandiri (`cmd/worker`).
- [**02. Fluent Mail System**](./05-background-jobs-and-mail/mail-system.md) — Pengiriman email sinkron/antrean dengan driver Log, SMTP, dan Resend.

### [06. Panduan Modul Domain & Scaffolding](./06-domain-and-crud-guide/)
Pola pengembangan fitur domain baru yang rapi dan konsisten:
- [**01. Domain Module Pattern**](./06-domain-and-crud-guide/domain-module-pattern.md) — Anatomi struktur modul (Model, DTO, Repository, Service, Controllers, Route).
- [**02. Scaffolding Modul Baru**](./06-domain-and-crud-guide/scaffolding-new-module.md) — Panduan instan membuat fitur baru dengan perintah `make:module`.

### [07. Pengujian & Deployment Produksi](./07-testing-and-deployment/)
Menjaga kualitas kode dan kesiapan rilis produksi:
- [**01. Strategi Pengujian (Testing)**](./07-testing-and-deployment/testing-strategy.md) — Unit test, integration test, SQLite in-memory, dan race detector.
- [**02. Build & Deployment Produksi**](./07-testing-and-deployment/production-build.md) — Kompilasi single binary, Docker containerization, dan CI/CD GitHub Actions.

---

## ⚡ Quick Start (Jalankan dalam 1 Menit)

```bash
# 1. Salin konfigurasi environment
cp .env.example .env

# 2. Jalankan migrasi database SQLite lokal
make migrate

# 3. Jalankan HTTP Server (default port: 3000)
make dev
```

Buka peramban (*browser*) Anda di:
- **Web App (SSR):** [http://localhost:3000/example/cruds](http://localhost:3000/example/cruds)
- **API Health Check:** [http://localhost:3000/api/v1/health](http://localhost:3000/api/v1/health)
- **Prometheus Metrics:** [http://localhost:3000/metrics](http://localhost:3000/metrics)
