package health

import (
	"github.com/gofiber/fiber/v3"
)

type RouteHandler struct {
	ctrl *Controller
}

func NewRouteHandler(ctrl *Controller) *RouteHandler {
	return &RouteHandler{ctrl: ctrl}
}

func (h *RouteHandler) RegisterRoutes(web fiber.Router, api fiber.Router) {
	// 🔌 REST API Group (/api/v1/health)
	api.Get("/health", h.ctrl.Check)

	// 🌐 Web Health Probe (/health)
	web.Get("/health", h.ctrl.Check)
}
