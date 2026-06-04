package logger

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"os"
)

type contextKey struct{}

// Generate returns a random 16-character hex trace ID.
func Generate() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// WithTraceID attaches a trace ID to the context.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, contextKey{}, traceID)
}

// TraceIDFrom extracts the trace ID from a context.
func TraceIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(contextKey{}).(string)
	return id
}

// FromCtx returns the default logger pre-populated with the request's trace_id.
func FromCtx(ctx context.Context) *slog.Logger {
	return slog.Default().With("trace_id", TraceIDFrom(ctx))
}

// Init configures the global slog logger. Reads LOG_LEVEL (debug|info) from env.
// Set json=true for JSON output (production), false for text (development).
func Init(json bool) {
	level := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "debug" {
		level = slog.LevelDebug
	}
	opts := &slog.HandlerOptions{Level: level}
	var h slog.Handler
	if json {
		h = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		h = slog.NewTextHandler(os.Stdout, opts)
	}
	slog.SetDefault(slog.New(h))
}
