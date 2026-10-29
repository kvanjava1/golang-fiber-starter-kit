package middleware_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang-starter-kit/internal/middleware"
	"golang-starter-kit/pkg/response"
)

func TestMiddleware_Trace(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.Trace())

	app.Get("/test", func(c fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	traceID := resp.Header.Get("X-Correlation-ID")
	assert.NotEmpty(t, traceID)
}

func TestMiddleware_CustomTracePreserved(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.Trace())

	app.Get("/test", func(c fiber.Ctx) error {
		return c.SendString("ok")
	})

	customID := "my-custom-trace-id-123"
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Correlation-ID", customID)

	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, customID, resp.Header.Get("X-Correlation-ID"))
}

func TestMiddleware_GlobalErrorHandler_API(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.GlobalErrorHandler,
	})
	app.Use(middleware.Trace())

	app.Get("/api/v1/fail", func(c fiber.Ctx) error {
		return response.ErrNotFound
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/fail", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), `"success":false`)
	assert.Contains(t, string(body), `"Resource not found"`)
}

func TestMiddleware_GlobalErrorHandler_500(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.GlobalErrorHandler,
	})
	app.Use(middleware.Trace())

	app.Get("/api/v1/crash", func(c fiber.Ctx) error {
		return errors.New("unexpected database crash")
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/crash", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), `"success":false`)
	assert.Contains(t, string(body), `"Internal server error"`)
}
