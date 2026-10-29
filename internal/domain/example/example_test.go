package example_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang-starter-kit/internal/domain/example"
	"golang-starter-kit/internal/middleware"
	"golang-starter-kit/pkg/cache"
	"golang-starter-kit/pkg/database"
	"golang-starter-kit/pkg/response"
)

func setupTestApp(t *testing.T) (*fiber.App, *database.DBManager, func()) {
	db, err := database.ConnectSQLite(":memory:")
	require.NoError(t, err)

	ctx := context.Background()
	_, err = db.ExecContext(ctx, `
		CREATE TABLE cruds (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			status TEXT NOT NULL DEFAULT 'active',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(t, err)

	dbMgr := database.NewManager("sqlite1")
	dbMgr.RegisterSQL("sqlite1", db)

	memCache := cache.NewMemoryCache()

	repo := example.NewRepository(db)
	svc := example.NewService(repo, memCache, nil)
	webCtrl := example.NewWebController(svc)
	apiCtrl := example.NewAPIController(svc)
	route := example.NewRouteHandler(webCtrl, apiCtrl)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.GlobalErrorHandler,
	})
	app.Use(middleware.Trace())

	web := app.Group("/")
	api := app.Group("/api/v1")
	route.RegisterRoutes(web, api)

	teardown := func() {
		_ = db.Close()
		_ = memCache.Close()
	}

	return app, dbMgr, teardown
}

func TestExampleCRUD_RESTAPI_Flow(t *testing.T) {
	app, _, teardown := setupTestApp(t)
	defer teardown()

	// 1. Create CRUD item (POST /api/v1/example/cruds)
	createReq := example.CreateCrudRequest{
		Title:       "Test CRUD Title",
		Description: "Detailed description for test",
	}
	body, _ := json.Marshal(createReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/example/cruds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createEnv response.Envelope
	respBytes, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(respBytes, &createEnv)
	require.NoError(t, err)
	assert.True(t, createEnv.Success)

	dataMap, ok := createEnv.Data.(map[string]any)
	require.True(t, ok)
	itemID := dataMap["id"].(string)
	assert.NotEmpty(t, itemID)
	assert.Equal(t, "Test CRUD Title", dataMap["title"])

	// 2. List CRUD items (GET /api/v1/example/cruds)
	reqList := httptest.NewRequest(http.MethodGet, "/api/v1/example/cruds", nil)
	respList, err := app.Test(reqList)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respList.StatusCode)

	// 3. Get Detail (GET /api/v1/example/cruds/:id)
	reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/example/cruds/"+itemID, nil)
	respGet, err := app.Test(reqGet)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respGet.StatusCode)

	// 4. Update (PUT /api/v1/example/cruds/:id)
	updateReq := example.UpdateCrudRequest{
		Title:       "Updated Title",
		Description: "Updated description",
		Status:      "inactive",
	}
	updateBody, _ := json.Marshal(updateReq)
	reqUpdate := httptest.NewRequest(http.MethodPut, "/api/v1/example/cruds/"+itemID, bytes.NewReader(updateBody))
	reqUpdate.Header.Set("Content-Type", "application/json")
	respUpdate, err := app.Test(reqUpdate)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respUpdate.StatusCode)

	// 5. Delete (DELETE /api/v1/example/cruds/:id)
	reqDel := httptest.NewRequest(http.MethodDelete, "/api/v1/example/cruds/"+itemID, nil)
	respDel, err := app.Test(reqDel)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respDel.StatusCode)

	// 6. Verify 404 after delete
	reqGetAfterDel := httptest.NewRequest(http.MethodGet, "/api/v1/example/cruds/"+itemID, nil)
	respAfterDel, err := app.Test(reqGetAfterDel)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, respAfterDel.StatusCode)
}

func TestExampleCRUD_ValidationRejection(t *testing.T) {
	app, _, teardown := setupTestApp(t)
	defer teardown()

	// Short title (< 3 chars)
	invalidReq := example.CreateCrudRequest{
		Title: "ab",
	}
	body, _ := json.Marshal(invalidReq)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/example/cruds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var env response.Envelope
	respBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &env)
	assert.False(t, env.Success)
	assert.NotEmpty(t, env.Errors)
}

func TestExampleCRUD_WebSSR_Endpoints(t *testing.T) {
	app, _, teardown := setupTestApp(t)
	defer teardown()

	// 1. GET /example/cruds (Table index HTML)
	req := httptest.NewRequest(http.MethodGet, "/example/cruds", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "<!doctype html>")
	assert.Contains(t, string(body), "Example CRUD Items")

	// 2. GET /example/cruds/create (Form HTML)
	reqCreate := httptest.NewRequest(http.MethodGet, "/example/cruds/create", nil)
	respCreate, err := app.Test(reqCreate)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, respCreate.StatusCode)
	bodyCreate, _ := io.ReadAll(respCreate.Body)
	assert.Contains(t, string(bodyCreate), "Tambah Item Baru")
}
