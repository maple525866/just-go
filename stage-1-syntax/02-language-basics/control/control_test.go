package control

import (
	"strings"
	"testing"
)

func TestLetterGrade(t *testing.T) {
	cases := []struct {
		name  string
		score int
		want  string
	}{
		{name: "A 边界 90", score: 90, want: "A"},
		{name: "A 边界 89", score: 89, want: "B"},
		{name: "B 边界 80", score: 80, want: "B"},
		{name: "C 边界 70", score: 70, want: "C"},
		{name: "D 边界 60", score: 60, want: "D"},
		{name: "F 边界 59", score: 59, want: "F"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := LetterGrade(tc.score)
			if got != tc.want {
				t.Errorf("LetterGrade(%d) = %q, 期望 %q", tc.score, got, tc.want)
			}
		})
	}
}

func TestSummarize(t *testing.T) {
	cases := []struct {
		name   string
		scores []int
		want   string
	}{
		{name: "空列表", scores: nil, want: "count=0 sum=0 avg=0"},
		{name: "多个分数", scores: []int{92, 78, 88}, want: "count=3 sum=258 avg=86.0 detail=[92,78,88]"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Summarize(tc.scores)
			if got != tc.want {
				t.Errorf("Summarize(%v) = %q, 期望 %q", tc.scores, got, tc.want)
			}
		})
	}
}

func TestRunReport(t *testing.T) {
	body := "line1\nline2"
	got := RunReport(body)
	if !strings.Contains(got, body) {
		t.Errorf("RunReport 应包含正文，got = %q", got)
	}
	if !strings.Contains(got, reportFooter) {
		t.Errorf("RunReport 应通过 defer 追加页脚 %q，got = %q", reportFooter, got)
	}
}
