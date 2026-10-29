package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
)

// Logger mencatat ringkasan HTTP request ke log terstruktur
func Logger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		durationMs := float64(duration.Microseconds()) / 1000.0

		slog.InfoContext(c.Context(), "HTTP Request",
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"duration_ms", durationMs,
			"ip", c.IP(),
		)

		return err
	}
}
