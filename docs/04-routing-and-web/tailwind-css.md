# 03. Tailwind CSS On-Demand

Starter kit ini mendukung **Tailwind CSS Standalone CLI** yang bersifat opsional (*opt-in / on-demand*). Anda mendapatkan seluruh kekuatan styling modern Tailwind tanpa harus menginstall Node.js, npm, webpack, ataupun Vite.

---

## 🎯 Filosofi Desain

- **API-First Developer:** Jika developer hanya menggunakan starter kit ini untuk membuat REST API, mereka **tidak dibebani** oleh dependency frontend apa pun. Perintah `make dev` tetap murni menyalakan server Go.
- **Web Developer:** Jika developer ingin membuat UI web yang cantik, mereka cukup menjalankan `make tailwind-install` sekali saja untuk mengunduh binary standalone resmi dari Tailwind Labs.

---

## 🛠️ Alur Kerja & Perintah

### 1. Inisialisasi Pertama Kali
Unduh binary mandiri Tailwind CLI (hanya dilakukan sekali):
```bash
make tailwind-install
```
Perintah ini akan:
1. Mendeteksi OS dan arsitektur CPU host secara otomatis (`linux-x64`, `macos-arm64`, `windows-x64`, dll.).
2. Mengunduh binary ke `bin/tailwindcss`.
3. Menghasilkan file konfigurasi [`tailwind.config.js`](file:///home/melodelavic/Documents/golang/golang-fiber-starter-kit/tailwind.config.js).
4. Menghasilkan file input CSS [`resources/css/input.css`](file:///home/melodelavic/Documents/golang/golang-fiber-starter-kit/resources/css/input.css).

### 2. Kompilasi Sekali Jalan (Production / Testing)
Untuk membangun dan meminifikasi CSS ke `public/css/output.css`:
```bash
make css
```

### 3. Live Watcher saat Pengembangan Aktif
Saat Anda sedang aktif mengedit utility class Tailwind di file `.templ`, jalankan watcher di tab terminal terpisah:
```bash
make css-watch
```
Watcher akan memantau perubahan file `.templ` dan langsung merefresh `public/css/output.css` dalam hitungan milidetik.

### 4. Mode Web Terpadu (`make dev-web`)
Perintah ini otomatis memastikan binary terpasang, mengompilasi template `templ`, membangun CSS, lalu menyalakan server Go di port `3000`:
```bash
make dev-web
```

---

## ⚙️ Konfigurasi Pemindaian Template

File `tailwind.config.js` telah dikonfigurasi untuk memindai seluruh file template Go:

```javascript
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./resources/views/**/*.templ",
    "./resources/views/**/*.html",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
```
Hasil build disimpan di `public/css/output.css` dan di-link pada head HTML:
```html
<link rel="stylesheet" href="/public/css/output.css"/>
```
