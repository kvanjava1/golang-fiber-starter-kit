# 03. Transaksional Database Seeders

Pengisian data awal (*seeding*) dikelola oleh runner transaksional di direktori [`database/seeders/`](file:///home/melodelavic/Documents/golang/golang-fiber-starter-kit/database/seeders/runner.go).

---

## 🔒 Jaminan Keamanan Transaksi (`bun.IDB`)

Berbeda dengan runner seeder biasa, seeder pada starter kit ini menerima antarmuka `bun.IDB` (bukan sekadar pointer connection pool). Setiap seeder dijalankan di dalam blok transaksi `db.RunInTx(ctx, nil, ...)`:
- Jika seluruh operasi seeder berhasil, transaksi di-**commit**.
- Jika salah satu seeder mengalami error atau crash, seluruh data yang sudah dimasukkan otomatis di-**rollback**, menjaga database tetap bersih dan konsisten.

---

## ✍️ Cara Menulis Seeder Baru

1. Buat file seeder menggunakan CLI generator:
   ```bash
   go run cmd/make/main.go make:seeder user
   ```
2. Buka file yang dihasilkan di `database/seeders/user_seeder.go` dan implementasikan interface `Seeder`:

```go
package seeders

import (
    "context"
    "time"
    "github.com/uptrace/bun"
)

type UserSeeder struct{}

func (s *UserSeeder) Name() string {
    return "UserSeeder"
}

func (s *UserSeeder) Run(ctx context.Context, db bun.IDB) error {
    users := []map[string]interface{}{
        {
            "id":         "usr-001",
            "name":       "Admin Utama",
            "email":      "admin@starter.kit",
            "role":       "admin",
            "created_at": time.Now(),
            "updated_at": time.Now(),
        },
    }

    _, err := db.NewInsert().Model(&users).TableExpr("users").Exec(ctx)
    return err
}
```

3. Daftarkan seeder Anda ke dalam array registrasi di `database/seeders/runner.go`:

```go
func DefaultSeeders() []Seeder {
    return []Seeder{
        &UserSeeder{},
        // daftarkan seeder lainnya di sini
    }
}
```
