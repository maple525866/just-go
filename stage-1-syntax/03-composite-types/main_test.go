package main

import (
	"strings"
	"testing"
)

func TestReportTitle(t *testing.T) {
	cases := []struct {
		name  string
		class string
		want  string
	}{
		{name: "指定班级", class: "三年二班", want: "=== 班级花名册报告：三年二班 ==="},
		{name: "空名回退", class: "", want: "=== 班级花名册报告：Class ==="},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := reportTitle(tc.class)
			if got != tc.want {
				t.Errorf("reportTitle(%q) = %q, 期望 %q", tc.class, got, tc.want)
			}
		})
	}
}

func TestBuildReportContainsSections(t *testing.T) {
	got := buildReport("三年二班")
	sections := []string{
		"=== 班级花名册报告：三年二班 ===",
		"## 学生（model：struct + 嵌入）",
		"## 分数表（dict：map + comma-ok）",
		"## 切片（seq：append 扩容 + 共享底层数组）",
		"## 指针（ptr：值接收者 vs 指针接收者）",
		"Zoe: (未登记，comma-ok 区分缺失键)",
	}
	for _, s := range sections {
		if !strings.Contains(got, s) {
			t.Errorf("buildReport 应包含 %q，got:\n%s", s, got)
		}
	}
}
