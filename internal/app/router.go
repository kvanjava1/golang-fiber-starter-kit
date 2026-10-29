package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/static"

	"golang-starter-kit/internal/middleware"
	"golang-starter-kit/pkg/response"
)

// DomainRouter adalah interface yang diimplementasikan oleh setiap domain untuk registrasi route mandiri
type DomainRouter interface {
	RegisterRoutes(web fiber.Router, api fiber.Router)
}

func RegisterRoutes(app *fiber.App, c *Container) {
	// 📁 Static Assets
	app.Use("/uploads", static.New("./storages/uploads"))
	app.Use("/public", static.New("./public"))

	// 📊 Prometheus Metrics Endpoint
	app.Get("/metrics", middleware.PrometheusHandler())

	// 🌐 1. Master Group Web (/)
	web := app.Group("/")
	web.Get("/", func(ctx fiber.Ctx) error {
		return ctx.SendString("Welcome to Go Starter Kit")
	})

	// 🔌 2. Master Group REST API (/api/v1)
	api := app.Group("/api/v1", cors.New())
	api.Get("/ping", func(ctx fiber.Ctx) error {
		return ctx.JSON(response.Success("pong", map[string]string{"status": "ok"}))
	})

	// 🚀 3. Delegasi Domain Modules
	for _, domain := range c.Domains {
		domain.RegisterRoutes(web, api)
	}
}
