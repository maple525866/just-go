package main

import (
	"testing"
)

func TestBuildHandler(t *testing.T) {
	if h := buildHandler(); h == nil {
		t.Fatalf("buildHandler returned nil")
	}
}
