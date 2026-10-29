# 02. Instalasi & Setup Lingkungan

Panduan langkah demi langkah untuk mempersiapkan lingkungan pengembangan lokal Anda.

---

## 📋 Prasyarat Sistem

Sebelum memulai, pastikan perangkat Anda memiliki perangkat lunak berikut:

1. **Go (Golang):** Versi `1.22` atau yang lebih baru (direkomendasikan versi terbaru).
   Cek versi Anda dengan:
   ```bash
   go version
   ```
2. **Make:** Utilitas otomatisasi tugas bawaan Linux / macOS.
   ```bash
   make --version
   ```
3. **a-h/templ CLI:** Generator template HTML Go.
   Pasang dengan perintah:
   ```bash
   go install github.com/a-h/templ/cmd/templ@latest
   ```
   Pastikan direktori `$GOPATH/bin` (biasanya `~/go/bin`) sudah masuk ke dalam `$PATH` sistem Anda.

4. **Redis (Opsional):** Hanya diperlukan jika Anda ingin menggunakan caching terdistribusi atau mengaktifkan antrean background worker Asynq. Untuk pengujian lokal awal, sistem sudah menyediakan fallback *in-memory cache*.

---

## ⚙️ Langkah Instalasi Proyek

### 1. Salin Berkas Environment
Salin file template `.env.example` menjadi `.env`:
```bash
cp .env.example .env
```

Buka file [`.env`](file:///home/melodelavic/Documents/golang/golang-fiber-starter-kit/.env) dan sesuaikan pengaturannya jika diperlukan:
```env
# Server
APP_NAME="Go Starter Kit"
APP_ENV=local
PORT=3000
LOG_LEVEL=debug

# Database (Default menggunakan SQLite tanpa setup rumit)
DB_DEFAULT_CONNECTION=sqlite1
DB_SQLITE_ENABLED=true
DB_SQLITE1_PATH=storages/database/sqlite1.db

# Mail Driver (Gunakan 'log' untuk mencetak email ke konsol terminal)
MAIL_DRIVER=log
MAIL_FROM_ADDRESS=no-reply@starter.kit
MAIL_FROM_NAME="Go Starter Kit"
```

### 2. Download Dependensi Go
Unduh seluruh dependensi modul Go:
```bash
go mod download
```

### 3. Jalankan Migrasi Database
Secara default starter kit menggunakan SQLite murni yang disimpan di `storages/database/sqlite1.db`. Jalankan migrasi skema tabel awal:
```bash
make migrate
```

Output yang diharapkan:
```text
goose: successfully migrated database to version: 20261029000001
```

### 4. Kompilasi Template HTML (Jika Menggunakan Web)
Jalankan kompilasi awal template `templ`:
```bash
make templ
```

### 5. Inisialisasi Tailwind CSS (Opsional)
Jika Anda berencana menggunakan antarmuka web ber-Tailwind CSS, unduh standalone binary CLI:
```bash
make tailwind-install
```
Perintah ini akan mendownload binary `bin/tailwindcss` yang cocok dengan sistem operasi Anda secara otomatis.
