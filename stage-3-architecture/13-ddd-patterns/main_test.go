package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunCompleteConfirmationFlow(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	if err := run(&output); err != nil {
		t.Fatalf("run() error = %v", err)
	}
	for _, want := range []string{
		"order=order-2026",
		"status=confirmed",
		"total=10000 CNY",
		"inventory_reserved=2",
	} {
		if !strings.Contains(output.String(), want) {
			t.Fatalf("output %q does not contain %q", output.String(), want)
		}
	}
}
