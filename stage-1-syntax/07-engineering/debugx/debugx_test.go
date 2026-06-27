package debugx

import (
	"strings"
	"testing"
)

func TestLogExample(t *testing.T) {
	tests := []struct {
		name      string
		component string
		step      int
		contains  []string
	}{
		{name: "structured slog", component: "calc", step: 2, contains: []string{"level=INFO", "msg=\"debug checkpoint\"", "component=calc", "step=2"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LogExample(tt.component, tt.step)
			for _, part := range tt.contains {
				if !strings.Contains(got, part) {
					t.Fatalf("LogExample() = %q, want it to contain %q", got, part)
				}
			}
		})
	}
}

func TestDebugCommands(t *testing.T) {
	tests := []struct {
		name     string
		contains []string
	}{
		{name: "mentions dlv or breakpoint", contains: []string{"dlv test", "IDE breakpoint"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			joined := strings.Join(DebugCommands(), "\n")
			for _, part := range tt.contains {
				if !strings.Contains(joined, part) {
					t.Fatalf("DebugCommands() = %q, want it to contain %q", joined, part)
				}
			}
		})
	}
}
