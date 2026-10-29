package example

import (
	"github.com/gofiber/fiber/v3"
	"github.com/uptrace/bun"

	"golang-starter-kit/pkg/cache"
)

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

// NewModule adalah provider mandiri yang merakit seluruh komponen domain example.
func NewModule(db *bun.DB, c cache.Cache) *RouteHandler {
	repo := NewRepository(db)
	svc := NewService(repo, c, nil)
	webCtrl := NewWebController(svc)
	apiCtrl := NewAPIController(svc)
	return NewRouteHandler(webCtrl, apiCtrl)
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
	apiGroup.Post("/send-email", h.apiCtrl.SendEmail)
}
