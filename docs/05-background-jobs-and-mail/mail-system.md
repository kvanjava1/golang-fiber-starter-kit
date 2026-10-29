# 02. Fluent Mail System

Paket [`pkg/mail`](file:///home/melodelavic/Documents/golang/golang-fiber-starter-kit/pkg/mail/mail.go) menyediakan sistem pengiriman email terpadu dengan API bergaya *fluent interface* (mirip Laravel Mail) dan dukungan multi-driver.

---

## 📬 Fitur Utama

1. **Fluent API:** Sintaks yang bersih dan mudah dibaca:
   ```go
   mailClient.To("budi@example.com").Send(ctx, welcomeMail)
   ```
2. **Dua Mode Pengiriman:**
   - `.Send(ctx, mailable)` : Pengiriman sinkron langsung saat itu juga.
   - `.Queue(ctx, mailable)` : Menyelipkan pengiriman ke antrean Asynq Redis di latar belakang.
3. **Multi-Driver Swappable:**
   Dapat diganti kapan saja melalui variabel `MAIL_DRIVER` di `.env`:
   - `log`: Mencetak format email lengkap ke konsol terminal (sangat cocok untuk development lokal & unit test tanpa butuh koneksi internet).
   - `smtp`: Mengirim lewat server SMTP standar (Gmail, Mailtrap, SES, Sendgrid, dll.).
   - `resend`: Mengirim lewat HTTP API modern dari [Resend](https://resend.com).

---

## ✉️ Anatomi Class Mailable

Setiap jenis email direpresentasikan oleh struct yang mengimplementasikan interface `Mailable`:

```go
type Mailable interface {
    Subject() string
    BuildHTML(ctx context.Context) (string, error)
}
```

Contoh implementasi class mailable ([`pkg/mail/mailables/welcome_mail.go`](file:///home/melodelavic/Documents/golang/golang-fiber-starter-kit/pkg/mail/mailables/welcome_mail.go)):

```go
package mailables

import (
    "bytes"
    "context"
    "golang-starter-kit/resources/views/emails"
)

type WelcomeMail struct {
    UserName string
}

func NewWelcomeMail(userName string) *WelcomeMail {
    return &WelcomeMail{UserName: userName}
}

func (m *WelcomeMail) Subject() string {
    return "Selamat Datang di Go Starter Kit!"
}

func (m *WelcomeMail) BuildHTML(ctx context.Context) (string, error) {
    var buf bytes.Buffer
    // Render menggunakan template templ yang type-safe
    component := emailviews.WelcomeEmail(m.UserName)
    err := component.Render(ctx, &buf)
    return buf.String(), err
}
```

---

## 🪄 Generator Mailable Baru

Untuk membuat mailable baru secara instan:
```bash
go run cmd/make/main.go make:mail invoice_paid
```
File akan otomatis terbuat di `pkg/mail/mailables/invoice_paid_mail.go`.
