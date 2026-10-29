# 02. Server-Side Rendering (a-h/templ)

Starter kit ini menggunakan [`a-h/templ`](https://templ.guide/) sebagai mesin rendering tampilan HTML.

---

## 💡 Mengapa `a-h/templ`?

1. **Type-Safe saat Compile-Time:**
   Setiap variabel yang dikirimkan ke dalam template dijamin tipe datanya oleh Go compiler. Jika ada typo nama field atau tipe data salah, kompilasi akan langsung gagal sebelum aplikasi berjalan.
2. **Kinerja Ekstrem:**
   Template `.templ` dikompilasi menjadi kode Go native murni (`*_templ.go`) yang langsung menulis string HTML ke buffer tanpa parsing teks runtime yang lambat.
3. **Ergonomis:**
   Mendukung if-else, for loops, modular components, dan passing parameter persis seperti memanggil fungsi Go biasa.

---

## 📁 Struktur File Template

Semua file template disimpan di direktori [`resources/views/`](file:///home/melodelavic/Documents/golang/golang-fiber-starter-kit/resources/views/):
- `resources/views/web/<domain>/` : Halaman Web per modul (misal: `example/index.templ`, `example/create.templ`).
- `resources/views/emails/` : Template email notifikasi.

---

## ✍️ Contoh Penulisan Template `.templ`

Contoh dari file [resources/views/web/example/index.templ](file:///home/melodelavic/Documents/golang/golang-fiber-starter-kit/resources/views/web/example/index.templ):

```templ
package exampleviews

import (
    "fmt"
    "time"
)

type CrudItem struct {
    ID          string
    Title       string
    Description string
    Status      string
    CreatedAt   time.Time
}

templ Index(items []CrudItem) {
    <!DOCTYPE html>
    <html lang="en">
        <head>
            <meta charset="utf-8"/>
            <title>Daftar Item</title>
            <link rel="stylesheet" href="/public/css/output.css"/>
        </head>
        <body class="bg-slate-50 p-10 font-sans">
            <h1 class="text-2xl font-bold text-slate-900">Daftar Item</h1>
            if len(items) == 0 {
                <p class="text-slate-500">Belum ada item.</p>
            } else {
                <ul class="mt-4 space-y-2">
                    for _, item := range items {
                        <li class="p-4 bg-white rounded shadow-sm">
                            <span class="font-semibold">{ item.Title }</span>
                        </li>
                    }
                </ul>
            }
        </body>
    </html>
}
```

---

## 🚀 Cara Merender Template di Web Controller Fiber

Gunakan adapter resmi Fiber v3 untuk menyajikan komponen templ:

```go
package example

import (
    "github.com/gofiber/fiber/v3"
    "golang-starter-kit/resources/views/web/example"
)

func (c *WebController) Index(ctx fiber.Ctx) error {
    items, err := c.service.FindAll(ctx.Context())
    if err != nil {
        return err
    }

    // Set Content-Type ke text/html
    ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)

    // Render komponen templ langsung ke response stream Fiber
    component := exampleviews.Index(items)
    return component.Render(ctx.Context(), ctx.Response().BodyWriter())
}
```

---

## ⚡ Kompilasi Template

Setiap kali Anda mengubah file `.templ`, jalankan perintah:
```bash
make templ
# Atau: templ generate
```
Jika Anda menggunakan `make dev` atau `make dev-web`, perintah ini otomatis dijalankan sebelum server menyala.
