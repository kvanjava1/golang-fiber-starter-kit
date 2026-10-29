# 💡 Opsi Solusi 02: Database Layer, ORM, Transaksi & Migrasi

> **Status Keputusan:** ✅ **DECIDED: Opsi C (`uptrace/bun`) + Multi-Connection Manager + `pressly/goose`**  
> **Keputusan:** Dipilih **`uptrace/bun`** sebagai SQL-first query builder terpadu untuk relational database, dengan dukungan **Multiple Named Connections** penuh untuk setiap jenis database: **SQLite (`sqlite1..3`)**, **MySQL (`mysql1..3`)**, **PostgreSQL (`pg1..3`)**, serta official driver MongoDB untuk **MongoDB (`mongo1..3`)**. Untuk migrasi skema SQL digunakan **`pressly/goose`**.

---

## 🎯 Standar Production-Grade yang Harus Dipenuhi
1. **Full Multi-Connection Support per Engine:**
   - **SQLite:** `sqlite1`, `sqlite2`, `sqlite3`
   - **PostgreSQL:** `pg1`, `pg2`, `pg3`
   - **MySQL:** `mysql1`, `mysql2`, `mysql3`
   - **MongoDB:** `mongo1`, `mongo2`, `mongo3`
2. **Repository Interface at Consumer Boundary:** Interface repository didefinisikan di dalam paket domain pengguna, bukan di paket implementasi database (sesuai idiom Go: *"Define interfaces where you use them, not where you implement them"*).
3. **Context-Aware Database Queries:** Seluruh kueri wajib menerima `ctx context.Context` (`QueryContext`, `ExecContext`, atau Drizzle/GORM with context). Jika client membatalkan request HTTP, kueri database yang sedang berjalan otomatis dibatalkan di server database.
4. **Pola Transaksi Database yang Bersih:** Memungkinkan Service mengeksekusi operasi multi-tabel dalam satu transaksi ACID tanpa mengotori domain logic dengan objek driver database (misal `*sql.Tx`).
5. **Migrasi Terkelola per Database:** Migrasi skema database berbasis file versioning terurut (`YYYYMMDDHHMMSS_create_table.sql`) yang dapat ditargetkan ke database spesifik (`--db=sqlite1`, `--db=pg1`, `--db=mysql1`).

---

## 🏛️ Arsitektur Multi-Connection Database Manager

Sama persis seperti kemampuan di `hono-skeleton-started-kit`, starter kit Go akan memiliki **Registry Connection Pool** terpusat di `pkg/database/`:

```mermaid
flowchart TD
    Config[".env Configuration"] --> Manager["DBManager (pkg/database/manager.go)"]
    
    subgraph RelationalPools ["Relational Pools (uptrace/bun)"]
        Manager --> SQLitePool["SQLite Connections: sqlite1, sqlite2, sqlite3"]
        Manager --> PGPool["PostgreSQL Connections: pg1, pg2, pg3"]
        Manager --> MySQLPool["MySQL Connections: mysql1, mysql2, mysql3"]
    end
    
    subgraph MongoPools ["NoSQL Pools (official mongo-driver)"]
        Manager --> MongoMgr["MongoManager: mongo1, mongo2, mongo3"]
    end
    
    SQLitePool --> Repo1["User Repository (menggunakan pg1)"]
    PGPool --> Repo2["Order Repository (menggunakan pg1)"]
    SQLitePool --> Repo3["Audit Repository (menggunakan sqlite2)"]
    MongoPools --> Repo4["Log Document Repository (menggunakan mongo1)"]
```

---

## ⚙️ 1. Format Konfigurasi Multi-Connection (`.env`)

```env
# Relational Default Connection
DB_DEFAULT_CONNECTION="sqlite1"

# --- SQLite Connections ---
DB_SQLITE_ENABLED=true
DB_SQLITE_1_PATH="storages/database/sqlite1.dev.db"
DB_SQLITE_2_PATH="storages/database/sqlite2.dev.db"
DB_SQLITE_3_PATH="storages/database/sqlite3.dev.db"

# --- PostgreSQL Connections ---
DB_PG_ENABLED=false
DB_PG_1_DSN="postgres://postgres:password@localhost:5432/db_main?sslmode=disable"
DB_PG_2_DSN="postgres://postgres:password@localhost:5432/db_analytics?sslmode=disable"
DB_PG_3_DSN="postgres://postgres:password@localhost:5432/db_tenant?sslmode=disable"

# --- MySQL Connections ---
DB_MYSQL_ENABLED=false
DB_MYSQL_1_DSN="root:password@tcp(localhost:3306)/db_store?parseTime=true"
DB_MYSQL_2_DSN="root:password@tcp(localhost:3306)/db_logs?parseTime=true"
DB_MYSQL_3_DSN="root:password@tcp(localhost:3306)/db_archive?parseTime=true"

# --- MongoDB Connections ---
DB_MONGO_ENABLED=false
DB_MONGO_1_URI="mongodb://localhost:27017/db_documents"
DB_MONGO_2_URI="mongodb://localhost:27017/db_events"
DB_MONGO_3_URI="mongodb://localhost:27017/db_metrics"
```

---

## 🛠️ 2. Desain Database Manager di Go (`pkg/database/manager.go`)

Karena `uptrace/bun` mendukung PostgreSQL, MySQL, dan SQLite melalui interface yang sama (`*bun.DB`), semua database relasional dapat dikelola secara seragam:

```go
package database

import (
    "fmt"
    "sync"
    "github.com/uptrace/bun"
    "go.mongodb.org/mongo-driver/v2/mongo"
)

type DBManager struct {
    mu          sync.RWMutex
    defaultName string
    sqlDBs      map[string]*bun.DB
    mongoDBs    map[string]*mongo.Database
}

func NewManager(defaultName string) *DBManager {
    return &DBManager{
        defaultName: defaultName,
        sqlDBs:      make(map[string]*bun.DB),
        mongoDBs:    make(map[string]*mongo.Database),
    }
}

// Mengambil koneksi Bun (bisa sqlite1, mysql1, pg1, dll)
func (m *DBManager) SQL(name string) (*bun.DB, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    db, ok := m.sqlDBs[name]
    if !ok {
        return nil, fmt.Errorf("sql connection '%s' not found or disabled", name)
    }
    return db, nil
}

// Mengambil koneksi default
func (m *DBManager) DefaultSQL() *bun.DB {
    db, _ := m.SQL(m.defaultName)
    return db
}

// Mengambil koneksi MongoDB (mongo1, mongo2, mongo3)
func (m *DBManager) Mongo(name string) (*mongo.Database, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    db, ok := m.mongoDBs[name]
    if !ok {
        return nil, fmt.Errorf("mongo connection '%s' not found or disabled", name)
    }
    return db, nil
}
```

---

## 🚀 3. Fleksibilitas Injeksi ke Domain Repositories

Setiap repository dapat dihubungkan ke koneksi database mana pun secara eksplisit saat *Composition Root* (`internal/app/app.go`):

```go
// User data disimpan di PostgreSQL master (pg1)
userRepo := user.NewRepository(dbManager.SQL("pg1"))

// Category data disimpan di SQLite lokal (sqlite1)
categoryRepo := category.NewRepository(dbManager.SQL("sqlite1"))

// Data Event/Document disimpan di MongoDB (mongo1)
eventRepo := event.NewRepository(dbManager.Mongo("mongo1"))
```

---

## 📋 4. Migrasi Skema Multi-Database dengan Goose

Tool CLI `cmd/make/main.go` mendukung migrasi terarah untuk masing-masing koneksi:
```bash
# Migrasi ke SQLite 1
go run cmd/make/main.go migrate:up --db=sqlite1

# Migrasi ke PostgreSQL 1
go run cmd/make/main.go migrate:up --db=pg1

# Migrasi ke MySQL 1
go run cmd/make/main.go migrate:up --db=mysql1
```

Struktur folder migrasi ditata per target database:
```text
resources/migrations/
├── sqlite1/
│   └── 20261029_create_users_table.sql
├── pg1/
│   └── 20261029_create_orders_table.sql
└── mysql1/
    └── 20261029_create_products_table.sql
```

---

## 🌱 5. Seeder System (Laravel-Style)

Sama seperti starter kit eksisting (`scripts/seed.ts`), sistem seeder di Golang dirancang modular per target engine dengan dukungan transaksi aman (*idempotent*):

### Struktur Folder Seeder:
```text
database/seeders/
├── seeder.go                 # Interface & Registry Runner Seeder
├── sqlite/
│   └── crud_seeder.go        # Seeder khusus SQLite
├── pg/
│   └── user_seeder.go        # Seeder khusus PostgreSQL
├── mysql/
│   └── store_seeder.go       # Seeder khusus MySQL
└── mongo/
    └── event_seeder.go       # Seeder khusus MongoDB
```

### Kontrak Interface Seeder di Go:
```go
package seeders

import (
    "context"
    "github.com/uptrace/bun"
)

type SQLSeeder interface {
    Name() string
    Run(ctx context.Context, db *bun.DB) error
}

type MongoSeeder interface {
    Name() string
    Run(ctx context.Context, db *mongo.Database) error
}
```

### Perintah CLI Seeder:
```bash
# Menjalankan seluruh seeder aktif
go run cmd/make/main.go seed

# Menjalankan seeder khusus per jenis database
go run cmd/make/main.go seed:sqlite
go run cmd/make/main.go seed:pg
go run cmd/make/main.go seed:mysql
go run cmd/make/main.go seed:mongo

# Scaffold file seeder baru via CLI
go run cmd/make/main.go make:seeder sqlite UserSeeder
go run cmd/make/main.go make:seeder pg OrderSeeder
```

### Fitur Keamanan Seeder:
1. **Transactional Seeding:** Seluruh eksekusi seeder SQL dibungkus di dalam `db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error { ... })`. Jika satu data gagal, seluruh insert di-rollback secara otomatis.
2. **Idempotent Insertion:** Menggunakan klausa `On("CONFLICT (id) DO NOTHING")` agar seeder aman dijalankan berulang kali tanpa membuat duplicate key error.

---

## ⚖️ Pilihan Engine Akses Data (ORM vs Query Builder vs SQLC)

```mermaid
flowchart TD
    subgraph ServiceLayer ["internal/domain/example/service.go"]
        Svc["CrudService"]
    end

    subgraph RepoInterface ["internal/domain/example/repository.go"]
        Repo["CrudRepository Interface"]
    end

    subgraph RepoImpl ["internal/repository/sqlite/ (or postgres)"]
        Impl["crud_repository.go"]
    end

    subgraph DataEngine ["Pilihan Engine"]
        EngineA["Opsi A: sqlc (Compiled Type-Safe SQL)"]
        EngineB["Opsi B: GORM (Eloquent-like ORM)"]
        EngineC["Opsi C: uptrace/bun (Fluent Query Builder)"]
    end

    Svc --> Repo
    Impl -.->|Implements| Repo
    Impl --> DataEngine
```

### 🔹 Opsi A: `sqlc` + `jackc/pgx` / `go-sqlite3` (Filosofi Identik dengan Drizzle ORM)
- **Cara Kerja:** Developer menulis SQL mentah di file `.sql`. Tool `sqlc` mengompilasi query tersebut menjadi kode Go struct dan method yang 100% type-safe.
- **Kelebihan:**
  - **Zero Reflection & Performa Tercepat:** Kinerja setara dengan memanggil driver database langsung.
  - **Compile-Time Safety:** Typo pada nama kolom atau kesalahan tipe data SQL langsung terdeteksi saat proses build.
- **Kekurangan:** Penulisan query dinamis dengan filter yang sangat banyak membutuhkan teknik SQL khusus (`COALESCE` / CASE).

### 🔹 Opsi B: GORM (`gorm.io/gorm`) (Familiaritas Tertinggi dengan Laravel Eloquent)
- **Cara Kerja:** Menggunakan struct tags untuk relasi dan skema, query ditulis menggunakan method chaining (`db.WithContext(ctx).Where("status = ?", "active").Find(&records)`).
- **Kelebihan:**
  - Sangat mirip dengan Laravel Eloquent (Hooks, Soft Deletes, Preload relasi).
  - Multi-database switching sangat mudah (cukup ganti package dialector: `gorm.io/driver/sqlite` atau `gorm.io/driver/postgres`).
  - Sangat mudah menangani query dinamis.
- **Kekurangan:** Menggunakan reflection runtime; kesalahan penulisan field seringkali baru terdeteksi saat runtime.

### 🔹 Opsi C: `uptrace/bun` (Jalan Tengah Modern)
- **Cara Kerja:** Query builder SQL modern yang cepat, elegan, dan mempertahankan kendali penuh atas SQL tanpa overhead reflection sebesar GORM.
- **Kelebihan:** Sangat cocok untuk developer yang menginginkan sintaks fluent SQL yang aman, cepat, dan transparan.

---

## 🔄 Standar Desain Repository & Transaksi di Go

### 1. Kontrak Interface Repository (`internal/domain/example/repository.go`):
```go
package example

import (
    "context"
)

type CrudRepository interface {
    Create(ctx context.Context, item *Crud) error
    FindByID(ctx context.Context, id string) (*Crud, error)
    FindAll(ctx context.Context, filter Filter) ([]Crud, int64, error)
    Update(ctx context.Context, item *Crud) error
    Delete(ctx context.Context, id string) error
    WithTx(ctx context.Context, fn func(repo CrudRepository) error) error
}
```

### 2. Pola Transaksi Bersih:
Service dapat menjalankan transaksi multi-tabel tanpa mengetahui detail `*sql.Tx`:
```go
func (s *CrudService) ProcessTransfer(ctx context.Context, fromID, toID string, amount int64) error {
    return s.repo.WithTx(ctx, func(txRepo CrudRepository) error {
        // Operasi 1
        if err := txRepo.Deduct(ctx, fromID, amount); err != nil {
            return err // Otomatis ROLLBACK
        }
        // Operasi 2
        if err := txRepo.Credit(ctx, toID, amount); err != nil {
            return err // Otomatis ROLLBACK
        }
        return nil // Otomatis COMMIT
    })
}
```

---

## 🛠️ Tool Migrasi Skema: `pressly/goose`
- Standar industri terpopuler di Go.
- Mendukung file migrasi `.sql` standar:
  ```sql
  -- +goose Up
  CREATE TABLE cruds (
      id TEXT PRIMARY KEY,
      title TEXT NOT NULL,
      description TEXT,
      created_at TIMESTAMP NOT NULL,
      updated_at TIMESTAMP NOT NULL
  );

  -- +goose Down
  DROP TABLE cruds;
  ```
- Dapat dijalankan via CLI starter kit: `go run cmd/make/main.go migrate:up`.

---

## 📋 Rekomendasi Solusi & Bahan Diskusi:

| Kriteria | Opsi A (sqlc) | Opsi B (GORM) | Opsi C (uptrace/bun) |
|---|---|---|---|
| **Kemiripan dengan Laravel Eloquent** | ⭐⭐⭐ (SQL murni) | ⭐⭐⭐⭐⭐ (Sangat Mirip) | ⭐⭐⭐⭐ (Query Builder) |
| **Keamanan Tipe & Efisiensi** | ⭐⭐⭐⭐⭐ (Compile-time) | ⭐⭐⭐ (Runtime reflection) | ⭐⭐⭐⭐ (Sangat Cepat) |
| **Kemudahan Pemula Go** | ⭐⭐⭐⭐ (Tahu SQL biasa) | ⭐⭐⭐⭐⭐ (Sangat Cepat Adaptasi) | ⭐⭐⭐⭐ (Mudah Dipahami) |

> 💬 **Bahan Diskusi Topik 2:**
> Apakah Anda lebih condong ke **Opsi B (GORM)** agar developer yang terbiasa dengan Eloquent merasa *"at home"*, atau **Opsi A (sqlc)** / **Opsi C (bun)** untuk performa maksimal dan kueri SQL yang eksplisit?
