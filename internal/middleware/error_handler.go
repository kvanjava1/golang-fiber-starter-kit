package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"

	"golang-starter-kit/pkg/logger"
	"golang-starter-kit/pkg/response"
)

// GlobalErrorHandler menangkap seluruh error yang belum tertangani dari handler HTTP
func GlobalErrorHandler(c fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	traceID := logger.CorrelationIDFromContext(c.Context())
	meta := response.NewMeta(traceID)

	code := http.StatusInternalServerError
	msg := "Internal Server Error"

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
		msg = fiberErr.Message
	} else {
		code, msg = response.MapError(err)
	}

	// Catat ke server logger untuk 5xx
	if code >= 500 {
		slog.ErrorContext(c.Context(), "Internal server error occurred",
			"error", err.Error(),
			"path", c.Path(),
			"method", c.Method(),
		)
	}

	// Jika route diawali /api/, selalu kirim response format JSON Envelope
	if strings.HasPrefix(c.Path(), "/api/") || c.Accepts("application/json") != "" {
		return c.Status(code).JSON(response.Error(msg, nil, meta))
	}

	// Fallback untuk Web browser
	return c.Status(code).SendString(msg)
}
