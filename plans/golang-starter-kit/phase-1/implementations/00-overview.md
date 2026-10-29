# 🏛️ Implementasi: 00 - Blueprint Teknis & Arsitektur Sistem

Dokumen ini memuat cetak biru teknis komprehensif, diagram relasi paket, dan alur data end-to-end untuk pembangunan **Golang Starter Kit** berbasis **Go-Fiber v3**, **uptrace/bun**, **hibiken/asynq**, **a-h/templ**, dan **log/slog**.

---

## 1. Diagram Arsitektur & Relasi Paket (Layer Dependency)

```mermaid
flowchart TD
    subgraph Binaries ["Entry Points (cmd/)"]
        ServerBin["cmd/server/main.go (HTTP Web & API)"]
        WorkerBin["cmd/worker/main.go (Asynq Background Worker)"]
        MakeBin["cmd/make/main.go (CLI Scaffolder & Migrator)"]
    end

    subgraph CompositionRoot ["Composition Root (internal/app/)"]
        App["internal/app/app.go (BuildContainer & Bootstrap)"]
        Router["internal/app/router.go (Master Route Delegator)"]
    end

    subgraph DomainModules ["Domain Modules (internal/domain/<name>/)"]
        Route["route.go (Self-Registering Routes)"]
        CtrlWeb["controller_web.go (HTML SSR Handler)"]
        CtrlAPI["controller_api.go (REST JSON Handler)"]
        Service["service.go (Pure Domain Logic)"]
        Repo["repository.go (Interface & Bun Queries)"]
        Model["model.go & dto.go"]
    end

    subgraph TechnicalPackages ["Technical Packages (pkg/)"]
        PkgDB["pkg/database/ (Multi-DB Manager: Bun & Mongo)"]
        PkgCache["pkg/cache/ (Redis & Tag Invalidation)"]
        PkgQueue["pkg/queue/ (Asynq Client & Tasks)"]
        PkgMail["pkg/mail/ (Mailable & Drivers)"]
        PkgResp["pkg/response/ (JSON Envelope & Error Mapper)"]
        PkgLog["pkg/logger/ (slog Trace Handler)"]
        PkgVal["pkg/validator/ (go-playground/validator)"]
    end

    subgraph ResourcesAndStorages ["Resources & Storages"]
        Views["resources/views/ (a-h/templ Templates)"]
        Migrations["resources/migrations/ (Goose SQL Migrations)"]
        Seeders["database/seeders/ (Transactional Seeders)"]
        Storages["storages/ (database/, logs/, uploads/)"]
    end

    ServerBin --> App
    WorkerBin --> PkgQueue
    MakeBin --> Migrations
    MakeBin --> Seeders
    
    App --> Router
    Router --> Route
    Route --> CtrlWeb
    Route --> CtrlAPI
    
    CtrlWeb --> Service
    CtrlAPI --> Service
    CtrlWeb --> Views
    
    Service --> Repo
    Service --> PkgCache
    Service --> PkgMail
    Service --> PkgQueue
    
    Repo --> PkgDB
    PkgDB --> Storages
    PkgMail --> PkgQueue
```

---

## 2. Alur Request HTTP (Web SSR & REST API)

```mermaid
sequenceDiagram
    autonumber
    actor Client as Web Browser / Mobile App
    participant Fib as Fiber HTTP Pipeline
    participant MW as Middlewares (Trace, Logger, RateLimiter)
    participant Ctrl as Domain Controller (Web / API)
    participant Svc as Domain Service
    participant Repo as Domain Repository
    participant DB as Database Manager (*bun.DB)

    Client->>Fib: HTTP Request (Headers: X-Correlation-ID / Form / JSON)
    Fib->>MW: Jalankan Middleware Chain
    MW->>MW: Inject correlation_id ke context.Context
    MW->>Ctrl: Panggil Handler
    Ctrl->>Ctrl: Validasi input DTO (pkg/validator)
    alt Input Tidak Valid
        Ctrl-->>Client: 400 Bad Request { success: false, errors: [...] }
    end
    Ctrl->>Svc: Panggil Domain Method (ctx, input)
    Svc->>Repo: Query Database (ctx, params)
    Repo->>DB: Eksekusi SQL via uptrace/bun
    DB-->>Repo: Hasil Query
    Repo-->>Svc: Domain Model
    Svc-->>Ctrl: Result / Domain Error
    alt Tipe Request: Web SSR
        Ctrl->>Ctrl: Render template a-h/templ ke HTML
        Ctrl-->>Client: HTTP 200 text/html
    else Tipe Request: REST API
        Ctrl->>Ctrl: Bungkus ke pkg/response.Envelope
        Ctrl-->>Client: HTTP 200 application/json
    end
```

---

## 3. Alur Background Worker & Mail Dispatch

```mermaid
sequenceDiagram
    autonumber
    actor User as User / Service
    participant Mail as pkg/mail Facade
    participant Templ as a-h/templ Engine
    participant AsynqClient as Asynq Task Queue (Redis)
    participant Worker as Asynq Worker Process
    participant Driver as Mail Driver (Log / SMTP / Resend)

    User->>Mail: mail.To("user@test.com").Queue(ctx, welcomeMail)
    Mail->>Templ: welcomeMail.Render(ctx) -> Pre-render ke String HTML
    Templ-->>Mail: "<!DOCTYPE html><html>...</html>"
    Mail->>AsynqClient: Enqueue Task ("email:send", payload, opts)
    AsynqClient-->>User: Nil (Enqueued dalam < 3ms)
    
    Note over AsynqClient,Worker: Pemrosesan Asinkron di Worker Process
    AsynqClient->>Worker: Dequeue Task
    Worker->>Worker: Restore correlation_id ke context.Context
    Worker->>Driver: Driver.Send(ctx, message)
    alt Pengiriman Berhasil
        Driver-->>Worker: nil
        Worker->>AsynqClient: Mark Completed
    else Gagal (SMTP Down / Timeout)
        Driver-->>Worker: error
        Worker->>AsynqClient: Schedule Retry (Exponential Backoff)
    end
```

---

## 4. Matriks Teknologi & Pustaka Utama

| Komponen Arsitektur | Teknologi / Pustaka | Alasan Pemilihan & Keunggulan |
|---|---|---|
| **Language & Runtime** | Golang 1.22+ | Static compilation (single binary), hemat RAM (~25MB), goroutines native |
| **HTTP Web Framework** | **Go-Fiber v3** | Engine `fasthttp` super cepat, sintaks 95% mirip Hono/Express, ekosistem middleware lengkap |
| **Relational Database** | **`uptrace/bun`** | SQL query builder fluent, mendukung `sqlite1..3`, `mysql1..3`, `pg1..3` dalam satu interface |
| **NoSQL Database** | **`mongo-driver/v2`** | Official MongoDB Go driver untuk `mongo1..3` |
| **Database Migrations** | **`pressly/goose`** | Standar industri migrasi SQL, mendukung migrasi terarah per database target |
| **Database Seeders** | **Custom Transactional Runner** | Seeder modular, dibungkus `RunInTx` dan *idempotent* (anti-duplikasi) |
| **Task Queue & Worker** | **`hibiken/asynq` + Redis** | Klon 1:1 BullMQ di Go, retries, exponential backoff, Web Dashboard `asynqmon` |
| **Email HTML Templates** | **`a-h/templ`** | Template berbasis komponen mirip JSX, 100% type-safe saat kompilasi |
| **Email Transports** | **`wneessen/go-mail` + Resend** | Multi-driver transport (`log`, `smtp`, `resend`), aman untuk dev/test |
| **Structured Logging** | **`log/slog` (Stdlib)** | Logging berkecepatan tinggi tanpa dependensi luar, otomatis membawa `correlation_id` |
| **Caching Layer** | **`redis/go-redis/v9`** | Pola Cache-Aside (`Remember`) dengan *Tag-based Invalidation* via Redis Sets |
| **Observabilitas** | **`prometheus/client_golang`** | Metrik standar Prometheus (`/metrics`) dan deep health check (`/api/v1/health`) |
| **Config & Env** | **`ilyakaznacheev/cleanenv`** | Fail-fast startup: otomatis validasi file `.env` dan tipe data variabel |
| **CLI Scaffolder** | **`cmd/make`** | Tool pembuat kode instan (module, controller, service, repo, mail, job, seeder) |
