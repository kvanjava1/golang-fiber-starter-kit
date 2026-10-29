package integration_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang-starter-kit/tests/helpers"
)

func TestIntegration_HealthAndMetrics_EndToEnd(t *testing.T) {
	server, teardown := helpers.SetupTestApp()
	defer teardown()

	// 1. Health check probe
	reqHealth := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	respHealth, err := server.Test(reqHealth)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respHealth.StatusCode)

	bodyHealth, _ := io.ReadAll(respHealth.Body)
	assert.Contains(t, string(bodyHealth), `"status":"healthy"`)

	// 2. Prometheus metrics
	reqMetrics := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	respMetrics, err := server.Test(reqMetrics)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respMetrics.StatusCode)

	bodyMetrics, _ := io.ReadAll(respMetrics.Body)
	assert.Contains(t, string(bodyMetrics), "http_requests_total")
	assert.Contains(t, string(bodyMetrics), "http_request_duration_seconds")
}
