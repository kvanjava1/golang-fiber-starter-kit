# 🚀 Go Fiber Starter Kit

Enterprise-grade Go Starter Kit featuring **Go Fiber v3**, **Bun Multi-DB ORM**, **Asynq Task Queue**, **Fluent Mail System**, **a-h/templ Server-Side Rendering**, **On-Demand Standalone Tailwind CSS**, and **Domain-Driven Clean Architecture**.

---

## ⚡ Quick Start

```bash
# 1. Salin berkas konfigurasi environment
cp .env.example .env

# 2. Jalankan migrasi database awal (SQLite)
make migrate

# 3. Jalankan HTTP Server (port 3000)
make dev
```

Buka peramban di:
- **Web App Example:** [http://localhost:3000/example/cruds](http://localhost:3000/example/cruds)
- **Web App User:** [http://localhost:3000/user](http://localhost:3000/user)
- **API Health Check:** [http://localhost:3000/api/v1/health](http://localhost:3000/api/v1/health)
- **Prometheus Metrics:** [http://localhost:3000/metrics](http://localhost:3000/metrics)

---

## 📚 Dokumentasi Lengkap (`docs/`)

Dokumentasi komprehensif disusun secara tematik di dalam folder [`docs/`](./docs/README.md):

1. **[01. Getting Started](./docs/01-getting-started/)**
   - [Overview & Fitur Unggulan](./docs/01-getting-started/overview.md)
   - [Instalasi & Setup Lingkungan](./docs/01-getting-started/installation.md)
   - [Panduan Perintah CLI & Makefile](./docs/01-getting-started/commands-and-scripts.md)

2. **[02. Architecture & Core Utilities](./docs/02-architecture-and-core/)**
   - [Struktur Folder (Directory Layout)](./docs/02-architecture-and-core/directory-layout.md)
   - [Manajemen Konfigurasi Type-Safe (`cleanenv`)](./docs/02-architecture-and-core/configuration.md)
   - [Format Response Standar & Error Handling](./docs/02-architecture-and-core/response-and-error-handling.md)
   - [Logging (`log/slog`) & Observability](./docs/02-architecture-and-core/logging-and-observability.md)

3. **[03. Database, Migrations & Seeders](./docs/03-database-and-migrations/)**
   - [Multi-Database Setup (`uptrace/bun`)](./docs/03-database-and-migrations/multi-database-setup.md)
   - [Migrasi Database (`goose v3`)](./docs/03-database-and-migrations/migrations.md)
   - [Transaksional Database Seeders](./docs/03-database-and-migrations/seeders.md)

4. **[04. Routing & Web UI](./docs/04-routing-and-web/)**
   - [Routing & Middlewares](./docs/04-routing-and-web/routing-and-middleware.md)
   - [Server-Side Rendering (`a-h/templ`)](./docs/04-routing-and-web/web-ssr-and-templ.md)
   - [Tailwind CSS On-Demand (No Node.js)](./docs/04-routing-and-web/tailwind-css.md)

5. **[05. Background Jobs & Mail System](./docs/05-background-jobs-and-mail/)**
   - [Asynq Queue & Worker Daemon](./docs/05-background-jobs-and-mail/queue-and-worker.md)
   - [Fluent Mail System (Log, SMTP, Resend)](./docs/05-background-jobs-and-mail/mail-system.md)

6. **[06. Panduan Modul Domain & Scaffolding](./docs/06-domain-and-crud-guide/)**
   - [Domain Module Pattern (Clean Architecture)](./docs/06-domain-and-crud-guide/domain-module-pattern.md)
   - [Panduan Scaffolding Modul Baru (`make:module`)](./docs/06-domain-and-crud-guide/scaffolding-new-module.md)

7. **[07. Pengujian & Deployment Produksi](./docs/07-testing-and-deployment/)**
   - [Strategi Pengujian (Testing Strategy)](./docs/07-testing-and-deployment/testing-strategy.md)
   - [Build & Deployment Produksi](./docs/07-testing-and-deployment/production-build.md)

---

## 🛠️ Ringkasan Perintah Penting

| Perintah | Deskripsi |
| :--- | :--- |
| `make dev` | Menjalankan server HTTP (API-First, ringan) di port `3000` |
| `make dev-web` | Menjalankan server Web terpadu (generate `templ` + build CSS) |
| `make worker` | Menjalankan Asynq background worker daemon |
| `make test` | Menjalankan seluruh test suite dengan Race Detector |
| `make migrate` | Menjalankan migrasi database pending |
| `make tailwind-install` | Mengunduh standalone Tailwind CLI tanpa Node.js |
| `make build` | Membangun binary produksi mandiri ke folder `bin/` |
| `go run cmd/make/main.go make:module <name>` | Membuat boilerplate modul domain baru secara instan |
