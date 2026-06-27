package profile

import (
	"strings"
	"testing"
)

func TestTypes(t *testing.T) {
	tests := []struct {
		name     string
		contains []string
	}{
		{name: "common profiles", contains: []string{"CPU", "memory", "blocking"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profiles := Types()
			joined := ""
			for _, p := range profiles {
				joined += p.Name + " " + p.Purpose + "\n"
			}
			for _, part := range tt.contains {
				if !strings.Contains(joined, part) {
					t.Fatalf("Types() = %#v, want text containing %q", profiles, part)
				}
			}
		})
	}
}

func TestCommands(t *testing.T) {
	tests := []struct {
		name     string
		contains []string
	}{
		{name: "pprof commands", contains: []string{"go test -bench", "-cpuprofile", "go tool pprof"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			joined := strings.Join(Commands(), "\n")
			for _, part := range tt.contains {
				if !strings.Contains(joined, part) {
					t.Fatalf("Commands() = %q, want it to contain %q", joined, part)
				}
			}
		})
	}
}
