package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang-starter-kit/pkg/response"
	"golang-starter-kit/tests/helpers"
)

func TestIntegration_CRUD_EndToEnd(t *testing.T) {
	server, teardown := helpers.SetupTestApp()
	defer teardown()

	// 1. Create CRUD via API (POST /api/v1/example/cruds)
	payload := map[string]string{
		"title":       "Integration Test Item",
		"description": "Created during end-to-end integration test",
	}
	body, _ := json.Marshal(payload)

	reqPost := httptest.NewRequest(http.MethodPost, "/api/v1/example/cruds", bytes.NewReader(body))
	reqPost.Header.Set("Content-Type", "application/json")
	respPost, err := server.Test(reqPost)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, respPost.StatusCode)

	var envPost response.Envelope
	postBytes, _ := io.ReadAll(respPost.Body)
	err = json.Unmarshal(postBytes, &envPost)
	require.NoError(t, err)
	assert.True(t, envPost.Success)

	dataMap := envPost.Data.(map[string]any)
	itemID := dataMap["id"].(string)
	assert.NotEmpty(t, itemID)

	// 2. Fetch via API List (GET /api/v1/example/cruds)
	reqList := httptest.NewRequest(http.MethodGet, "/api/v1/example/cruds", nil)
	respList, err := server.Test(reqList)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respList.StatusCode)

	listBytes, _ := io.ReadAll(respList.Body)
	assert.Contains(t, string(listBytes), "Integration Test Item")

	// 3. Verify on Web SSR Browser Page (GET /example/cruds)
	reqWeb := httptest.NewRequest(http.MethodGet, "/example/cruds", nil)
	respWeb, err := server.Test(reqWeb)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respWeb.StatusCode)

	webHTML, _ := io.ReadAll(respWeb.Body)
	assert.Contains(t, string(webHTML), "<!doctype html>")
	assert.Contains(t, string(webHTML), "Integration Test Item")

	// 4. Delete via API (DELETE /api/v1/example/cruds/:id)
	reqDel := httptest.NewRequest(http.MethodDelete, "/api/v1/example/cruds/"+itemID, nil)
	respDel, err := server.Test(reqDel)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respDel.StatusCode)

	// 5. Verify 404 after deletion
	reqCheck := httptest.NewRequest(http.MethodGet, "/api/v1/example/cruds/"+itemID, nil)
	respCheck, err := server.Test(reqCheck)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, respCheck.StatusCode)
}
