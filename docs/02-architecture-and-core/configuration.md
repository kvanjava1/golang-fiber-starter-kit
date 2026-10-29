# 02. Manajemen Konfigurasi (Type-Safe with cleanenv)

Pengelolaan konfigurasi aplikasi ditangani oleh pustaka [`ilyakaznacheev/cleanenv`](https://github.com/ilyakaznacheev/cleanenv) yang berlokasi di [`internal/config/config.go`](file:///home/melodelavic/Documents/golang/golang-fiber-starter-kit/internal/config/config.go).

---

## 🛡️ Keunggulan Pendekatan Ini

1. **Type-Safe:** Variabel lingkungan dipetakan langsung ke tipe data Go yang tepat (`int`, `string`, `bool`, `time.Duration`).
2. **Default Fallback:** Nilai default dapat ditentukan langsung pada tag struct menggunakan `env-default:"..."`. Jika variabel tidak ada di file `.env`, sistem tidak akan crash.
3. **Validasi Wajib (`env-required:"true"`):** Variabel penting dapat ditandai wajib sehingga server akan gagal menyala (*fail fast*) jika kunci kritis tidak ditemukan.
4. **Hierarki Prioritas:** Variabel lingkungan OS (*environment variables*) memiliki prioritas lebih tinggi daripada nilai di file `.env`, mempermudah konfigurasi di lingkungan container Docker/Kubernetes.

---

## 🧩 Anatomi Struct Konfigurasi

Berikut adalah cuplikan deklarasi konfigurasi di `internal/config/config.go`:

```go
type Config struct {
    // Pengaturan Server
    AppName  string `env:"APP_NAME" env-default:"Go Starter Kit"`
    AppEnv   string `env:"APP_ENV" env-default:"local"`
    Port     int    `env:"PORT" env-default:"3000"`
    LogLevel string `env:"LOG_LEVEL" env-default:"info"`

    // Database Driver
    DBDefaultConnection string `env:"DB_DEFAULT_CONNECTION" env-default:"sqlite1"`
    DBSqliteEnabled     bool   `env:"DB_SQLITE_ENABLED" env-default:"true"`
    DBSqlite1Path       string `env:"DB_SQLITE1_PATH" env-default:"storages/database/sqlite1.db"`

    // Redis & Cache
    RedisEnabled bool   `env:"REDIS_ENABLED" env-default:"false"`
    RedisURL     string `env:"REDIS_URL" env-default:"redis://localhost:6379/0"`

    // Mailer
    MailDriver      string `env:"MAIL_DRIVER" env-default:"log"`
    MailFromAddress string `env:"MAIL_FROM_ADDRESS" env-default:"no-reply@starter.kit"`
    MailFromName    string `env:"MAIL_FROM_NAME" env-default:"Go Starter Kit"`
}
```

---

## 📖 Cara Menggunakan Konfigurasi di Kode

Untuk memuat konfigurasi saat aplikasi dimulai:

```go
package main

import (
    "log/slog"
    "golang-starter-kit/internal/config"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        slog.Error("Gagal memuat konfigurasi", "error", err)
        return
    }

    // Mengakses properti konfigurasi
    slog.Info("Server dimulai", "app", cfg.AppName, "port", cfg.Port)

    // Helper pengecekan environment
    if cfg.IsProduction() {
        // Logika khusus produksi
    }
}
```

---

## 🧪 Helper Method Bawaan

Struct `Config` dilengkapi method pembantu untuk mempermudah percabangan logika environment:

```go
// Memeriksa apakah aplikasi berjalan di environment production
func (c *Config) IsProduction() bool {
    return c.AppEnv == "production" || c.AppEnv == "prod"
}

// Memeriksa apakah aplikasi berjalan di environment pengujian
func (c *Config) IsTesting() bool {
    return c.AppEnv == "testing" || c.AppEnv == "test"
}
```
