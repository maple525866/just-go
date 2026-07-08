package tracex

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

type traceKey struct{}
type spanKey struct{}

// Span is a teaching-friendly trace span model.
type Span struct {
	TraceID  string
	SpanID   string
	ParentID string
	Name     string
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceKey{}, traceID)
}

func TraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(traceKey{}).(string); ok {
		return traceID
	}
	return ""
}

func StartSpan(ctx context.Context, name string) (context.Context, Span) {
	traceID := TraceID(ctx)
	if traceID == "" {
		traceID = randomHex(16)
		ctx = WithTraceID(ctx, traceID)
	}
	parent := CurrentSpan(ctx)
	span := Span{TraceID: traceID, SpanID: randomHex(8), ParentID: parent.SpanID, Name: name}
	return context.WithValue(ctx, spanKey{}, span), span
}

func CurrentSpan(ctx context.Context) Span {
	if span, ok := ctx.Value(spanKey{}).(Span); ok {
		return span
	}
	return Span{}
}

func randomHex(bytesLen int) string {
	b := make([]byte, bytesLen)
	if _, err := rand.Read(b); err != nil {
		return "0000000000000000"
	}
	return hex.EncodeToString(b)
}
