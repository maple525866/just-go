package dict

import "testing"

func TestLookup(t *testing.T) {
	scores := map[string]int{
		"Ada": 92,
		"Zoe": 0, // 值恰为零值，用于区分"零值"与"缺失键"
	}

	cases := []struct {
		name      string
		key       string
		wantScore int
		wantFound bool
	}{
		{name: "存在的键", key: "Ada", wantScore: 92, wantFound: true},
		{name: "值为零值的键", key: "Zoe", wantScore: 0, wantFound: true},
		{name: "缺失的键", key: "Nobody", wantScore: 0, wantFound: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotScore, gotFound := Lookup(scores, tc.key)
			if gotScore != tc.wantScore || gotFound != tc.wantFound {
				t.Errorf("Lookup(%q) = (%d, %v), 期望 (%d, %v)",
					tc.key, gotScore, gotFound, tc.wantScore, tc.wantFound)
			}
		})
	}
}

func TestTotal(t *testing.T) {
	cases := []struct {
		name   string
		scores map[string]int
		want   int
	}{
		{name: "空 map", scores: map[string]int{}, want: 0},
		{name: "nil map", scores: nil, want: 0},
		{name: "多个分数", scores: map[string]int{"a": 92, "b": 78, "c": 88}, want: 258},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Total(tc.scores); got != tc.want {
				t.Errorf("Total = %d, 期望 %d", got, tc.want)
			}
		})
	}
}

func TestCountAtLeast(t *testing.T) {
	scores := map[string]int{"a": 92, "b": 78, "c": 59}
	cases := []struct {
		name      string
		threshold int
		want      int
	}{
		{name: "及格线 60", threshold: 60, want: 2},
		{name: "优秀线 90", threshold: 90, want: 1},
		{name: "全部 0", threshold: 0, want: 3},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := CountAtLeast(scores, tc.threshold); got != tc.want {
				t.Errorf("CountAtLeast(%d) = %d, 期望 %d", tc.threshold, got, tc.want)
			}
		})
	}
}
