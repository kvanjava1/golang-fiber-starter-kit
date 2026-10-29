package helpers

import (
	"path/filepath"
	"runtime"

	"github.com/gofiber/fiber/v3"
	"github.com/pressly/goose/v3"

	"golang-starter-kit/internal/app"
	"golang-starter-kit/internal/config"
	"golang-starter-kit/pkg/cache"
	"golang-starter-kit/pkg/database"
)

func findProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	// tests/helpers/app.go -> root is 2 levels up
	return filepath.Join(filepath.Dir(b), "..", "..")
}

func SetupTestApp() (*fiber.App, func()) {
	root := findProjectRoot()

	// 1. Buat koneksi SQLite in-memory terisolasi
	db, err := database.ConnectSQLite(":memory:")
	if err != nil {
		panic(err)
	}

	// 2. Jalankan migrasi Goose ke database in-memory
	if err := goose.SetDialect("sqlite3"); err == nil {
		migrationsDir := filepath.Join(root, "resources", "migrations", "sqlite1")
		_ = goose.Up(db.DB, migrationsDir)
	}

	// 3. Register ke DBManager
	dbMgr := database.NewManager("sqlite1")
	dbMgr.RegisterSQL("sqlite1", db)

	// 4. Memory Cache terisolasi
	memCache := cache.NewMemoryCache()

	// 5. Inisialisasi Mock Config
	cfg := &config.Config{
		AppName:             "Integration Test Server",
		AppEnv:              "testing",
		Port:                3000,
		DBDefaultConnection: "sqlite1",
	}

	// 6. Build App Server
	application := app.BuildApp(cfg, dbMgr, memCache)

	teardown := func() {
		_ = db.Close()
		_ = memCache.Close()
	}

	return application.Server, teardown
}
