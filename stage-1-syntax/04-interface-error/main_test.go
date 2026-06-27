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
		{name: "default title", want: "=== 接口、错误与泛型学习报告 ==="},
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
				"interface（隐式实现 + 小接口）",
				"error（包装 + Is/As）",
				"generic（类型参数 + 约束）",
				"errors.Is(ErrUserNotFound): true",
				"Sum with Number constraint: 60",
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
