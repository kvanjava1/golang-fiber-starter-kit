# 04. Logging & Observability

Observabilitas dalam starter kit ini dibangun menggunakan paket standar Go 1.21+ [`log/slog`](file:///home/melodelavic/Documents/golang/golang-fiber-starter-kit/pkg/logger/logger.go) dengan penelusuran korelasi terpusat (*Trace Correlation ID*) dan metrik Prometheus.

---

## 🪵 Structured Logging (`log/slog`)

Seluruh log dicatat dalam format JSON terstruktur saat produksi (*production*) dan text berwarna saat pengembangan lokal (*local*).

### Fitur Kunci:
- **Trace Correlation:** Setiap request HTTP yang masuk otomatis diberikan `correlation_id` unik (atau meneruskan header `X-Correlation-ID` / `X-Request-ID` dari upstream API gateway).
- **Context-Aware:** Dengan menyematkan `ctx` saat melakukan logging, `correlation_id` otomatis muncul pada baris log:
  ```go
  slog.InfoContext(ctx, "Memproses pesanan pengguna", "order_id", orderID, "user_id", userID)
  ```
  Output log JSON yang dihasilkan:
  ```json
  {
    "time": "2026-10-29T15:00:00Z",
    "level": "INFO",
    "msg": "Memproses pesanan pengguna",
    "correlation_id": "a1b2c3d4-e5f6-7890-abcd-1234567890ab",
    "order_id": "ORD-9912",
    "user_id": "USR-100"
  }
  ```

---

## 📊 Metrik Telemetri Prometheus

Starter kit menyediakan endpoint bawaan di:
```text
GET http://localhost:3000/metrics
```
Endpoint ini siap di-scrape oleh server Prometheus dan divisualisasikan menggunakan Grafana untuk memantau:
- `http_requests_total`: Total hit per metode HTTP dan status code.
- `http_request_duration_seconds`: Histogram latensi request server.
- Metrik internal runtime Go (penggunaan goroutine, memori heap, frekuensi Garbage Collection).

---

## 🩺 Endpoint Pemantauan Kesehatan (Health Probes)

Endpoint kesehatan sistem berlokasi di:
```text
GET http://localhost:3000/api/v1/health
```

Contoh respons yang dihasilkan:
```json
{
  "success": true,
  "message": "Sistem berjalan normal",
  "data": {
    "status": "healthy",
    "app": "Go Starter Kit",
    "timestamp": "2026-10-29T15:30:00Z",
    "components": {
      "database": {
        "status": "up",
        "connection": "sqlite1"
      },
      "cache": {
        "status": "up",
        "driver": "memory"
      }
    }
  }
}
```

Endpoint ini sangat ideal digunakan sebagai *Readiness* dan *Liveness Probe* di Kubernetes atau Docker Compose healthcheck.
