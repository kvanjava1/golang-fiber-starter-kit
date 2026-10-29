# 📊 Analisis Referensi: Caching, Logging & Observability

Dokumen ini membedah subsistem caching terdistribusi, structured logging, metrik Prometheus, dan health checks pada starter kit eksisting.

---

## 1. Caching Utility & Invalidation (`src/utils/cache.util.ts`)

Starter kit menyediakan abstraksi caching seragam di atas Redis (menggunakan `ioredis`):

### Fitur Kunci Caching:
- **Metode Operasi Standar:** `get(key)`, `set(key, value, ttlSeconds)`, `del(key)`, `remember(key, ttlSeconds, fallbackFn)`.
- **Pola `remember` (Cache-Aside):** Mengambil data dari cache jika ada; jika tidak ada (*cache miss*), mengeksekusi fungsi `fallbackFn`, menyimpan hasilnya ke cache dengan TTL yang ditentukan, lalu mengembalikan data.
- **Tag-Based Invalidation:** Mengelompokkan key cache berdasarkan domain (misal: tag `cruds`), sehingga saat ada mutasi data (Create/Update/Delete), seluruh cache terkait dapat di-*flush* secara instan tanpa mengganggu cache domain lain.

---

## 2. Structured Logging dengan Correlation ID (`src/utils/logger.util.ts`)

Menggunakan **Pino** sebagai mesin logger berkecepatan tinggi:
- **Format JSON di Production:** Log dicetak dalam format JSON satu baris standar agar mudah di-*ingest* oleh Elasticsearch, Datadog, atau Grafana Loki.
- **Pretty Print di Development:** Menghasilkan output berwarna dan terbaca rapi oleh developer di terminal lokal.
- **Automatic Context Metadata:** Setiap baris log otomatis memuat `requestId`, `correlationId`, level log, timestamp ISO, dan nama service.

---

## 3. Prometheus Telemetry & Metrics (`src/utils/metrics.util.ts`)

Starter kit mengekspos metrik standar Prometheus pada endpoint `GET /metrics`:
- **`http_request_duration_seconds` (Histogram):** Mengukur distribusi latensi request HTTP berdasarkan label `method`, `route`, dan `status_code`.
- **`http_requests_total` (Counter):** Menghitung total akumulasi request HTTP yang masuk.
- **Runtime Metrics:** Otomatis mengumpulkan metrik runtime Bun/Node (memori heap, garbage collection, event loop lag, CPU usage).

---

## 4. Health Check Probes (`GET /api/health`)

Menyediakan endpoint observabilitas untuk container orchestration (Kubernetes liveness/readiness probe atau Docker healthcheck):
- Memeriksa status koneksi aktif ke Database (ping SQLite/Postgres/MySQL).
- Memeriksa status koneksi aktif ke Redis.
- Menampilkan waktu operasional (*uptime*) server dan timestamp server.
- Mengembalikan HTTP status `200 OK` jika semua dependensi sehat, atau `503 Service Unavailable` jika ada dependensi krusial yang mati.
