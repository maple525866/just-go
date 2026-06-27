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
		{name: "default title", want: "=== 标准库精要学习报告 ==="},
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
				"fmt（格式化）",
				"io / bufio（流式读写）",
				"os / os/exec（文件与进程）",
				"net/http（handler + client）",
				"encoding/json / xml（序列化）",
				"time（格式化 + ticker）",
				"reflect（只读元数据）",
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
