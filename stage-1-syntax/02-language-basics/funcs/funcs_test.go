package funcs

import "testing"

func TestMinMax(t *testing.T) {
	cases := []struct {
		name    string
		nums    []int
		wantMin int
		wantMax int
		wantOK  bool
	}{
		{name: "空参数", nums: nil, wantOK: false},
		{name: "单个元素", nums: []int{42}, wantMin: 42, wantMax: 42, wantOK: true},
		{name: "多个元素", nums: []int{92, 78, 88}, wantMin: 78, wantMax: 92, wantOK: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotMin, gotMax, gotOK := MinMax(tc.nums...)
			if gotOK != tc.wantOK {
				t.Errorf("MinMax ok = %v, 期望 %v", gotOK, tc.wantOK)
			}
			if gotMin != tc.wantMin || gotMax != tc.wantMax {
				t.Errorf("MinMax = (%d, %d), 期望 (%d, %d)", gotMin, gotMax, tc.wantMin, tc.wantMax)
			}
		})
	}
}

func TestAverage(t *testing.T) {
	cases := []struct {
		name string
		nums []int
		want float64
	}{
		{name: "空参数", nums: nil, want: 0},
		{name: "三个分数", nums: []int{92, 78, 88}, want: 86},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Average(tc.nums...)
			if got != tc.want {
				t.Errorf("Average(%v) = %v, 期望 %v", tc.nums, got, tc.want)
			}
		})
	}
}

func TestMakeGrader(t *testing.T) {
	grader := MakeGrader(60)
	cases := []struct {
		name  string
		score int
		want  string
	}{
		{name: "及格", score: 60, want: "pass"},
		{name: "不及格", score: 59, want: "fail"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := grader(tc.score)
			if got != tc.want {
				t.Errorf("grader(%d) = %q, 期望 %q", tc.score, got, tc.want)
			}
		})
	}

	strict := MakeGrader(90)
	if strict(89) != "fail" || strict(90) != "pass" {
		t.Errorf("不同闭包实例应持有各自 threshold")
	}
}
