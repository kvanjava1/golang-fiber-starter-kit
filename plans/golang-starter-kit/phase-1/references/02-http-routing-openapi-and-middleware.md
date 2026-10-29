# 🌐 Analisis Referensi: HTTP Routing, OpenAPI & Middleware Pipeline

Dokumen ini membedah arsitektur routing ganda (API vs Web SSR), integrasi otomatis spesifikasi OpenAPI / Scalar Docs, serta rantai (*pipeline*) middleware.

---

## 1. Dual Routing Structure (API vs Web SSR)

Starter kit memisahkan secara fisik antara endpoint RESTful API dan antarmuka Web SSR:

```text
src/routes/
├── index.ts          # Mount agregat: menggabungkan /api dan /web
├── api/
│   ├── index.ts      # Mount router API (/api/*)
│   ├── docs/         # OpenAPI JSON spec (/api/docs/openapi.json) & UI (/api/docs)
│   ├── health/       # Health check probes (/api/health)
│   ├── metrics/      # Prometheus telemetry (/metrics)
│   └── example/      # Domain endpoint (misal: /api/example/cruds)
└── web/
    ├── index.ts      # Mount router Web (SSR)
    └── example/      # Tampilan Web SSR Hono JSX
```

---

## 2. OpenAPI Specification & Schema-First Validation

Starter kit menggunakan `@hono/zod-openapi` untuk menghubungkan validasi input, dokumentasi OpenAPI, dan definisi route dalam satu sumber kebenaran (*single source of truth*).

### Alur Definisi Route OpenAPI:
1. **Schema DTO (`src/validations/`):** Didefinisikan menggunakan Zod schemas dengan metadata OpenAPI (`openapi: { example: ..., description: ... }`).
2. **OpenAPI Route Spec (`createRoute`):**
   ```typescript
   export const createCrudRoute = createRoute({
     method: "post",
     path: "/api/example/cruds",
     request: {
       body: {
         content: {
           "application/json": { schema: CreateCrudBody }
         }
       }
     },
     responses: {
       201: {
         description: "Resource created",
         content: { "application/json": { schema: CrudResponseSchema } }
       },
       400: {
         description: "Validation error",
         content: { "application/json": { schema: ErrorResponseSchema } }
       }
     }
   });
   ```
3. **Dokumentasi UI:** Disajikan melalui **Scalar API Reference** (`@scalar/hono-api-reference`) pada route `/api/docs` yang membaca langsung OpenAPI 3.1 JSON dari `/api/docs/openapi.json`.

---

## 3. Middleware Pipeline

Pipeline middleware disusun dengan urutan eksekusi (*order of execution*) yang ketat:

```mermaid
flowchart TD
    Req(["Incoming HTTP Request"]) --> M1["1. Request ID & Correlation ID"]
    M1 --> M2["2. Context Initializer (AsyncLocalStorage)"]
    M2 --> M3["3. Metrics Middleware (Prometheus Timer)"]
    M3 --> M4["4. Logger Middleware (Pino log request/response)"]
    M4 --> M5["5. Security Headers (CSP, HSTS, X-Frame-Options)"]
    M5 --> M6["6. CORS Middleware"]
    M6 --> M7["7. Rate Limiter (Memory Token Bucket)"]
    M7 --> Router{"Router Matching"}
    
    Router -->|Path /api/*| M8["API Specific: Content-Type & Error Envelope"]
    Router -->|Path /*| M9["Web Specific: Static Assets & Vite Manifest"]
    
    M8 --> Handler["Controller Handler"]
    M9 --> Handler
```

### Rincian Middleware:
- **`context.middleware.ts`**: Menangkap header `X-Request-Id` (atau menghasilkan UUID baru) dan menyimpannya ke request context.
- **`logger.middleware.ts`**: Mencatat waktu mulai request, status HTTP, dan latency (ms) saat request selesai.
- **`rateLimiter.middleware.ts`**: Membatasi laju request per IP menggunakan in-memory store dengan mekanisme sliding-window / token-bucket.
- **`metrics.middleware.ts`**: Mengukur durasi HTTP request (`http_request_duration_seconds`) dan menghitung jumlah request berdasarkan path dan status code (`http_requests_total`).
- **`errorHandler.ts`**: Menangkap kegagalan di seluruh titik pipeline dan memformat output menjadi JSON terstandarisasi.
