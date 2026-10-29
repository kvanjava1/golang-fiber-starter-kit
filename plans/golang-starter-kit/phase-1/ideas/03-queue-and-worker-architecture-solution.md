# 💡 Opsi Solusi 03: Queue System, Worker Lifecycle & Tracing

> **Status Keputusan:** ✅ **DECIDED: Opsi A (`hibiken/asynq` + Redis)**  
> **Keputusan:** Dipilih **`hibiken/asynq`** dengan Redis broker sebagai sistem antrean asinkron utama. Memberikan kesetaraan fitur 1:1 dengan BullMQ (auto-retries exponential backoff, delayed scheduling, priority queues, tracing context propagation), serta dashboard web pemantau antrean resmi (`asynqmon`).

---

## 🎯 Standar Production-Grade yang Harus Dipenuhi
1. **Persistent Task Broker:** Tugas kritis (seperti email atau pemrosesan pembayaran) tidak boleh hilang saat server mati mendadak (*crash-resilient*).
2. **Context & Correlation ID Propagation:** ID pelacakan (`X-Correlation-ID`) dari request HTTP wajib dimasukkan ke dalam metadata payload task, dan dihidupkan kembali ke dalam `ctx context.Context` saat task dieksekusi oleh worker.
3. **Graceful Draining:** Saat menerima sinyal OS (`SIGINT`, `SIGTERM`), worker berhenti mengambil tugas baru dari Redis dan memberikan waktu tenggang (*timeout*) bagi tugas yang sedang berjalan untuk selesai secara normal.
4. **Retry Eksponensial & Dead-Letter Queue (DLQ):** Otomatis mencoba ulang tugas yang gagal dengan jeda waktu eksponensial (misal: 1 detik, 2 detik, 4 detik, dst.) hingga batas maksimum percobaan, lalu memindahkannya ke antrean gagal (*archived/dead-letter*) untuk diaudit.

---

## ⚖️ Pilihan Engine Antrean Asinkron

```mermaid
sequenceDiagram
    autonumber
    actor Client as HTTP Client
    participant API as HTTP Server (cmd/server)
    participant Redis as Redis Broker
    participant Worker as Background Worker (cmd/worker)
    participant Svc as Domain Logic

    Client->>API: POST /api/example/cruds/send-email (Header: X-Correlation-ID=abc)
    API->>API: Ekstrak Correlation ID & Encode Payload JSON
    API->>Redis: asynq.Enqueue(Task: "email:send", Payload + Meta, MaxRetry: 3)
    Redis-->>API: Task Enqueued
    API-->>Client: HTTP 202 Accepted { success: true, message: "Email queued" }
    
    Note over Redis,Worker: Asynchronous Processing (di luar HTTP cycle)
    Redis->>Worker: Dequeue Task
    Worker->>Worker: Inject Correlation ID ke dalam worker context.Context
    Worker->>Svc: Execute SendEmail(ctx, payload)
    alt Berhasil
        Svc-->>Worker: nil
        Worker->>Redis: Mark Task Complete
    alt Gagal (Error)
        Svc-->>Worker: error
        Worker->>Redis: Schedule Retry with Exponential Backoff
    end
```

### 🔹 Opsi A: `hibiken/asynq` + Redis (Klon 1:1 BullMQ di Go - SANGAT DIREKOMENDASIKAN)
- **Karakteristik:** Library antrean terdistribusi berbasis Redis paling populer dan tangguh di Go.
- **Kelebihan:**
  - Fitur 100% setara dengan BullMQ: Retries bertingkat, delayed scheduling (misal kirim email 1 jam lagi), priority queues (`critical`, `default`, `low`), task deduplication.
  - Menyediakan Web Dashboard resmi (`asynqmon`) yang siap dipasang untuk memantau throughput antrean secara visual.
  - Sangat efisien: Mampu memproses puluhan ribu job per detik dengan jejak memori sangat kecil di Go.
- **Kekurangan:** Memerlukan Redis aktif sebagai broker.

### 🔹 Opsi B: In-Memory Channel Worker Pool (Zero-Dependency)
- **Karakteristik:** Menggunakan buffered channel bawaan Go (`chan Task`) dan goroutine worker pool.
- **Kelebihan:** Tidak butuh Redis sama sekali, cocok untuk kebutuhan embedded atau deployment single-server yang sangat minim sumber daya.
- **Kekurangan:** Antrean hilang jika aplikasi di-restart, tidak mendukung multi-server.

### 🔹 Opsi C: Hybrid Driver Abstraction (`pkg/queue/`)
Menyediakan interface terpadu sehingga aplikasi dapat beralih antara driver `sync` (in-memory) saat testing/dev ringan dan driver `asynq` (Redis) saat staging/production via `.env`.

---

## 🧬 Pola Propagasi Tracing & Context di Worker

Untuk menjamin observabilitas production-grade, metadata context dimasukkan ke dalam header task:

### 1. Payload Task Wrapper (`pkg/queue/task.go`):
```go
package queue

type TaskHeader struct {
    CorrelationID string `json:"correlation_id"`
    UserID        string `json:"user_id,omitempty"`
    EnqueuedAt    int64  `json:"enqueued_at"`
}

type EmailTaskPayload struct {
    Header TaskHeader `json:"header"`
    To     string     `json:"to"`
    HTML   string     `json:"html"`
}
```

### 2. Context Reconstruction di Sisi Worker:
```go
func HandleSendEmailTask(ctx context.Context, t *asynq.Task) error {
    var p EmailTaskPayload
    if err := json.Unmarshal(t.Payload(), &p); err != nil {
        return fmt.Errorf("failed to parse task payload: %w", err)
    }

    // Bangun context baru yang membawa CorrelationID dari request awal
    ctxWithTrace := context.WithValue(ctx, "correlation_id", p.Header.CorrelationID)
    
    // Semua log di bawah ini otomatis mencatat CorrelationID yang sama
    slog.InfoContext(ctxWithTrace, "Processing email task", "to", p.To)
    
    return emailService.SendRaw(ctxWithTrace, p.To, p.HTML)
}
```

---

## 🏗️ Topologi Eksekusi (Development vs Production)

Untuk kemudahan Developer Experience (DX):
1. **Mode Development:** Perintah `go run cmd/server/main.go --with-worker` dapat menjalankan HTTP server dan background worker dalam satu proses yang sama.
2. **Mode Production:**
   - Container 1 (Web Pod): Menjalankan `cmd/server/main.go` (Khusus menerima HTTP traffic).
   - Container 2 (Worker Pod): Menjalankan `cmd/worker/main.go` (Khusus menyedot antrean Redis).

---

## 📋 Rekomendasi Solusi & Bahan Diskusi:

> 💬 **Bahan Diskusi Topik 3:**
> Apakah Anda setuju menggunakan **`hibiken/asynq`** sebagai broker antrean utama karena kemampuannya yang setara 1:1 dengan BullMQ, dengan opsi menjalankan HTTP dan Worker dalam satu proses saat lokal dev agar tidak merepotkan developer?
