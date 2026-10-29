# ⚠️ Analisis Masalah: Paradigma & Pola Idiomatik Golang

Dokumen ini memetakan perbedaan mendasar antara model eksekusi **TypeScript (Node/Bun)** dan **Golang**, serta potensi kendala arsitektural yang muncul saat mem-porting starter kit.

---

## 1. Perbedaan Paradigma: Event Loop vs Goroutines

| Aspek | Bun / TypeScript (Eksisting) | Golang (Target) |
|---|---|---|
| **Model Konkurensi** | Single-threaded Event Loop dengan Promise asinkron (`async/await`) | Multi-threaded M:N scheduler dengan **Goroutines** & Channels |
| **I/O Blocking** | Non-blocking via runtime async I/O | Blocking secara sintaksis, namun non-blocking secara OS thread berkat Go runtime scheduler |
| **Background Processing** | Wajib menggunakan antrean eksternal (Redis/BullMQ) agar tidak membebani loop utama | Bisa menggunakan *in-process worker pools* via goroutines untuk tugas ringan, atau Redis antrean untuk tugas kritis |

### Implikasi Desain:
- Di Go, kita tidak memerlukan `async/await`. Semua pemanggilan fungsi tampak sinkron (`res, err := repo.FindByID(id)`), namun runtime Go secara otomatis mengelola peralihan thread di latar belakang.
- Kita harus berhati-hati dalam menangani *race condition* dan data sharing antar-goroutine menggunakan mutex atau channel.

---

## 2. Error Handling: Exception Throwing vs Explicit Error Values

- **Eksisting (TypeScript):**
  Service melempar exception:
  ```typescript
  if (!user) throw new NotFoundError("User not found");
  ```
  Exception ini otomatis ditangkap oleh global error middleware dan diubah menjadi response HTTP 404.
- **Tantangan di Go:**
  Go menolak konsep exception untuk alur kontrol normal. Go mewajibkan nilai balik eksplisit `(data, error)`:
  ```go
  user, err := svc.repo.FindByID(ctx, id)
  if err != nil {
      return nil, err
  }
  ```
- **Masalah yang Harus Diselesaikan:**
  Bagaimana Service di Go dapat mengembalikan jenis error yang spesifik (misal: Not Found, Validation, Conflict) tanpa membocorkan status code HTTP ke dalam domain Service, namun Controller tetap dapat memetakannya secara otomatis ke status code HTTP yang benar?

---

## 3. Request Context: `AsyncLocalStorage` vs `context.Context`

- **Eksisting (TypeScript):**
  Menggunakan `AsyncLocalStorage` (`requestContext`). Fungsi apa pun di layer Service atau Util dapat memanggil `requestContext.getStore()?.requestId` tanpa menerima argumen context dari Controller.
- **Tantangan di Go:**
  Go tidak memiliki implicit thread-local context bawaan. Standar idiomatik Go mewajibkan parameter `ctx context.Context` dioper secara eksplisit sebagai **parameter pertama** pada setiap method:
  ```go
  func (s *CrudService) Create(ctx context.Context, input CreateInput) (*Crud, error)
  ```
- **Implikasi:**
  Seluruh layer (Controller -> Service -> Repository -> Mail -> Queue) harus secara konsisten menerima dan meneruskan `ctx context.Context` untuk memastikan timeout, pembatalan request (*cancellation*), dan pelacakan *Correlation ID* berfungsi dengan benar.

---

## 4. Class & Inheritance vs Struct Composition & Interfaces

- **Eksisting (TypeScript):**
  Menggunakan class berorientasi objek: `class WelcomeMail extends Mailable`, `class AppError extends Error`.
- **Tantangan di Go:**
  Go tidak memiliki `class` dan `extends`. Go menggunakan `struct`, *composition/embedding*, dan *implicit interfaces*.
- **Masalah yang Harus Diselesaikan:**
  Bagaimana merancang pola `Mailable` di Go agar tetap memiliki kemudahan penggunaan (*ergonomics*) yang setara dengan Laravel / Hono kit eksisting tanpa melanggar prinsip idiomatik Go.
