# 🛠️ Step 06: Asynq Queue, Worker Process & Mail System

Dokumen ini memuat panduan implementasi antrean asinkron terdistribusi berbasis **`hibiken/asynq`**, proses background worker (`cmd/worker`), serta sistem email ala Laravel dengan template **`a-h/templ`** dan multi-driver transport.

---

## 🎯 Tujuan Langkah
1. Mengimplementasikan `pkg/queue/` untuk enqueue task ke Redis via Asynq dengan propagasi `correlation_id`.
2. Membangun entry point background worker di `cmd/worker/main.go` dengan *graceful draining*.
3. Mengimplementasikan `pkg/mail/` dengan interface `Mailable`, template email JSX-like via `a-h/templ`, driver (`log`, `smtp`, `resend`), dan facade fluent `mail.To().Queue()`.

---

## 📂 Struktur File yang Dibuat
```text
pkg/
├── queue/
│   ├── client.go         # Asynq client wrapper
│   ├── tasks.go          # Definisi nama task & payload
│   └── context.go        # Injeksi correlation_id ke task header
├── mail/
│   ├── mailable.go       # Interface Mailable & Envelope
│   ├── driver.go         # Interface Driver & factory (Log, SMTP, Resend)
│   ├── facade.go         # Fluent builder: mail.To().Queue() & .Send()
│   └── drivers/
│       ├── log_driver.go
│       ├── smtp_driver.go
│       └── resend_driver.go
cmd/
└── worker/
    └── main.go           # Asynq Worker server runner
resources/
└── views/
    └── emails/
        ├── layouts/
        │   └── email.templ
        └── welcome.templ
```

---

## 📝 Detail Spesifikasi Kode

### 1. Kontrak Mailable & Fluent Facade (`pkg/mail/`)
```go
package mail

import (
    "context"
    "fmt"
)

type Envelope struct {
    Subject string
    From    string
    To      []string
}

type Mailable interface {
    Envelope() Envelope
    Render(ctx context.Context) (string, error)
}

type PendingMail struct {
    client  *MailClient
    to      string
}

func (c *MailClient) To(recipient string) *PendingMail {
    return &PendingMail{client: c, to: recipient}
}

// Send: Langsung kirim via driver aktif secara sinkron
func (p *PendingMail) Send(ctx context.Context, m Mailable) error {
    html, err := m.Render(ctx)
    if err != nil {
        return fmt.Errorf("failed to render email: %w", err)
    }

    env := m.Envelope()
    env.To = []string{p.to}

    return p.client.driver.Send(ctx, &Message{
        Envelope: env,
        HTML:     html,
    })
}

// Queue: Pre-render ke HTML statis -> enqueue ke Asynq / Redis
func (p *PendingMail) Queue(ctx context.Context, m Mailable) error {
    html, err := m.Render(ctx)
    if err != nil {
        return fmt.Errorf("failed to render email: %w", err)
    }

    env := m.Envelope()
    return p.client.queue.EnqueueSendEmail(ctx, p.to, env.Subject, html)
}
```

### 2. Log Driver untuk Dev/Test (`pkg/mail/drivers/log_driver.go`)
```go
package drivers

import (
    "context"
    "fmt"
    "log/slog"
)

type LogDriver struct{}

func NewLogDriver() *LogDriver {
    return &LogDriver{}
}

func (d *LogDriver) Send(ctx context.Context, msg *Message) error {
    banner := "==================== [EMAIL DISPATCHED (LOG DRIVER)] ===================="
    fmt.Printf("\n%s\nTo: %s\nSubject: %s\n\n%s\n%s\n\n",
        banner, msg.Envelope.To, msg.Envelope.Subject, msg.HTML, banner)

    slog.InfoContext(ctx, "Email logged to console",
        "to", msg.Envelope.To,
        "subject", msg.Envelope.Subject,
    )
    return nil
}
```

### 3. Background Worker Server (`cmd/worker/main.go`)
```go
package main

import (
    "log"
    "github.com/hibiken/asynq"
    "golang-starter-kit/internal/config"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Parse Redis URI lengkap (misal redis://localhost:6379 atau redis://:secret@host:port/0)
    redisOpt, err := asynq.ParseRedisURI(cfg.RedisURL)
    if err != nil {
        log.Fatalf("Invalid Redis URL: %v", err)
    }

    srv := asynq.NewServer(
        redisOpt,
        asynq.Config{
            Concurrency: 10,
            Queues: map[string]int{
                "critical": 6,
                "default":  3,
                "low":      1,
            },
        },
    )

    mux := asynq.NewServeMux()
    // Registrasi task handlers
    mux.HandleFunc("email:send", handleSendEmailTask)

    log.Println("Background worker running...")
    if err := srv.Run(mux); err != nil {
        log.Fatalf("Worker error: %v", err)
    }
}
```

---

## ✅ Kriteria Verifikasi
- [ ] Pemanggilan `mail.To().Send()` mencetak email dengan box rapi di terminal saat mode `MAIL_DRIVER=log`.
- [ ] Pemanggilan `mail.To().Queue()` memasukkan payload task ke Redis Asynq dalam waktu < 5 ms.
- [ ] Worker di `cmd/worker/main.go` berhasil menyedot task dari Redis dan mengeksekusi pengiriman email.
- [ ] Task yang gagal otomatis di-retry hingga 3 kali dengan jeda eksponensial.
