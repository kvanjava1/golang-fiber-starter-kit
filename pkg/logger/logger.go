package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// ParseLevel mengubah string level log ke slog.Level
func ParseLevel(lvl string) slog.Level {
	switch strings.ToLower(lvl) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// New menginisialisasi slog.Logger dengan TraceHandler terpasang
func New(levelStr string, isProd bool, out ...io.Writer) *slog.Logger {
	var writer io.Writer = os.Stdout
	if len(out) > 0 && out[0] != nil {
		writer = out[0]
	}

	lvl := ParseLevel(levelStr)
	opts := &slog.HandlerOptions{
		Level: lvl,
	}

	var baseHandler slog.Handler
	if isProd {
		baseHandler = slog.NewJSONHandler(writer, opts)
	} else {
		baseHandler = slog.NewTextHandler(writer, opts)
	}

	traceHandler := NewTraceHandler(baseHandler)
	return slog.New(traceHandler)
}

// SetDefault menginisialisasi dan menyetel logger global default bawaan Go
func SetDefault(levelStr string, isProd bool) *slog.Logger {
	l := New(levelStr, isProd)
	slog.SetDefault(l)
	return l
}
