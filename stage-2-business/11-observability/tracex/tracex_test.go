package tracex

import (
	"context"
	"testing"
)

func TestTraceContextAndChildSpan(t *testing.T) {
	ctx := WithTraceID(context.Background(), "trace-123")
	if got := TraceID(ctx); got != "trace-123" {
		t.Fatalf("TraceID() = %q, want trace-123", got)
	}

	ctx, root := StartSpan(ctx, "http.request")
	childCtx, child := StartSpan(ctx, "repo.query")

	if root.TraceID != "trace-123" || child.TraceID != "trace-123" {
		t.Fatalf("spans should reuse trace id: root=%+v child=%+v", root, child)
	}
	if child.ParentID != root.SpanID {
		t.Fatalf("child parent = %q, want %q", child.ParentID, root.SpanID)
	}
	if CurrentSpan(childCtx).Name != "repo.query" {
		t.Fatalf("current span not propagated")
	}
}
