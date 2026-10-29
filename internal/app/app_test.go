package app_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang-starter-kit/internal/app"
	"golang-starter-kit/internal/config"
	"golang-starter-kit/pkg/database"
)

func TestApp_RouteDispatch(t *testing.T) {
	cfg := &config.Config{
		AppName: "Test Server",
		AppEnv:  "testing",
		Port:    3000,
	}

	dbMgr := database.NewManager("sqlite1")
	application := app.BuildApp(cfg, dbMgr)

	// 1. Test Web Root /
	reqWeb := httptest.NewRequest(http.MethodGet, "/", nil)
	respWeb, err := application.Server.Test(reqWeb)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respWeb.StatusCode)
	bodyWeb, _ := io.ReadAll(respWeb.Body)
	assert.Contains(t, string(bodyWeb), "Welcome to Go Starter Kit")

	// 2. Test REST API /api/v1/ping
	reqAPI := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	respAPI, err := application.Server.Test(reqAPI)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respAPI.StatusCode)
	bodyAPI, _ := io.ReadAll(respAPI.Body)
	assert.Contains(t, string(bodyAPI), `"success":true`)
	assert.Contains(t, string(bodyAPI), `"pong"`)

	// 3. Test X-Correlation-ID injected
	assert.NotEmpty(t, respAPI.Header.Get("X-Correlation-ID"))

	// 4. Test Health Check /api/v1/health
	reqHealth := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	respHealth, err := application.Server.Test(reqHealth)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respHealth.StatusCode)

	// 5. Test Prometheus Metrics /metrics
	reqMetrics := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	respMetrics, err := application.Server.Test(reqMetrics)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respMetrics.StatusCode)
	bodyMetrics, _ := io.ReadAll(respMetrics.Body)
	assert.Contains(t, string(bodyMetrics), "http_requests_total")
}
