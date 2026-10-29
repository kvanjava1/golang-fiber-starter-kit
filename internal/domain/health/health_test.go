package health_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang-starter-kit/internal/domain/health"
	"golang-starter-kit/pkg/database"
	"golang-starter-kit/pkg/response"
)

func TestHealthCheck_Healthy(t *testing.T) {
	app := fiber.New()
	dbMgr := database.NewManager("sqlite1")

	db, err := database.ConnectSQLite(":memory:")
	require.NoError(t, err)
	defer db.Close()

	dbMgr.RegisterSQL("sqlite1", db)

	ctrl := health.NewController(dbMgr, nil)
	route := health.NewRouteHandler(ctrl)

	apiGroup := app.Group("/api/v1")
	webGroup := app.Group("/")
	route.RegisterRoutes(webGroup, apiGroup)

	// 1. Test /api/v1/health
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	var env response.Envelope
	err = json.Unmarshal(bodyBytes, &env)
	require.NoError(t, err)
	assert.True(t, env.Success)
	assert.Equal(t, "healthy", env.Message)

	// 2. Test /health
	reqWeb := httptest.NewRequest(http.MethodGet, "/health", nil)
	respWeb, err := app.Test(reqWeb)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respWeb.StatusCode)
}
