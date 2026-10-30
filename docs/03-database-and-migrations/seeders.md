# 03. Multi-Database & Transaksional Database Seeders

Pengisian data awal (*seeding*) dikelola secara terisolasi per target koneksi database di direktori [`database/seeders/`](file:///home/melodelavic/Documents/golang/golang-fiber-starter-kit/database/seeders/).

---

## 🎯 Terisolasi Penuh per Database Target

Setiap database target (`sqlite1`, `mysql1`, `pg1`, `mongo1`) memiliki sub-direktori seeder masing-masing:
```text
database/seeders/
├── seeder.go                 ← Interface Seeder, MongoSeeder, & Runner transaksional
├── sqlite1/                  ← Seeder KHUSUS SQLite 1
│   └── user_seeder.go
├── mysql1/                   ← Seeder KHUSUS MySQL 1
│   └── seeder.go
├── pg1/                      ← Seeder KHUSUS PostgreSQL 1
│   └── seeder.go
└── mongo1/                   ← Seeder KHUSUS MongoDB 1
    └── seeder.go
```

### 🔒 Jaminan Keamanan Isolasi
Ketika Anda menjalankan seeder untuk `mysql1`:
- **HANYA** koneksi MySQL yang dibuka.
- **HANYA** seeder yang terdaftar di `mysql1` yang dieksekusi.
- Database **SQLite, PostgreSQL, dan MongoDB sama sekali TIDAK disentuh, TIDAK dibuka, dan TIDAK ikut terkena seed**.

---

## 🔒 Jaminan Keamanan Transaksi (`bun.IDB`)

Setiap seeder SQL menerima antarmuka `bun.IDB` dan otomatis dibungkus di dalam transaksi atomik `db.RunInTx(ctx, nil, ...)`:
- Jika seluruh operasi seeder sukses, transaksi di-**commit**.
- Jika salah satu seeder mengalami error atau crash, seluruh perubahan otomatis di-**rollback**, menjaga database tetap bersih dan konsisten.

---

## 🛠️ Cara Menjalankan Seeder (`db:seed`)

```bash
# 1. Jalankan seeder SQLite saja (default: sqlite1)
make seed
# Atau: go run cmd/make/main.go db:seed --db=sqlite1

# 2. Jalankan seeder MySQL saja (sqlite & db lain tidak disentuh)
make seed-mysql
# Atau: go run cmd/make/main.go db:seed --db=mysql1

# 3. Jalankan seeder PostgreSQL saja
make seed-pg
# Atau: go run cmd/make/main.go db:seed --db=pg1

# 4. Jalankan seeder MongoDB saja
go run cmd/make/main.go db:seed --db=mongo1
```

---

## ✍️ Cara Menghasilkan Seeder Baru (`make:seeder`)

Gunakan CLI Scaffolder dengan menyertakan flag `--db=<target>`:

```bash
# Membuat seeder untuk MySQL
go run cmd/make/main.go make:seeder customer --db=mysql1
# File terbuat di: database/seeders/mysql1/customer_seeder.go

# Membuat seeder untuk PostgreSQL
go run cmd/make/main.go make:seeder transaction --db=pg1
# File terbuat di: database/seeders/pg1/transaction_seeder.go

# Membuat seeder untuk MongoDB
go run cmd/make/main.go make:seeder log --db=mongo1
# File terbuat di: database/seeders/mongo1/log_seeder.go

# Membuat seeder untuk SQLite (default)
go run cmd/make/main.go make:seeder role
# File terbuat di: database/seeders/sqlite1/role_seeder.go
```

---

## 📝 Contoh Implementasi Seeder

### Seeder SQL (MySQL / PostgreSQL / SQLite)
```go
package mysql1

import (
    "context"
    "time"
    "github.com/uptrace/bun"
    "golang-starter-kit/database/seeders"
)

func init() {
    seeders.Register("mysql1", NewCustomerSeeder())
}

type CustomerSeeder struct{}

func NewCustomerSeeder() *CustomerSeeder {
    return &CustomerSeeder{}
}

func (s *CustomerSeeder) Name() string {
    return "mysql1:CustomerSeeder"
}

func (s *CustomerSeeder) Run(ctx context.Context, db bun.IDB) error {
    customers := []map[string]any{
        {"name": "Customer PT ABC", "status": "active", "created_at": time.Now()},
    }
    _, err := db.NewInsert().Model(&customers).TableExpr("customers").Exec(ctx)
    return err
}
```

### Seeder NoSQL (MongoDB)
```go
package mongo1

import (
    "context"
    "time"
    "go.mongodb.org/mongo-driver/v2/mongo"
    "golang-starter-kit/database/seeders"
)

func init() {
    seeders.RegisterMongo("mongo1", NewActivityLogSeeder())
}

type ActivityLogSeeder struct{}

func NewActivityLogSeeder() *ActivityLogSeeder {
    return &ActivityLogSeeder{}
}

func (s *ActivityLogSeeder) Name() string {
    return "mongo1:ActivityLogSeeder"
}

func (s *ActivityLogSeeder) Run(ctx context.Context, db *mongo.Database) error {
    logs := []any{
        map[string]any{"action": "system_init", "timestamp": time.Now()},
    }
    _, err := db.Collection("activity_logs").InsertMany(ctx, logs)
    return err
}
```
