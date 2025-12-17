package trace

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

// Define context keys as constants
const (
	ContextTrace ContextKey = "trace_id"
	ContextSpan  ContextKey = "span_id"
)

// WithTraceID adds trace ID to context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, ContextTrace, traceID)
}

// GetTraceID extracts trace ID from context
func GetTraceID(ctx context.Context) (string, bool) {
	if traceID := ctx.Value(ContextTrace); traceID != nil {
		if id, ok := traceID.(string); ok {
			return id, true
		}
	}
	return "", false
}

// WithSpanID adds span ID to context
func WithSpanID(ctx context.Context, spanID string) context.Context {
	return context.WithValue(ctx, ContextSpan, spanID)
}

// GetSpanID extracts span ID from context
func GetSpanID(ctx context.Context) (string, bool) {
	if spanID := ctx.Value(ContextSpan); spanID != nil {
		if id, ok := spanID.(string); ok {
			return id, true
		}
	}
	return "", false
}

// WithTrace adds both trace ID and span ID to context
func WithTrace(ctx context.Context, traceID, spanID string) context.Context {
	return WithSpanID(WithTraceID(ctx, traceID), spanID)
}

// NewDefaultTraceContext creates a new background context with a fresh trace ID and span ID
func NewDefaultTraceContext() context.Context {
	return NewTraceContext(context.Background())
}

// NewTraceContext creates a new context with a fresh trace ID and span ID
func NewTraceContext(ctx context.Context) context.Context {
	return WithTrace(ctx, NewTraceID(), NewSpanID())
}

// NewChildSpan creates a new child span context with the same trace ID but new span ID
func NewChildSpan(ctx context.Context) context.Context {
	if traceID, ok := GetTraceID(ctx); ok {
		return WithTrace(ctx, traceID, NewSpanID())
	}
	// If no trace ID exists, create a new trace
	return NewTraceContext(ctx)
}

// NewTraceID generates a new trace ID using high-performance random generation
// Format: 16-character hex string (64-bit)
func NewTraceID() string {
	return IDGen(16)
}

// NewSpanID generates a new span ID using high-performance random generation
// Format: 8-character hex string (32-bit)
func NewSpanID() string {
	return IDGen(8)
}

// IDGen generates a new ID using high-performance random generation
// Format: [size*4]-character hex string [size*16]bit )
func IDGen(size int) string {
	bytes := make([]byte, size)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
