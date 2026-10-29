# 01. Multi-Database Setup (`uptrace/bun`)

Starter kit ini menggunakan [`uptrace/bun`](https://bun.uptrace.dev/) sebagai lapisan ORM / query builder SQL, dengan pengelola koneksi terpusat [`DBManager`](file:///home/melodelavic/Documents/golang/golang-fiber-starter-kit/pkg/database/database.go).

---

## 🔌 Driver yang Didukung

Aplikasi dapat terhubung secara simultan ke beberapa jenis database sekaligus:

1. **SQLite (Default):**
   Menggunakan driver pure-Go `modernc.org/sqlite` via `sqliteshim.ShimName`. Tidak memerlukan dependensi CGO sama sekali.
2. **PostgreSQL:**
   Didukung via driver `pgdriver`.
3. **MySQL:**
   Didukung via driver resmi `go-sql-driver/mysql`.
4. **MongoDB (v2):**
   Didukung via `go.mongodb.org/mongo-driver/v2/mongo`.

---

## ⚙️ Konfigurasi di `.env`

Anda dapat mengaktifkan database yang diinginkan melalui variabel lingkungan:

```env
# Koneksi aktif utama
DB_DEFAULT_CONNECTION=sqlite1

# SQLite 1
DB_SQLITE_ENABLED=true
DB_SQLITE1_PATH=storages/database/sqlite1.db

# PostgreSQL (Opsional)
DB_POSTGRES_ENABLED=false
DB_POSTGRES1_DSN="postgres://user:password@localhost:5432/my_db?sslmode=disable"

# MySQL (Opsional)
DB_MYSQL_ENABLED=false
DB_MYSQL1_DSN="user:password@tcp(localhost:3306)/my_db?parseTime=true"
```

---

## 📖 Mengakses Koneksi Database di Kode

Dalam dependensi injeksi service atau repository, Anda dapat mengambil instance koneksi `*bun.DB`:

```go
package main

import (
    "golang-starter-kit/pkg/database"
)

func RegisterServices(dbManager *database.DBManager) {
    // 1. Mengambil database default (sesuai DB_DEFAULT_CONNECTION di .env)
    defaultDB := dbManager.DefaultSQL()

    // 2. Mengambil koneksi database spesifik berdasarkan nama registrasi
    pgDB, err := dbManager.GetSQL("postgres1")
    if err == nil {
        // Query ke PostgreSQL
    }

    // 3. Mengambil koneksi MongoDB
    mongoClient, err := dbManager.GetMongo("mongo1")
    if err == nil {
        // Operasi dokumen MongoDB
    }
}
```

---

## 🧹 Graceful Shutdown

Saat server menerima sinyal terminasi (`SIGINT` atau `SIGTERM`), method `dbManager.CloseAll(ctx)` otomatis dipanggil untuk menutup seluruh connection pool SQL dan mendiskoneksikan client MongoDB secara bersih tanpa kebocoran memori atau koneksi menggantung.
