# 🛠️ Step 01: Inisialisasi Modul, Tata Letak Direktori & Konfigurasi Fail-Fast

Dokumen ini memuat panduan implementasi inisialisasi modul Go, pembuatan tata letak direktori standar industri, alokasi folder `storages/`, serta sistem konfigurasi `cleanenv`.

---

## 🎯 Tujuan Langkah
1. Menginisialisasi `go.mod` dengan Go versi 1.22+.
2. Membuat kerangka folder *Standard Go Project Layout*: `cmd/`, `internal/`, `pkg/`, `resources/`, `database/`, dan `storages/`.
3. Mengonfigurasi `internal/config/config.go` dengan `ilyakaznacheev/cleanenv` untuk validasi *fail-fast* variabel `.env`.
4. Menyediakan template `.env.example` lengkap untuk multi-database, queue, mail, dan app settings.

---

## 📂 Struktur File yang Dibuat
```text
golang-starter-kit/
├── go.mod
├── .env.example
├── .gitignore
├── storages/
│   ├── database/.gitkeep
│   ├── logs/app/.gitkeep
│   ├── logs/cli/.gitkeep
│   └── uploads/.gitkeep
└── internal/
    └── config/
        ├── config.go
        └── config_test.go
```

---

## 📝 Detail Spesifikasi Kode

### 1. File `.gitignore`
```gitignore
# Binaries
/bin/
/dist/
server
worker
make

# Local Storages
storages/database/*.db*
storages/logs/*/*.log
storages/uploads/*
!storages/**/.gitkeep

# Environment
.env
.env.local

# IDE
.idea/
.vscode/
*.swp
```

### 2. Struct Konfigurasi (`internal/config/config.go`)
```go
package config

import (
    "fmt"
    "github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
    AppEnv      string `env:"APP_ENV" env-default:"development"`
    AppName     string `env:"APP_NAME" env-default:"Go Starter Kit"`
    Port        int    `env:"PORT" env-default:"3000"`
    LogLevel    string `env:"LOG_LEVEL" env-default:"info"`

    // Database Master Relational
    DBDefaultConnection string `env:"DB_DEFAULT_CONNECTION" env-default:"sqlite1"`

    // SQLite Connections
    DBSqliteEnabled bool   `env:"DB_SQLITE_ENABLED" env-default:"true"`
    DBSqlite1Path   string `env:"DB_SQLITE_1_PATH" env-default:"storages/database/sqlite1.dev.db"`
    DBSqlite2Path   string `env:"DB_SQLITE_2_PATH" env-default:"storages/database/sqlite2.dev.db"`
    DBSqlite3Path   string `env:"DB_SQLITE_3_PATH" env-default:"storages/database/sqlite3.dev.db"`

    // PostgreSQL Connections
    DBPGEnabled bool   `env:"DB_PG_ENABLED" env-default:"false"`
    DBPG1DSN    string `env:"DB_PG_1_DSN" env-default:""`
    DBPG2DSN    string `env:"DB_PG_2_DSN" env-default:""`
    DBPG3DSN    string `env:"DB_PG_3_DSN" env-default:""`

    // MySQL Connections
    DBMySQLEnabled bool   `env:"DB_MYSQL_ENABLED" env-default:"false"`
    DBMySQL1DSN    string `env:"DB_MYSQL_1_DSN" env-default:""`
    DBMySQL2DSN    string `env:"DB_MYSQL_2_DSN" env-default:""`
    DBMySQL3DSN    string `env:"DB_MYSQL_3_DSN" env-default:""`

    // MongoDB Connections
    DBMongoEnabled bool   `env:"DB_MONGO_ENABLED" env-default:"false"`
    DBMongo1URI    string `env:"DB_MONGO_1_URI" env-default:""`
    DBMongo2URI    string `env:"DB_MONGO_2_URI" env-default:""`
    DBMongo3URI    string `env:"DB_MONGO_3_URI" env-default:""`

    // Redis
    RedisURL string `env:"REDIS_URL" env-default:"redis://localhost:6379"`

    // Mail System
    MailDriver    string `env:"MAIL_DRIVER" env-default:"log"` // log, smtp, resend
    MailFromName  string `env:"MAIL_FROM_NAME" env-default:"Go Starter Kit"`
    MailFromEmail string `env:"MAIL_FROM_EMAIL" env-default:"noreply@example.com"`
    MailSMTPHost  string `env:"MAIL_SMTP_HOST" env-default:"localhost"`
    MailSMTPPort  int    `env:"MAIL_SMTP_PORT" env-default:"1025"`
    MailSMTPUser  string `env:"MAIL_SMTP_USER" env-default:""`
    MailSMTPPass  string `env:"MAIL_SMTP_PASS" env-default:""`
    MailResendKey string `env:"MAIL_RESEND_API_KEY" env-default:""`
}

func Load() (*Config, error) {
    var cfg Config
    // Baca dari file .env jika ada, fallback ke env OS
    if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
        if err := cleanenv.ReadEnv(&cfg); err != nil {
            return nil, fmt.Errorf("failed to load configuration: %w", err)
        }
    }
    return &cfg, nil
}
```

---

## ✅ Kriteria Verifikasi
- [ ] `go mod tidy` berhasil mengunduh dependensi `cleanenv`.
- [ ] Menjalankan unit test `internal/config/config_test.go` sukses memverifikasi default value dan parsing file `.env`.
- [ ] Struktur folder `storages/` terlindungi dari komit git via `.gitignore`.
