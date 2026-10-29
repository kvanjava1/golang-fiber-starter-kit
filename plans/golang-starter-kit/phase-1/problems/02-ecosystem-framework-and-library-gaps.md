# ⚠️ Analisis Masalah: Pemetaan Gap Ekosistem & Library

Dokumen ini memetakan ketiadaan kesetaraan 1:1 (*parity gaps*) antara pustaka yang digunakan di TypeScript/Hono dengan ekosistem Golang, serta tantangan dalam memilih alternatif terbaik.

---

## 1. Tabel Perbandingan Pustaka: TypeScript vs Golang

| Kebutuhan Fitur | Solusi di Hono Kit (TS) | Tantangan / Gap di Golang | Alternatif Kandidat di Go |
|---|---|---|---|
| **Web Routing Framework** | Hono | Go tidak memiliki framework dominan tunggal seperti Express/Hono | **Fiber v3**, **Echo v4**, **Chi**, **Gin** |
| **OpenAPI / Swagger Docs** | `@hono/zod-openapi` + Scalar | Di Go, OpenAPI sering dihasilkan via komentar swagger (`swag`) atau code-generator terpisah | **swag**, **huma.rocks**, **oapi-codegen** |
| **Validasi Schema & DTO** | Zod (Runtime + TS inference) | Go mengandalkan struct tags & reflection, atau code-generation | **go-playground/validator**, **ozzo-validation** |
| **ORM & Query Builder** | Drizzle ORM (Type-safe SQL) | ORM tradisional di Go (GORM) lambat & banyak reflection; query builder murni membutuhkan banyak boilerplate | **sqlc** (Type-safe SQL), **GORM**, **uptrace/bun** |
| **Database Migrations** | Drizzle Migrator | Perlu tool migrasi database berbasis CLI/Go terpisah | **pressly/goose**, **golang-migrate/migrate** |
| **Async Queue / Jobs** | BullMQ (Redis-based) | BullMQ eksklusif untuk Node/TS; Go butuh engine Redis task queue yang setara | **hibiken/asynq**, **River** (Postgres), **Machinery** |
| **Email Transports** | Nodemailer + Resend SDK | Go butuh library SMTP tangguh dan wrapper REST API | **wneessen/go-mail**, `net/smtp`, **resend-go** |
| **Email HTML Templates** | Hono JSX Components | Go tidak memiliki JSX bawaan; string concatenation rentan XSS | **a-h/templ** (JSX-like), `html/template` standar |
| **Structured Logging** | Pino | Perlu logger berkinerja tinggi yang kompatibel dengan context | `log/slog` (Stdlib 1.21+), **rs/zerolog** |
| **Distributed Caching** | ioredis | Perlu driver Redis yang efisien dan type-safe | **redis/go-redis/v9** |

---

## 2. Masalah Utama: Skema Terpadu vs Pemisahan DTO & OpenAPI

### Kondisi di Hono:
Satu definisi Zod schema langsung menghasilkan:
1. Validasi runtime input HTTP.
2. Type inference TypeScript untuk request body.
3. Dokumen OpenAPI JSON 3.1 otomatis tanpa menulis anotasi tambahan.

### Gap di Golang:
- Pendekatan konvensional di Go (seperti Gin/Fiber + `swag`) mewajibkan developer menulis komentar anotasi yang panjang di atas setiap controller:
  ```go
  // @Summary Create CRUD
  // @Param request body CreateCrudRequest true "Payload"
  // @Success 201 {object} CrudResponse
  // @Router /api/example/cruds [post]
  ```
- Pendekatan ini rentan *out-of-sync* (dokumentasi tidak sesuai dengan implementasi nyata jika developer lupa memperbarui anotasi).
- **Pertanyaan Desain:** Apakah kita ingin menggunakan framework modern berbasis schema seperti **Huma** (`huma.rocks`) yang memberikan schema-first OpenAPI otomatis ala Hono, atau menggunakan pendekatan standar industri Go via **Echo/Fiber + swag**?

---

## 3. Masalah Utama: Drizzle ORM vs Ekosistem Database Go

- Di TypeScript, Drizzle ORM sangat digemari karena memberikan pengalaman "menulis SQL tetapi ber-tipe aman (*type-safe*)".
- Di Go:
  - **GORM:** Mirip Laravel Eloquent, sangat mudah digunakan, tetapi menggunakan banyak `interface{}` dan reflection yang tidak disukai sebagian gopher purist.
  - **sqlc:** Menulis SQL mentah di file `.sql`, lalu compiler sqlc menghasilkan kode Go struct dan query yang 100% type-safe tanpa reflection. Ini adalah alternatif paling mendekati filosofi performa dan type-safety Drizzle.
  - **uptrace/bun:** Query builder SQL di Go yang berada di tengah-tengah antara GORM dan sqlc.
