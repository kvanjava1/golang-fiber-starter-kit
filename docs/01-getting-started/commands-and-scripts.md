# 03. Panduan Perintah CLI & Makefile

Starter kit ini dilengkapi dengan serangkaian perintah praktis melalui `Makefile` dan CLI Scaffolder mandiri di [`cmd/make`](file:///home/melodelavic/Documents/golang/golang-fiber-starter-kit/cmd/make/main.go) untuk memaksimalkan kenyamanan dan kecepatan pengembangan.

---

## 🛠️ Perintah Utama `Makefile`

Seluruh perintah umum dapat dieksekusi menggunakan utilitas `make`:

### 1. Menjalankan Server & Aplikasi

| Perintah | Deskripsi | Kapan Digunakan |
| :--- | :--- | :--- |
| `make dev` | Menjalankan server HTTP (API & Web) secara langsung. Mengompilasi `templ` sebelum start. | Digunakan untuk pengembangan harian REST API maupun Web SSR standar. |
| `make dev-web` | Mengompilasi `templ`, membangun minified CSS Tailwind, lalu menyalakan server. | Digunakan jika Anda aktif mengembangkan tampilan antarmuka web ber-Tailwind. |
| `make worker` | Menjalankan background worker antrean Asynq (`cmd/worker/main.go`). | Digunakan saat Anda perlu menguji pemrosesan background job / antrean email. |

### 2. Database & Migrasi

| Perintah | Deskripsi |
| :--- | :--- |
| `make migrate` | Menjalankan seluruh file migrasi pending pada database aktif (`sqlite1`). |
| `make migrate-status` | Menampilkan tabel status migrasi (versi mana saja yang sudah/belum dieksekusi). |
| `make migrate-rollback`| Melakukan rollback 1 batch migrasi terakhir. |

### 3. Frontend & Styling (On-Demand)

| Perintah | Deskripsi |
| :--- | :--- |
| `make templ` | Menjalankan `templ generate` untuk mengompilasi file `.templ` menjadi file `.go`. |
| `make tailwind-install` | Mendownload binary standalone Tailwind CLI ke `bin/tailwindcss`. |
| `make css` | Melakukan build & minifikasi CSS sekali jalan ke `public/css/output.css`. |
| `make css-watch` | Menjalankan live CSS watcher interaktif di terminal terpisah. |

### 4. Pengujian & Kompilasi Binary

| Perintah | Deskripsi |
| :--- | :--- |
| `make test` | Menjalankan seluruh test suite dengan Race Detector aktif (`go test -v -race ./...`). |
| `make build` | Mengompilasi binary mandiri untuk produksi ke dalam folder `bin/` (`server`, `worker`, `make`). |
| `make clean` | Menghapus binary build aplikasi yang ada di folder `bin/`. |

---

## 🪄 Perintah CLI Scaffolding (`cmd/make`)

Perintah generator kode dapat dijalankan langsung melalui `go run cmd/make/main.go <perintah>`:

### 1. Membuat Modul Domain Baru (`make:module`)
Secara instan menghasilkan boilerplate modul domain lengkap dengan clean architecture:
```bash
go run cmd/make/main.go make:module product
```
File yang otomatis dibuat:
- `internal/domain/product/model.go` (Definisi entitas database)
- `internal/domain/product/dto.go` (Struct request validation & response)
- `internal/domain/product/repository.go` (Interface & implementasi query Bun)
- `internal/domain/product/service.go` (Logika bisnis & orkestrasi cache)
- `internal/domain/product/controller_api.go` (Handler REST API JSON)
- `internal/domain/product/controller_web.go` (Handler Web SSR HTML)
- `internal/domain/product/route.go` (Pendaftaran endpoint Web & API)

### 2. Membuat File Migrasi Database (`make:migration`)
Membuat template migrasi SQL berformat Goose:
```bash
go run cmd/make/main.go make:migration create_products_table
```
File dihasilkan di: `resources/migrations/sqlite1/<timestamp>_create_products_table.sql`.

### 3. Membuat Seeder Database (`make:seeder`)
Membuat file seeder transaksional:
```bash
go run cmd/make/main.go make:seeder product
```
File dihasilkan di: `database/seeders/product_seeder.go`.

### 4. Membuat Class Email Notifikasi (`make:mail`)
Membuat class mailable baru untuk pengiriman email:
```bash
go run cmd/make/main.go make:mail invoice_paid
```
File dihasilkan di: `pkg/mail/mailables/invoice_paid_mail.go`.

### 5. Mengatur Target Database Migrasi
Jika menggunakan arsitektur multi-database, Anda dapat menunjuk database target secara eksplisit menggunakan flag `--db`:
```bash
go run cmd/make/main.go migrate:up --db=postgres1
go run cmd/make/main.go migrate:status --db=mysql1
go run cmd/make/main.go migrate:down --db=postgres1
```
