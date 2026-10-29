package user

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

type WebController struct {
	service *Service
}

func NewWebController(service *Service) *WebController {
	return &WebController{service: service}
}

// Index menampilkan halaman daftar data User
func (ctrl *WebController) Index(c fiber.Ctx) error {
	ctx := c.Context()
	items, err := ctrl.service.List(ctx)
	if err != nil {
		return err
	}

	c.Set(fiber.HeaderContentType, "text/html; charset=utf-8")
	return c.SendString(fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
    <title>User List - Go Starter Kit</title>
    <link rel="stylesheet" href="/public/css/output.css"/>
</head>
<body class="min-h-screen bg-slate-50 font-sans text-slate-900 p-6 sm:p-10">
    <div class="max-w-4xl mx-auto">
        <div class="bg-white rounded-xl shadow-sm border border-slate-200 p-6 mb-6 flex items-center justify-between">
            <div>
                <h1 class="text-2xl font-bold text-slate-900">User List</h1>
                <p class="text-sm text-slate-500 mt-1">Total Data: %d item</p>
            </div>
            <a href="/user/create" class="px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg shadow-sm transition">
                + Tambah User
            </a>
        </div>
        <div class="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
            <p class="text-slate-600 text-sm">Modul <strong>User</strong> berhasil dibuat. Anda dapat mengganti tampilan ini menggunakan template <code>a-h/templ</code> di folder <code>resources/views/web/user/</code>.</p>
        </div>
    </div>
</body>
</html>`, len(items)))
}

// CreateForm menampilkan formulir tambah data User
func (ctrl *WebController) CreateForm(c fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, "text/html; charset=utf-8")
	return c.SendString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
    <title>Tambah User - Go Starter Kit</title>
    <link rel="stylesheet" href="/public/css/output.css"/>
</head>
<body class="min-h-screen bg-slate-50 font-sans text-slate-900 p-6 sm:p-10">
    <div class="max-w-xl mx-auto bg-white rounded-xl shadow-sm border border-slate-200 p-8">
        <h1 class="text-2xl font-bold text-slate-900 mb-6">Tambah User Baru</h1>
        <form action="/user" method="POST" class="space-y-4">
            <div>
                <label class="block text-sm font-medium text-slate-700 mb-1">Title / Name *</label>
                <input type="text" name="title" required class="block w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500" placeholder="Contoh: Item baru"/>
            </div>
            <div class="flex items-center gap-3 pt-2">
                <button type="submit" class="px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg shadow-sm transition">
                    Simpan User
                </button>
                <a href="/user" class="px-4 py-2 text-sm font-medium text-slate-700 bg-slate-100 hover:bg-slate-200 rounded-lg transition">
                    Batal
                </a>
            </div>
        </form>
    </div>
</body>
</html>`)
}

// Store memproses data kiriman formulir web
func (ctrl *WebController) Store(c fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(http.StatusBadRequest).SendString("Input formulir tidak valid")
	}

	ctx := c.Context()
	if _, err := ctrl.service.Create(ctx, req); err != nil {
		return err
	}

	return c.Redirect().To("/user")
}
