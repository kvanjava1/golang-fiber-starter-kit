# 01. Overview & Fitur Unggulan

**Go Fiber Starter Kit** adalah blueprint aplikasi backend yang dirancang untuk memberikan fleksibilitas maksimal bagi pengembang:
1. **API-First:** Sangat ringan dan cepat ketika digunakan untuk membangun REST API murni berkinerja tinggi.
2. **Full-Stack SSR Ready:** Siap menyajikan antarmuka pengguna web interaktif berbasis Server-Side Rendering (SSR) menggunakan `a-h/templ` dan Tailwind CSS tanpa perlu repot mengatur bundler JavaScript yang rumit.

---

## 🚀 Fitur Utama & Matriks Teknologi

| Kategori | Teknologi | Deskripsi / Alasan Pemilihan |
| :--- | :--- | :--- |
| **HTTP Engine** | [GoFiber v3](https://docs.gofiber.io/) | Framework web berbasis Fasthttp dengan performa ekstrem, zero memory allocation pada hot path, dan ergonomi ala Express. |
| **Database ORM** | [uptrace/bun](https://bun.uptrace.dev/) | SQL-first ORM yang cepat, fleksibel, mendukung pagination, subquery, dan multi-database tanpa bloat. |
| **Drivers Database** | SQLite, PostgreSQL, MySQL, MongoDB | Menggunakan driver pure-Go (`modernc.org/sqlite` / `sqliteshim`) sehingga CGO dinonaktifkan (`CGO_ENABLED=0`). |
| **Migrasi Database**| [pressly/goose v3](https://github.com/pressly/goose) | Standard industri migrasi database berbasis file SQL murni yang mendukung up/down rollback. |
| **Templating Web**  | [a-h/templ](https://templ.guide/) | Template HTML modern yang dikompilasi langsung menjadi kode Go murni dengan *type-safety* saat compile-time. |
| **Styling UI**     | [Tailwind CSS v3 Standalone](https://tailwindcss.com/) | Utilitas CSS modern yang diunduh langsung sebagai binary mandiri tanpa memerlukan Node.js atau npm. |
| **Task Queue**     | [hibiken/asynq](https://github.com/hibiken/asynq) | Antrean tugas latar belakang (*background jobs*) berbasis Redis dengan fitur retry, delay, dan dead-letter queue. |
| **Mail System**    | [go-mail/mail](https://github.com/go-mail/mail) | Sistem pengiriman email berbasis interface fluent yang mendukung driver `log` (pengujian lokal), `smtp`, dan `resend`. |
| **Caching**        | `redis` & In-Memory | Lapisan cache dengan fallback otomatis ke memory lokal jika Redis tidak diaktifkan, dilengkapi tag-based invalidation. |
| **Observability**  | `log/slog` & Prometheus | Structured logging standar Go 1.21+ dengan Trace Correlation ID terintegrasi, serta endpoint metrik `/metrics`. |
| **Config Loader**  | [cleanenv](https://github.com/ilyakaznacheev/cleanenv) | Membaca variabel environment dan file `.env` dengan validasi tag struct yang aman. |
| **Validasi Input** | [go-playground/validator v10](https://github.com/go-playground/validator) | Validasi struct request dengan pesan error yang rapi dan terstandardisasi. |

---

## 🏛️ Filosofi Desain

1. **Zero-CGO & Single Binary Portability:**
   Aplikasi dapat dikompilasi menjadi satu file binary statis murni yang dapat langsung dijalankan di server Linux mana pun tanpa dependensi C runtime (`libc/musl`).
2. **Domain-Driven Modular Structure:**
   Setiap fitur bisnis dikelompokkan ke dalam satu folder domain (misal: `internal/domain/user/` atau `internal/domain/example/`) yang memiliki kontrak model, DTO, repository, service, controller, dan routing sendiri.
3. **Pemisahan Jalur Web & API:**
   Router aplikasi mendistribusikan lalu lintas secara terpisah:
   - Jalur Web melayani HTML SSR pada root path (misal: `/example/cruds`).
   - Jalur API melayani JSON response berstandar envelope pada `/api/v1/example/cruds`.
4. **On-Demand Tooling:**
   Jika Anda hanya ingin membuat REST API, Anda tidak perlu menjalankan watcher CSS atau menginstall alat bantu frontend.
