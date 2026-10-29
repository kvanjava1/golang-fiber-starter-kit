package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"golang-starter-kit/internal/config"
	"golang-starter-kit/internal/middleware"
)

// NewServer menginisialisasi Fiber instance dengan konfigurasi fail-safe dan global middleware
func NewServer(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		ErrorHandler: middleware.GlobalErrorHandler,
	})

	// Pipeline Global Middleware
	app.Use(middleware.Trace())
	app.Use(middleware.Prometheus())
	app.Use(middleware.Logger())
	app.Use(recover.New(recover.Config{
		EnableStackTrace: cfg.IsDevelopment(),
	}))

	return app
}
