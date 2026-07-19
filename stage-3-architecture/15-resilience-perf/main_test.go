package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

func TestRunDemo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var out bytes.Buffer
	if err := RunDemo(ctx, &out); err != nil {
		t.Fatalf("RunDemo: %v\noutput:\n%s", err, out.String())
	}
	text := out.String()
	for _, want := range []string{"success status=200", "fallback status=200", "rate-limit status=429"} {
		if !strings.Contains(text, want) {
			t.Fatalf("output missing %q:\n%s", want, text)
		}
	}
}
