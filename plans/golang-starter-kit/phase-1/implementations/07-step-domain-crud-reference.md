# 🛠️ Step 07: Living Reference Domain (`internal/domain/example/`)

Dokumen ini memuat panduan implementasi modul referensi hidup (*living reference domain*) `example` menggunakan **Gaya A Co-location**, yang mencakup REST API JSON, Web SSR HTML (`a-h/templ`), Service domain murni, Repository `uptrace/bun`, dan *Self-Registering Routes*.

---

## 🎯 Tujuan Langkah
1. Membangun modul domain referensi lengkap di `internal/domain/example/`.
2. Menyediakan endpoint REST API: Create, Read (List & Detail), Update, Delete, dan Send Email (`POST /api/v1/example/cruds/send-email`).
3. Menyediakan antarmuka Web SSR HTML browser di `/example/cruds`.
4. Menghubungkan registrasi route mandiri via `route.go` ke master router.

---

## 📂 Struktur File yang Dibuat
```text
internal/
└── domain/
    └── example/
        ├── model.go          # Bun struct entity
        ├── dto.go            # Request / Response DTOs & validasi
        ├── repository.go     # Interface & Bun queries
        ├── service.go        # Pure business logic
        ├── controller_api.go # REST API JSON handler
        ├── controller_web.go # Web SSR HTML handler
        └── route.go          # RegisterRoutes(web, api)
resources/
└── views/
    └── web/
        └── example/
            ├── index.templ   # Halaman daftar tabel
            └── create.templ  # Form input data
```

---

## 📝 Detail Spesifikasi Kode

### 1. Model Entitas Database (`internal/domain/example/model.go`)
```go
package example

import (
    "time"
    "github.com/uptrace/bun"
)

type Crud struct {
    bun.BaseModel `bun:"table:cruds,alias:c"`
    ID          string    `bun:",pk"`
    Title       string    `bun:",notnull"`
    Description string    `bun:""`
    Status      string    `bun:",notnull,default:'active'"`
    CreatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp"`
    UpdatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
```

### 2. DTO & Validation (`internal/domain/example/dto.go`)
```go
package example

type CreateCrudRequest struct {
    Title       string `json:"title" validate:"required,min=3,max=100"`
    Description string `json:"description" validate:"max=500"`
}

type SendEmailRequest struct {
    Email string `json:"email" validate:"required,email"`
    Name  string `json:"name" validate:"required,min=2"`
}
```

### 3. Self-Registering Routes (`internal/domain/example/route.go`)
```go
package example

import "github.com/gofiber/fiber/v3"

type RouteHandler struct {
    webCtrl *WebController
    apiCtrl *APIController
}

func NewRouteHandler(webCtrl *WebController, apiCtrl *APIController) *RouteHandler {
    return &RouteHandler{
        webCtrl: webCtrl,
        apiCtrl: apiCtrl,
    }
}

func (h *RouteHandler) RegisterRoutes(web fiber.Router, api fiber.Router) {
    // 🖥️ 1. Web SSR Routes (Browser) -> Base URL: /example/cruds
    webGroup := web.Group("/example/cruds")
    webGroup.Get("/", h.webCtrl.Index)
    webGroup.Get("/create", h.webCtrl.Create)
    webGroup.Post("/", h.webCtrl.Store)

    // 🔌 2. REST API Routes (JSON) -> Base URL: /api/v1/example/cruds
    apiGroup := api.Group("/example/cruds")
    apiGroup.Get("/", h.apiCtrl.List)
    apiGroup.Post("/", h.apiCtrl.Create)
    apiGroup.Get("/:id", h.apiCtrl.GetByID)
    apiGroup.Put("/:id", h.apiCtrl.Update)
    apiGroup.Delete("/:id", h.apiCtrl.Delete)
    apiGroup.Post("/send-email", h.apiCtrl.SendEmail) // Endpoint contoh dispatch email queue
}
```

### 4. API Send Email Handler (`internal/domain/example/controller_api.go`)
```go
func (ctrl *Controller) APISendEmail(c fiber.Ctx) error {
    var req SendEmailRequest
    if err := c.Bind().JSON(&req); err != nil {
        return c.Status(400).JSON(response.Error(400, "Invalid JSON body"))
    }

    if err := validator.Validate(req); err != nil {
        return c.Status(400).JSON(response.ValidationError(err))
    }

    ctx := c.Context()
    if err := ctrl.service.SendExampleEmail(ctx, req.Email, req.Name); err != nil {
        return response.HandleError(c, err)
    }

    return c.Status(202).JSON(response.Envelope{
        Success: true,
        Message: "Example email queued successfully",
    })
}
```

---

## ✅ Kriteria Verifikasi
- [ ] Endpoint `POST /api/v1/example/cruds` berhasil membuat record baru dan mengembalikan HTTP 201.
- [ ] Endpoint `GET /example/cruds` di browser merender halaman HTML berbasis `a-h/templ`.
- [ ] Endpoint `POST /api/v1/example/cruds/send-email` mengembalikan HTTP 202 dan men-dispatch email ke antrean Asynq.
- [ ] Input yang tidak valid (misal email salah format) otomatis ditolak dengan pesan error yang jelas.
