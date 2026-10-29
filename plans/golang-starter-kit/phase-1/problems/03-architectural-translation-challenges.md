# ⚠️ Analisis Masalah: Tantangan Translasi Arsitektur

Dokumen ini membedah tantangan dalam mempertahankan struktur arsitektur modular, *Dependency Injection*, dan *Developer Experience (DX)* saat berpindah ke Golang.

---

## 1. Tantangan Struktur Direktori: Go Standard Project Layout

Di proyek TypeScript/Bun, kita bebas menaruh seluruh kode aplikasi di dalam `src/`. Namun dalam ekosistem Go, terdapat konvensi komunitas yang kuat (**Standard Go Project Layout**):

```text
go-starter-kit/
├── cmd/               # Entrypoint aplikasi yang dapat di-build menjadi binary terpisah
│   ├── api/           # Entrypoint HTTP Web & REST API
│   ├── worker/        # Entrypoint Background Job Worker
│   └── make/          # Entrypoint CLI Scaffolder tool
├── internal/          # Kode private domain aplikasi (tidak dapat diimpor oleh proyek eksternal)
│   ├── app/           # Inisialisasi & wiring dependency
│   ├── domain/        # Domain modules (controllers, services, repositories)
│   └── middleware/    # HTTP middleware
├── pkg/               # Paket utilitas & library teknis yang reusable (mail, cache, logger)
└── resources/         # Template HTML, asset statis, file migrasi
```

### Masalah yang Perlu Ditentukan:
- Bagaimana memetakan layer `routes/`, `controllers/`, `services/`, dan `repositories/` dari starter kit TS ke dalam `internal/` di Go agar tetap rapi, tidak terjadi *circular dependency* (siklus impor terlarang di Go), dan mudah dieksplorasi oleh developer.

---

## 2. Tantangan Dependency Injection (DI)

- **Eksisting di TypeScript:**
  Menggunakan singleton instance atau export langsung:
  ```typescript
  export const crudService = new CrudService(crudRepository);
  ```
- **Tantangan di Go:**
  - Di Go, membuat global package variable (`var DB *sql.DB`, `var Service = ...`) dianggap sebagai *anti-pattern* karena mempersulit unit testing secara paralel dan membuat state tersembunyi.
  - Cara idiomatik Go adalah menggunakan **Constructor Injection**:
    ```go
    type CrudService struct {
        repo   domain.CrudRepository
        logger *slog.Logger
        mail   pkg.MailClient
    }

    func NewCrudService(repo domain.CrudRepository, logger *slog.Logger, mail pkg.MailClient) *CrudService {
        return &CrudService{repo: repo, logger: logger, mail: mail}
    }
    ```
- **Masalah:**
  Jika proyek bertambah besar, wiring constructor manual di `main.go` bisa menjadi sangat panjang. Apakah kita perlu DI framework seperti **Google Wire** (compile-time code generation) atau **Uber Fx** (runtime reflection), atau cukup manual wiring yang bersih?

---

## 3. Tantangan Desain Mailable Pattern di Go

Di TypeScript, kita dapat membuat subclass dari `Mailable`:
```typescript
export class WelcomeMail extends Mailable {
  envelope() { return { subject: "Welcome" }; }
  content() { return <WelcomeEmail name={this.name} />; }
}
```

### Tantangan di Go:
- Go tidak memiliki pewarisan class.
- Jika menggunakan interface:
  ```go
  type Mailable interface {
      Envelope() Envelope
      RenderHTML() (string, error)
  }
  ```
  Bagaimana struct konkret mengimplementasikannya secara ringkas tanpa menulis boilerplate berulang untuk setiap email baru?
- Bagaimana merender template email: apakah menggunakan **`a-h/templ`** (yang sangat mirip komponen JSX dan ber-tipe aman saat kompilasi) atau `html/template` standar Go?
