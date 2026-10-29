# 🗺️ Master Plan: Golang Starter Kit Port

Dokumen ini memuat analisis mendalam terhadap **`hono-skeleton-started-kit`** (Bun + Hono + TypeScript) serta spesifikasi teknis lengkap untuk membangun kembali starter kit ini dalam bahasa pemrograman **Golang**.

Perencanaan disusun dengan standar skill **`software-engineering-and-architecture`** ([E10] Technical Specification & Feature Planning Pipeline).

---

## 📂 Struktur Rencana Kerja (`plans/golang-starter-kit/phase-1/`)

```text
plans/golang-starter-kit/phase-1/
├── README.md                           ← Entry point fase 1 & indeks dokumen
├── references/                         ← Analisa tematik sistem eksisting (Bun + Hono) [7 File]
│   ├── 01-core-architecture-and-lifecycle.md
│   ├── 02-http-routing-openapi-and-middleware.md
│   ├── 03-database-layer-and-repository-pattern.md
│   ├── 04-queue-system-and-background-workers.md
│   ├── 05-mail-system-and-templating.md
│   ├── 06-caching-logging-and-observability.md
│   └── 07-cli-scaffolder-and-dx.md
├── problems/                           ← Analisa tantangan & gap konversi ke Golang [3 File]
│   ├── 01-go-paradigm-and-idiomatic-patterns.md
│   ├── 02-ecosystem-framework-and-library-gaps.md
│   └── 03-architectural-translation-challenges.md
├── ideas/                              ← Opsi solusi tematik & keputusan final [6 File - ALL DECIDED]
│   ├── 01-http-framework-and-openapi-solution.md    [✅ Go-Fiber v3]
│   ├── 02-database-orm-and-migrations-solution.md   [✅ uptrace/bun + Multi-DB + Goose + Seeders]
│   ├── 03-queue-and-worker-architecture-solution.md [✅ hibiken/asynq + Redis]
│   ├── 04-mail-system-and-html-templates-solution.md [✅ a-h/templ + Multi-Driver]
│   ├── 05-cache-logging-and-observability-solution.md [✅ log/slog + go-redis + prometheus]
│   └── 06-project-layout-and-cli-scaffolding-solution.md [✅ Standard Go + storages/ + cleanenv]
└── implementations/                    ← Rencana implementasi teknis lengkap [10 File - READY]
    ├── README.md                       ← Master indeks langkah implementasi
    ├── STATUS.md                       ← Resume point sesi baru (anti-gemuk)
    ├── 00-overview.md                  ← Blueprint teknis, alur data Mermaid & matriks teknologi
    ├── 01-step-module-init-and-config.md
    ├── 02-step-pkg-core-utilities.md
    ├── 03-step-multi-database-and-migrations.md
    ├── 04-step-fiber-server-and-middlewares.md
    ├── 05-step-caching-and-observability.md
    ├── 06-step-queue-worker-and-mail-system.md
    ├── 07-step-domain-crud-reference.md
    ├── 08-step-cli-scaffolder-and-makefile.md
    └── 09-step-testing-and-verification.md
```

---

## 🚦 Status Fase 1: `READY_FOR_EXECUTION`
- [x] **Analisis Konsep Eksisting (`references/`):** 7 file tematik selesai.
- [x] **Pemetaan Masalah & Gap Golang (`problems/`):** 3 file tantangan arsitektur selesai.
- [x] **Diskusi Opsi Solusi (`ideas/`):** 6 topik arsitektur telah selesai disepakati secara utuh.
- [x] **Rencana Implementasi Teknis (`implementations/`):** Blueprint 00-overview dan 9 langkah teknis siap dieksekusi kapan saja.
