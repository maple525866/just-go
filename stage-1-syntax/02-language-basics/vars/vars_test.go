package vars

import (
	"strings"
	"testing"
)

func TestFormatScore(t *testing.T) {
	cases := []struct {
		name string
		in   struct {
			subject string
			score   int
		}
		wantContains []string
	}{
		{
			name: "满分边界",
			in: struct {
				subject string
				score   int
			}{"Math", 100},
			wantContains: []string{"Math: 100", "100.0%", "[pass]"},
		},
		{
			name: "及格线",
			in: struct {
				subject string
				score   int
			}{"English", 60},
			wantContains: []string{"English: 60", "60.0%", "[pass]"},
		},
		{
			name: "不及格",
			in: struct {
				subject string
				score   int
			}{"Go", 59},
			wantContains: []string{"Go: 59", "59.0%", "[fail]"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := FormatScore(tc.in.subject, tc.in.score)
			for _, part := range tc.wantContains {
				if !strings.Contains(got, part) {
					t.Errorf("FormatScore(%q, %d) = %q，应包含 %q", tc.in.subject, tc.in.score, got, part)
				}
			}
		})
	}
}

func TestZeroValueDemo(t *testing.T) {
	got := ZeroValueDemo()
	wantParts := []string{"int=0", "bool=false", "string=\"\"", "byte=0"}
	for _, part := range wantParts {
		if !strings.Contains(got, part) {
			t.Errorf("ZeroValueDemo() = %q，应包含 %q", got, part)
		}
	}
}
