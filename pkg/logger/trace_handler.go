package logger

import (
	"context"
	"log/slog"
)

type correlationIDKey struct{}

// WithCorrelationID menyuntikkan correlation_id ke context secara type-safe
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationIDKey{}, id)
}

// CorrelationIDFromContext mengekstrak correlation_id dari context
func CorrelationIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if id, ok := ctx.Value(correlationIDKey{}).(string); ok {
		return id
	}
	return ""
}

type TraceHandler struct {
	slog.Handler
}

func NewTraceHandler(handler slog.Handler) *TraceHandler {
	return &TraceHandler{Handler: handler}
}

func (h *TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	if traceID := CorrelationIDFromContext(ctx); traceID != "" {
		r.AddAttrs(slog.String("correlation_id", traceID))
	}
	return h.Handler.Handle(ctx, r)
}

// WithAttrs memastikan wrapper TraceHandler tidak hilang saat chaining logger.With(...)
func (h *TraceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TraceHandler{Handler: h.Handler.WithAttrs(attrs)}
}

// WithGroup memastikan wrapper TraceHandler tidak hilang saat chaining logger.WithGroup(...)
func (h *TraceHandler) WithGroup(name string) slog.Handler {
	return &TraceHandler{Handler: h.Handler.WithGroup(name)}
}
