# 🚀 Phase 1: Analisis Konsep Eksisting & Rencana Arsitektur Golang Port

Dokumen ini merupakan panduan utama porting **`hono-skeleton-started-kit`** dari ekosistem **TypeScript / Bun / Hono** ke **Golang**. 

Tujuan utama dari porting ini adalah mempertahankan **developer experience (DX) elegan ala Laravel/Hono**, modularitas arsitektur berlapis (*Layered Architecture*), serta kapabilitas production-grade (OpenAPI, background worker, mailer, caching, metrics) dengan memanfaatkan performa tinggi, konkurensi native (*goroutines*), dan stabilitas kompilasi statis dari **Go**.

---

## 📑 Indeks Dokumen

### 1. 🔍 Analisis Referensi Sistem Eksisting (`references/`)
Bedah konsep modular starter kit eksisting agar memiliki landasan konversi yang presisi:
- [01-core-architecture-and-lifecycle.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/references/01-core-architecture-and-lifecycle.md) — Pola Layered Architecture, boundaries, Error handling global, dan lifecycle aplikasi.
- [02-http-routing-openapi-and-middleware.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/references/02-http-routing-openapi-and-middleware.md) — Routing ganda (API & Web), Zod OpenAPI generation, Scalar UI docs, dan middleware pipeline.
- [03-database-layer-and-repository-pattern.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/references/03-database-layer-and-repository-pattern.md) — Drizzle ORM, multi-koneksi DB, schema relational, repository pattern, dan migrasi/seeding.
- [04-queue-system-and-background-workers.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/references/04-queue-system-and-background-workers.md) — BullMQ + Redis, worker process, background job dispatching, dan propagasi context.
- [05-mail-system-and-templating.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/references/05-mail-system-and-templating.md) — Sistem email ala Laravel: Mailable class, Hono JSX email views, driver transports (`log`, `smtp`, `resend`), dan async queue.
- [06-caching-logging-and-observability.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/references/06-caching-logging-and-observability.md) — Redis caching utility, Pino structured logging dengan correlation ID, Prometheus metrics, dan health checks.
- [07-cli-scaffolder-and-dx.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/references/07-cli-scaffolder-and-dx.md) — CLI generator (`bun run make <type> <name>`), stub templating, environment config, testing, dan CI.

### 2. ⚠️ Tantangan & Masalah di Golang (`problems/`)
Pemetaan gap paradigma dan ekosistem antara TypeScript/Node/Bun dan Golang:
- [01-go-paradigm-and-idiomatic-patterns.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/problems/01-go-paradigm-and-idiomatic-patterns.md) — Static compilation vs dynamic runtime, Goroutines vs single-thread event loop, Explicit `error` returns vs exception throwing, `context.Context` vs implicit `AsyncLocalStorage`.
- [02-ecosystem-framework-and-library-gaps.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/problems/02-ecosystem-framework-and-library-gaps.md) — Evaluasi perbandingan library: Web framework, ORM/query builder, validation library, queue system, template engine, dan mailer.
- [03-architectural-translation-challenges.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/problems/03-architectural-translation-challenges.md) — Tantangan implementasi Mailable pattern tanpa OOP inheritance, dependency injection yang bersih tanpa bloated container, dan standard project layout Go.

### 3. 💡 Opsi Solusi Tematik Siap Diskusi (`ideas/`)
Setiap dokumen mengupas opsi teknis, trade-off, dan rekomendasi konkrit untuk kita putuskan bersama satu per satu:
- [01-http-framework-and-openapi-solution.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/ideas/01-http-framework-and-openapi-solution.md) — Pilihan Web Framework (Fiber v3 vs Echo v4 vs Chi) & strategi OpenAPI docs (Swag vs Huma.rocks).
- [02-database-orm-and-migrations-solution.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/ideas/02-database-orm-and-migrations-solution.md) — Pilihan Data Layer (sqlc vs GORM vs uptrace/bun) & Migrasi (Goose / Golang-migrate).
- [03-queue-and-worker-architecture-solution.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/ideas/03-queue-and-worker-architecture-solution.md) — Pilihan Background Worker: hibiken/asynq (Redis 1:1 BullMQ) vs River (Postgres) vs In-memory Goroutine worker.
- [04-mail-system-and-html-templates-solution.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/ideas/04-mail-system-and-html-templates-solution.md) — Desain Mailable Go, multi-driver (`log`, `smtp`, `resend`), dan template engine (`a-h/templ` vs `html/template`).
- [05-cache-logging-and-observability-solution.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/ideas/05-cache-logging-and-observability-solution.md) — Redis wrapper dengan fallback memory, `slog` / `zerolog` ber-trace context, dan Prometheus metrics.
- [06-project-layout-and-cli-scaffolding-solution.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/ideas/06-project-layout-and-cli-scaffolding-solution.md) — Standar struktur direktori Go (`cmd/`, `internal/`, `pkg/`), Cobra CLI generator (`go run cmd/make`), dan konfigurasi Cleanenv/Viper.

### 4. 🛠️ Implementasi (`implementations/`)
- [README.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/implementations/README.md) — Dikosongkan sementara menunggu keputusan ide.
- [STATUS.md](file:///home/melodelavic/Documents/bun/hono-skeleton-started-kit/plans/golang-starter-kit/phase-1/implementations/STATUS.md) — Resume point ringkas status perencanaan.
