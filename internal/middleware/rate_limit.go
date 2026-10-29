package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"

	"golang-starter-kit/pkg/logger"
	"golang-starter-kit/pkg/response"
)

// RateLimiter membatasi jumlah request per IP
func RateLimiter(max int, exp time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: exp,
		LimitReached: func(c fiber.Ctx) error {
			traceID := logger.CorrelationIDFromContext(c.Context())
			meta := response.NewMeta(traceID)
			return c.Status(fiber.StatusTooManyRequests).JSON(
				response.Error("Too many requests, please try again later", nil, meta),
			)
		},
	})
}
