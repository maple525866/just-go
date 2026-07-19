package tracekit

import (
	"context"
	"testing"
)

func TestTraceParentRoundTripAndParentage(t *testing.T) {
	root := WithRemote(context.Background(), "00-0123456789abcdef0123456789abcdef-0123456789abcdef-01")
	exporter := &MemoryExporter{}
	child, end := Start(root, exporter, "gateway", "request")
	end(nil)
	parsed, ok := Parse(Format(child))
	if !ok || parsed.TraceID != "0123456789abcdef0123456789abcdef" {
		t.Fatalf("parsed trace = %#v ok=%v", parsed, ok)
	}
	spans := exporter.Spans()
	if len(spans) != 1 || spans[0].ParentSpanID != "0123456789abcdef" {
		t.Fatalf("spans = %#v", spans)
	}
}

func TestParseRejectsInvalidTraceParent(t *testing.T) {
	for _, value := range []string{
		"00-0123456789abcdef0123456789abcdef-0123456789abcdef-zz",
		"00-0123456789abcdef0123456789abcdef-0123456789abcdef-01-extra",
		"00-0123456789ABCDEF0123456789abcdef-0123456789abcdef-01",
	} {
		if _, ok := Parse(value); ok {
			t.Fatalf("Parse(%q) accepted invalid traceparent", value)
		}
	}
}
