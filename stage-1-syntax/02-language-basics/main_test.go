package main

import (
	"strings"
	"testing"
)

func TestReportTitle(t *testing.T) {
	cases := []struct {
		name    string
		student string
		want    string
	}{
		{name: "指定学生", student: "Ada", want: "=== 语法基础报告：Ada ==="},
		{name: "空名回退", student: "", want: "=== 语法基础报告：Student ==="},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := reportTitle(tc.student)
			if got != tc.want {
				t.Errorf("reportTitle(%q) = %q, 期望 %q", tc.student, got, tc.want)
			}
		})
	}
}

func TestBuildReportContainsSections(t *testing.T) {
	got := buildReport("Ada")
	sections := []string{
		"=== 语法基础报告：Ada ===",
		"## 变量与类型（vars）",
		"## 控制流（control）",
		"## 函数（funcs）",
		"--- end of report ---",
	}
	for _, s := range sections {
		if !strings.Contains(got, s) {
			t.Errorf("buildReport 应包含 %q，got:\n%s", s, got)
		}
	}
}
