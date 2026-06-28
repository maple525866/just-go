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
		{name: "default title", want: "=== 并发编程学习报告 ==="},
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
			name: "includes all sections",
			contains: []string{
				"goroutine（启动 + 等待 + 收集）",
				"channel（通信 + close + select）",
				"sync（Mutex + RWMutex + Once）",
				"context（取消 + 超时）",
				"并发坑（安全说明）",
				"data race, goroutine leak, channel deadlock",
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
