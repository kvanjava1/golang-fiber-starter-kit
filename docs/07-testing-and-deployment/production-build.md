# 02. Build & Deployment Produksi

Panduan untuk mengompilasi binary mandiri (*single binary*) dan mendeploy aplikasi ke lingkungan server produksi.

---

## 📦 Single Binary Compilation (`make build`)

Karena menggunakan driver database pure-Go (`CGO_ENABLED=0`) dan template `a-h/templ` yang dikompilasi ke kode Go murni, aplikasi ini dapat dibangun menjadi file binary statis tunggal:

```bash
make build
```

Perintah di atas menghasilkan 3 binary mandiri di folder `bin/`:
- `bin/server` : Binary server HTTP Web & REST API (~32 MB)
- `bin/worker` : Binary background job worker daemon (~15 MB)
- `bin/make`   : Binary utilitas migrasi & CLI database (~20 MB)

Binary ini **tidak memiliki ketergantungan runtime ke sistem operasi** dan dapat langsung dipindahkan ke VPS / server Linux mana pun.

---

## 🐳 Contoh Dockerfile Multi-Stage

Berikut adalah template Dockerfile optimal dan ringan (~25 MB) menggunakan base image scratch atau alpine:

```dockerfile
# Stage 1: Build binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependensi sistem dan templ CLI
RUN apk add --no-cache make git
RUN go install github.com/a-h/templ/cmd/templ@latest

# Cache go mod
COPY go.mod go.sum ./
RUN go mod download

# Copy source code dan build
COPY . .
RUN make build

# Stage 2: Runtime image yang sangat ramping
FROM alpine:3.20

WORKDIR /app

# Install ca-certificates untuk HTTPS request & timezone
RUN apk add --no-cache ca-certificates tzdata

# Copy binary dan folder statis/storages
COPY --from=builder /app/bin/server /app/server
COPY --from=builder /app/bin/worker /app/worker
COPY --from=builder /app/bin/make /app/make
COPY --from=builder /app/public /app/public
COPY --from=builder /app/resources/migrations /app/resources/migrations

# Setup folder penyimpanan data
RUN mkdir -p /app/storages/database /app/storages/logs /app/storages/uploads

EXPOSE 3000

CMD ["/app/server"]
```

---

## 🔄 Otomasi CI/CD (GitHub Actions)

Starter kit ini sudah dilengkapi pipeline otomatis di [`.github/workflows/ci.yml`](file:///home/melodelavic/Documents/golang/golang-fiber-starter-kit/.github/workflows/ci.yml).

Setiap kali Anda melakukan `git push` atau membuat `Pull Request` ke branch `main`/`master`, GitHub Actions akan:
1. Memeriksa formatting kode Go (`gofmt`).
2. Menjalankan linter kode.
3. Mengompilasi template `templ`.
4. Menjalankan seluruh test suite dengan Race Detector (`go test -race ./...`).
5. Memverifikasi bahwa seluruh binary (`server`, `worker`, `make`) berhasil dikompilasi tanpa error.
