package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang-starter-kit/internal/config"
)

func TestConfig_DefaultValues(t *testing.T) {
	// Memastikan load saat tidak ada .env menghasilkan nilai default yang tepat
	cfg, err := config.Load("non-existent-file.env")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "development", cfg.AppEnv)
	assert.Equal(t, "Go Starter Kit", cfg.AppName)
	assert.Equal(t, 3000, cfg.Port)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "sqlite1", cfg.DBDefaultConnection)
	assert.True(t, cfg.DBSqliteEnabled)
	assert.Equal(t, "redis://localhost:6379", cfg.RedisURL)
	assert.Equal(t, "log", cfg.MailDriver)

	assert.True(t, cfg.IsDevelopment())
	assert.False(t, cfg.IsProduction())
	assert.False(t, cfg.IsTesting())
}

func TestConfig_LoadFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, "test.env")

	content := `
APP_ENV=production
APP_NAME="Custom App"
PORT=8080
LOG_LEVEL=debug
DB_DEFAULT_CONNECTION=pg1
DB_PG_ENABLED=true
DB_PG_1_DSN=postgres://user:pass@localhost:5432/proddb
MAIL_DRIVER=smtp
`
	err := os.WriteFile(envPath, []byte(content), 0644)
	require.NoError(t, err)

	cfg, err := config.Load(envPath)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "production", cfg.AppEnv)
	assert.Equal(t, "Custom App", cfg.AppName)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, "pg1", cfg.DBDefaultConnection)
	assert.True(t, cfg.DBPGEnabled)
	assert.Equal(t, "postgres://user:pass@localhost:5432/proddb", cfg.DBPG1DSN)
	assert.Equal(t, "smtp", cfg.MailDriver)

	assert.True(t, cfg.IsProduction())
	assert.False(t, cfg.IsDevelopment())
}
