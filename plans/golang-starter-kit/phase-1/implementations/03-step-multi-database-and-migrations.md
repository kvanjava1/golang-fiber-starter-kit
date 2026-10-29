# 🛠️ Step 03: Multi-Database Manager, Migrasi Goose & Seeders

Dokumen ini memuat panduan implementasi konektivitas multi-database (SQLite 1..3, PostgreSQL 1..3, MySQL 1..3, MongoDB 1..3) menggunakan **`uptrace/bun`** dan official **`mongo-driver`**, integrasi migrasi **`pressly/goose`**, serta sistem seeder transaksional.

---

## 🎯 Tujuan Langkah
1. Mengimplementasikan `pkg/database/manager.go` untuk mengelola pool koneksi dinamik ke SQLite, PostgreSQL, MySQL, dan MongoDB.
2. Mengintegrasikan `pressly/goose` untuk migrasi skema SQL terarah per database target (`resources/migrations/`).
3. Mengimplementasikan runner seeder transaksional & *idempotent* di `database/seeders/`.

---

## 📂 Struktur File yang Dibuat
```text
pkg/
└── database/
    ├── manager.go            # Registry multi-koneksi bun.DB & mongo.Database
    ├── sqlite.go             # Inisialisasi driver SQLite (mattn/go-sqlite3)
    ├── postgres.go           # Inisialisasi driver PostgreSQL (pgdriver)
    ├── mysql.go              # Inisialisasi driver MySQL (mysqldriver)
    └── mongo.go              # Inisialisasi official mongo-driver
resources/
└── migrations/               # Direktori migrasi goose
    ├── sqlite1/
    ├── pg1/
    └── mysql1/
database/
└── seeders/
    ├── seeder.go             # Interface seeder & runner
    └── sqlite/
        └── example_seeder.go # Seeder contoh
```

---

## 📝 Detail Spesifikasi Kode

### 1. Inisialisasi SQLite (`pkg/database/sqlite.go`)
```go
package database

import (
    "database/sql"
    "os"
    "path/filepath"
    "github.com/uptrace/bun"
    "github.com/uptrace/bun/dialect/sqlitedialect"
    "github.com/uptrace/bun/driver/sqliteshim" // Pure Go fallback (modernc.org/sqlite) saat CGO nonaktif
)

func ConnectSQLite(dbPath string) (*bun.DB, error) {
    if dbPath != ":memory:" {
        dir := filepath.Dir(dbPath)
        if err := os.MkdirAll(dir, 0755); err != nil {
            return nil, err
        }
    }

    sqldb, err := sql.Open(sqliteshim.ShimName, dbPath+"?_journal=WAL&_timeout=10000")
    if err != nil {
        return nil, err
    }

    // Set connection pool limits
    sqldb.SetMaxOpenConns(1) // SQLite single-writer optimization
    sqldb.SetMaxIdleConns(1)

    db := bun.NewDB(sqldb, sqlitedialect.New())
    return db, nil
}
```

### 2. Multi-Connection Manager (`pkg/database/manager.go`)
```go
package database

import (
    "context"
    "fmt"
    "sync"
    "github.com/uptrace/bun"
    "go.mongodb.org/mongo-driver/v2/mongo"
)

type DBManager struct {
    mu           sync.RWMutex
    defaultName  string
    sqlConns     map[string]*bun.DB
    mongoConns   map[string]*mongo.Database
    mongoClients map[string]*mongo.Client
}

func NewManager(defaultName string) *DBManager {
    return &DBManager{
        defaultName:  defaultName,
        sqlConns:     make(map[string]*bun.DB),
        mongoConns:   make(map[string]*mongo.Database),
        mongoClients: make(map[string]*mongo.Client),
    }
}

func (m *DBManager) RegisterSQL(name string, db *bun.DB) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.sqlConns[name] = db
}

func (m *DBManager) RegisterMongo(name string, client *mongo.Client, db *mongo.Database) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.mongoClients[name] = client
    m.mongoConns[name] = db
}

func (m *DBManager) SQL(name string) (*bun.DB, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    db, ok := m.sqlConns[name]
    if !ok {
        return nil, fmt.Errorf("sql connection '%s' not found", name)
    }
    return db, nil
}

func (m *DBManager) DefaultSQL() *bun.DB {
    db, _ := m.SQL(m.defaultName)
    return db
}

func (m *DBManager) Mongo(name string) (*mongo.Database, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    db, ok := m.mongoConns[name]
    if !ok {
        return nil, fmt.Errorf("mongo connection '%s' not found", name)
    }
    return db, nil
}

// CloseAll menutup seluruh pool koneksi SQL dan mendiskoneksikan client MongoDB secara aman
func (m *DBManager) CloseAll(ctx context.Context) {
    m.mu.Lock()
    defer m.mu.Unlock()

    for _, db := range m.sqlConns {
        _ = db.Close()
    }

    for _, client := range m.mongoClients {
        _ = client.Disconnect(ctx)
    }
}
```

### 3. Transactional Seeder Runner (`database/seeders/seeder.go`)
```go
package seeders

import (
    "context"
    "fmt"
    "log/slog"
    "github.com/uptrace/bun"
)

// Seeder menerima bun.IDB sehingga kompatibel dengan *bun.DB maupun bun.Tx
type Seeder interface {
    Name() string
    Run(ctx context.Context, db bun.IDB) error
}

type Runner struct {
    seeders []Seeder
}

func NewRunner() *Runner {
    return &Runner{}
}

func (r *Runner) Register(s Seeder) {
    r.seeders = append(r.seeders, s)
}

func (r *Runner) RunAll(ctx context.Context, db *bun.DB) error {
    for _, s := range r.seeders {
        slog.InfoContext(ctx, "Running seeder", "name", s.Name())
        // Jalankan di dalam transaksi aman dengan rollback otomatis jika return error
        err := db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
            return s.Run(ctx, tx) // Mengoper tx, bukan db!
        })
        if err != nil {
            return fmt.Errorf("seeder %s failed: %w", s.Name(), err)
        }
        slog.InfoContext(ctx, "Seeder completed", "name", s.Name())
    }
    return nil
}
```

---

## ✅ Kriteria Verifikasi
- [ ] Koneksi ke `sqlite1` berhasil membuat file `.dev.db` di `storages/database/`.
- [ ] Multiple connection (`sqlite1`, `sqlite2`) dapat diakses secara simultan tanpa file locking error.
- [ ] Perintah migrasi Goose berhasil mengeksekusi tabel baru.
- [ ] Seeder berhasil mengisi data dummy dan aman di-run berkali-kali (*idempotent*).
