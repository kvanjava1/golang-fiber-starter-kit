# 🛠️ Rencana Implementasi: Golang Starter Kit Port

Dokumen ini memuat peta navigasi langkah-langkah implementasi teknis untuk membangun **Golang Starter Kit** dari nol hingga siap produksi.

---

## 📑 Daftar Dokumen Blueprint & Langkah Eksekusi

| No | Dokumen | Ruang Lingkup & Fokus Utama |
|---|---|---|
| **00** | [00-overview.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/implementations/00-overview.md) | Blueprint arsitektur teknis, diagram Mermaid alur data, dan matriks teknologi |
| **01** | [01-step-module-init-and-config.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/implementations/01-step-module-init-and-config.md) | Inisialisasi `go.mod`, layout folder, direktori `storages/`, dan konfigurasi `cleanenv` |
| **02** | [02-step-pkg-core-utilities.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/implementations/02-step-pkg-core-utilities.md) | Paket `pkg/`: Response JSON envelope, error mapper, slog trace handler, dan validator |
| **03** | [03-step-multi-database-and-migrations.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/implementations/03-step-multi-database-and-migrations.md) | Multi-DB Manager (`sqlite1..3`, `mysql1..3`, `pg1..3`, `mongo1..3`), migrasi Goose, dan Seeder |
| **04** | [04-step-fiber-server-and-middlewares.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/implementations/04-step-fiber-server-and-middlewares.md) | HTTP server Go-Fiber v3, pipeline middleware, static uploads, dan master router delegator |
| **05** | [05-step-caching-and-observability.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/implementations/05-step-caching-and-observability.md) | Redis cache (`Remember` & tag invalidation), Prometheus `/metrics`, dan `/health` probes |
| **06** | [06-step-queue-worker-and-mail-system.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/implementations/06-step-queue-worker-and-mail-system.md) | Antrean Asynq Redis, worker runner `cmd/worker`, Mailable interface, dan multi-driver mailer |
| **07** | [07-step-domain-crud-reference.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/implementations/07-step-domain-crud-reference.md) | Living reference domain `example` (Web SSR HTML, REST API JSON, Bun repo, send-email) |
| **08** | [08-step-cli-scaffolder-and-makefile.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/implementations/08-step-cli-scaffolder-and-makefile.md) | CLI Scaffolder `cmd/make` (make:module, make:mail, dll) dan shortcut Makefile |
| **09** | [09-step-testing-and-verification.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/implementations/09-step-testing-and-verification.md) | Test suite (unit & integration via Fiber test), test helpers, dan GitHub Actions CI |

---

## 🚦 Status Penunjuk Eksekusi
Status terkini dari progress implementasi selalu diperbarui di file ringkas [STATUS.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/implementations/STATUS.md).
