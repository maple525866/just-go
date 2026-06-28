package system

import (
	"strings"
	"testing"
)

func TestWriteReadTempFile(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{name: "round trip", content: "hello file"},
		{name: "empty", content: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WriteReadTempFile(tt.content)
			if err != nil {
				t.Fatalf("WriteReadTempFile() unexpected error: %v", err)
			}
			if got != tt.content {
				t.Fatalf("WriteReadTempFile() = %q, want %q", got, tt.content)
			}
		})
	}
}

func TestEnvOrDefault(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		fallback string
		set      bool
		want     string
	}{
		{name: "uses env", key: "JUST_GO_STDLIB_TEST", value: "present", fallback: "fallback", set: true, want: "present"},
		{name: "uses fallback", key: "JUST_GO_STDLIB_MISSING", fallback: "fallback", want: "fallback"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.set {
				t.Setenv(tt.key, tt.value)
			}
			if got := EnvOrDefault(tt.key, tt.fallback); got != tt.want {
				t.Fatalf("EnvOrDefault() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGoVersion(t *testing.T) {
	tests := []struct {
		name       string
		wantPrefix string
	}{
		{name: "go env version", wantPrefix: "go"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GoVersion()
			if err != nil {
				t.Fatalf("GoVersion() unexpected error: %v", err)
			}
			if !strings.HasPrefix(got, tt.wantPrefix) {
				t.Fatalf("GoVersion() = %q, want prefix %q", got, tt.wantPrefix)
			}
		})
	}
}
