# 02. Migrasi Database (Goose)

Pengelolaan skema database ditangani oleh [`pressly/goose v3`](https://github.com/pressly/goose). File migrasi ditulis dalam format SQL murni dengan anotasi `-- +goose Up` dan `-- +goose Down`.

---

## 📁 Lokasi File Migrasi

Seluruh file migrasi disimpan di direktori [`resources/migrations/<dbTarget>/`](file:///home/melodelavic/Documents/golang/golang-fiber-starter-kit/resources/migrations/sqlite1/) (misalnya: `resources/migrations/sqlite1/`).

Contoh nama file migrasi:
```text
resources/migrations/sqlite1/20261029000001_create_cruds_table.sql
```

---

## 📝 Format Penulisan File Migrasi SQL

Setiap file migrasi wajib mendefinisikan bagian `Up` (eksekusi perubahan skema) dan `Down` (rollback perubahan skema):

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cruds (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_cruds_status ON cruds(status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cruds;
-- +goose StatementEnd
```

---

## 🛠️ Perintah Eksekusi Migrasi

### 1. Menjalankan Migrasi yang Tertunda
```bash
make migrate
# Atau langsung melalui CLI:
go run cmd/make/main.go migrate:up --db=sqlite1
```

### 2. Memeriksa Status Migrasi
```bash
make migrate-status
# Output menampilkan daftar file migrasi beserta status apakah sudah diterapkan (Applied) atau belum (Pending)
```

### 3. Membatalkan (Rollback) Migrasi Terakhir
```bash
make migrate-rollback
# Atau langsung melalui CLI:
go run cmd/make/main.go migrate:down --db=sqlite1
```

### 4. Membuat File Migrasi Baru
Gunakan perintah generator:
```bash
go run cmd/make/main.go make:migration create_orders_table
```
Perintah ini akan membuat file migrasi baru dengan prefix timestamp saat ini secara otomatis.
