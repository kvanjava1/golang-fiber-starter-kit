# 🛠️ Analisis Referensi: CLI Scaffolder & Developer Experience (DX)

Dokumen ini membedah generator kode berbasis CLI (`make`), manajemen stub template, pengujian otomatis, dan pipeline CI/CD pada starter kit eksisting.

---

## 1. CLI Code Scaffolder (`scripts/make.ts`)

Salah satu fitur unggulan yang terinspirasi dari `php artisan make:*` di Laravel adalah generator kode otomatis:

### Perintah Generator yang Didukung:
```bash
bun run make controller <domain> <name>   # Membuat Controller (API & Web)
bun run make service <domain> <name>      # Membuat Service layer
bun run make repository <domain> <name>   # Membuat Repository layer
bun run make validation <domain> <name>   # Membuat Zod validation schema
bun run make mailable <name>              # Membuat Mailable class + JSX template
bun run make job <domain> <name>          # Membuat BullMQ Job handler
bun run make view <name>                  # Membuat Hono JSX view component
```

### Mekanisme Kerja:
1. Membaca template file mentah di folder `scripts/stubs/*.stub`.
2. Melakukan substitusi placeholder token: `{{Name}}`, `{{name}}`, `{{Domain}}`, `{{domain}}`, `{{PascalCase}}`.
3. Memastikan direktori target ada (`mkdir -p`), lalu menulis file baru tanpa menimpa file yang sudah ada (*safety check*).
4. Menambahkan import otomatis pada barrel file (`index.ts`) jika diperlukan.

---

## 2. Configuration & Environment Management (`src/configs/env.ts`)

- Seluruh environment variables didefinisikan secara deklaratif menggunakan skema Zod.
- Jika ada variabel `.env` yang wajib (*required*) namun tidak terisi atau salah format (misal: format port bukan angka), aplikasi akan langsung menolak start (*fail fast*) dengan pesan error yang jelas dan detail di terminal.

---

## 3. Testing Strategy & Continuous Integration

- **Testing Framework:** Menggunakan `bun:test` bawaan untuk kecepatan eksekusi tinggi (ratusan test selesai dalam hitungan ratusan milidetik).
- **Struktur Test:**
  - `tests/unit/`: Menguji logika murni Service, Util, Mail facade, dan formatting secara terisolasi.
  - `tests/integration/`: Menguji endpoint HTTP lengkap dari Router, Middleware, Service, hingga interaksi Database SQLite.
- **CI/CD Pipeline (`.github/workflows/ci.yml`):**
  - Menjalankan linter & typecheck (`tsc --noEmit`).
  - Menjalankan seluruh test suite dengan environment mock Redis & SQLite.
  - Memastikan 0 error sebelum kode dapat digabungkan (*merge*) ke branch utama.
