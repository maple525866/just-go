package goroutine

import (
	"reflect"
	"testing"
)

func TestRunSquares(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []Result
	}{
		{name: "three tasks", in: []int{2, 3, 4}, want: []Result{{ID: 0, Value: 4}, {ID: 1, Value: 9}, {ID: 2, Value: 16}}},
		{name: "empty", in: []int{}, want: []Result{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RunSquares(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("RunSquares() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestLifecycleSummary(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{name: "summary", want: []string{"start goroutine", "do work", "signal done", "wait and collect"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LifecycleSummary(); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("LifecycleSummary() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
