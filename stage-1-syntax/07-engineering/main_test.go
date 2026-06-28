package main

import (
	"strings"
	"testing"
)

func TestReportTitle(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "default title", want: "=== 工程化基础学习报告 ==="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reportTitle(); got != tt.want {
				t.Fatalf("reportTitle() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildReport(t *testing.T) {
	tests := []struct {
		name     string
		contains []string
	}{
		{
			name: "includes all topic sections",
			contains: []string{
				"module / go.work / semantic version",
				"testing / benchmark",
				"quality gates",
				"debug / slog / dlv",
				"pprof",
				"go test -bench=.",
				"golangci-lint run",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildReport()
			for _, part := range tt.contains {
				if !strings.Contains(got, part) {
					t.Fatalf("buildReport() = %q, want it to contain %q", got, part)
				}
			}
		})
	}
}
