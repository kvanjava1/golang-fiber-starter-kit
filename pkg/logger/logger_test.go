package logger_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"golang-starter-kit/pkg/logger"
)

func TestLogger_TraceCorrelationID(t *testing.T) {
	buf := new(bytes.Buffer)
	log := logger.New("info", true, buf) // JSON handler to buffer

	ctx := logger.WithCorrelationID(context.Background(), "req-12345")
	log.InfoContext(ctx, "Test log message")

	output := buf.String()
	assert.Contains(t, output, `"msg":"Test log message"`)
	assert.Contains(t, output, `"correlation_id":"req-12345"`)
}

func TestLogger_ChainingPreservesTrace(t *testing.T) {
	buf := new(bytes.Buffer)
	log := logger.New("info", true, buf)

	// Chaining .With() harus tetap mempertahankan TraceHandler
	childLog := log.With("module", "user")

	ctx := logger.WithCorrelationID(context.Background(), "req-chained-999")
	childLog.InfoContext(ctx, "Chained message")

	output := buf.String()
	assert.Contains(t, output, `"msg":"Chained message"`)
	assert.Contains(t, output, `"module":"user"`)
	assert.Contains(t, output, `"correlation_id":"req-chained-999"`)
}

func TestLogger_ParseLevel(t *testing.T) {
	assert.Equal(t, slog.LevelDebug, logger.ParseLevel("debug"))
	assert.Equal(t, slog.LevelInfo, logger.ParseLevel("info"))
	assert.Equal(t, slog.LevelWarn, logger.ParseLevel("warn"))
	assert.Equal(t, slog.LevelError, logger.ParseLevel("error"))
	assert.Equal(t, slog.LevelInfo, logger.ParseLevel("unknown"))
}
