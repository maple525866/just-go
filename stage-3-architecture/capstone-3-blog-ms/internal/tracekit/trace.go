package tracekit

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const metadataKey = "traceparent"

type Context struct {
	TraceID string
	SpanID  string
}

type Span struct {
	Name         string
	Service      string
	TraceID      string
	SpanID       string
	ParentSpanID string
	StartedAt    time.Time
	EndedAt      time.Time
	Error        string
}

type Exporter interface {
	Export(Span)
}

type MemoryExporter struct {
	mu    sync.Mutex
	spans []Span
}

func (e *MemoryExporter) Export(span Span) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.spans = append(e.spans, span)
}

func (e *MemoryExporter) Spans() []Span {
	e.mu.Lock()
	defer e.mu.Unlock()
	return append([]Span(nil), e.spans...)
}

type contextKey struct{}

func FromContext(ctx context.Context) (Context, bool) {
	value, ok := ctx.Value(contextKey{}).(Context)
	return value, ok && validHex(value.TraceID, 32) && validHex(value.SpanID, 16)
}

func WithRemote(ctx context.Context, traceparent string) context.Context {
	traceContext, ok := Parse(traceparent)
	if !ok {
		traceContext = Context{TraceID: randomHex(16), SpanID: randomHex(8)}
	}
	return context.WithValue(ctx, contextKey{}, traceContext)
}

func Start(ctx context.Context, exporter Exporter, service, name string) (context.Context, func(error)) {
	parent, ok := FromContext(ctx)
	if !ok {
		parent = Context{TraceID: randomHex(16)}
	}
	current := Context{TraceID: parent.TraceID, SpanID: randomHex(8)}
	started := time.Now().UTC()
	ctx = context.WithValue(ctx, contextKey{}, current)
	return ctx, func(err error) {
		if exporter == nil {
			return
		}
		span := Span{Name: name, Service: service, TraceID: current.TraceID, SpanID: current.SpanID, ParentSpanID: parent.SpanID, StartedAt: started, EndedAt: time.Now().UTC()}
		if err != nil {
			span.Error = err.Error()
		}
		exporter.Export(span)
	}
}

func Format(ctx context.Context) string {
	value, ok := FromContext(ctx)
	if !ok {
		value = Context{TraceID: randomHex(16), SpanID: randomHex(8)}
	}
	return fmt.Sprintf("00-%s-%s-01", value.TraceID, value.SpanID)
}

func Parse(value string) (Context, bool) {
	parts := strings.Split(value, "-")
	if len(parts) != 4 {
		return Context{}, false
	}
	version, traceID, spanID, flags := parts[0], parts[1], parts[2], parts[3]
	if version != "00" || !validHex(traceID, 32) || !validHex(spanID, 16) || !validHex(flags, 2) || allZero(traceID) || allZero(spanID) {
		return Context{}, false
	}
	return Context{TraceID: traceID, SpanID: spanID}, true
}

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, metadataKey, Format(ctx))
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func UnaryServerInterceptor(service string, exporter Exporter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (reply any, retErr error) {
		if values := metadata.ValueFromIncomingContext(ctx, metadataKey); len(values) > 0 {
			ctx = WithRemote(ctx, values[0])
		}
		ctx, end := Start(ctx, exporter, service, info.FullMethod)
		defer func() { end(retErr) }()
		return handler(ctx, req)
	}
}

func validHex(value string, size int) bool {
	if len(value) != size {
		return false
	}
	if value != strings.ToLower(value) {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}

func allZero(value string) bool {
	for _, ch := range value {
		if ch != '0' {
			return false
		}
	}
	return true
}

func randomHex(bytes int) string {
	buffer := make([]byte, bytes)
	if _, err := rand.Read(buffer); err != nil {
		panic(fmt.Sprintf("generate trace id: %v", err))
	}
	return hex.EncodeToString(buffer)
}
