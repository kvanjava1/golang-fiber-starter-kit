package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"golang-starter-kit/pkg/logger"
)

const HeaderCorrelationID = "X-Correlation-ID"

// Trace memastikan setiap request HTTP memiliki correlation_id unik dan diinjeksikan ke context
func Trace() fiber.Handler {
	return func(c fiber.Ctx) error {
		traceID := c.Get(HeaderCorrelationID)
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Set ke response header agar client dapat mereferensikan request
		c.Set(HeaderCorrelationID, traceID)

		// Simpan ke request context
		ctx := logger.WithCorrelationID(c.Context(), traceID)
		c.SetContext(ctx)

		return c.Next()
	}
}
