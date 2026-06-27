package generic

import (
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []string
	}{
		{name: "numbers to labels", in: []int{1, 2, 3}, want: []string{"#1", "#2", "#3"}},
		{name: "empty", in: []int{}, want: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Map(tt.in, func(n int) string { return "#" + string(rune('0'+n)) })
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Map() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{name: "keep long words", in: []string{"Go", "error", "generic"}, want: []string{"error", "generic"}},
		{name: "none kept", in: []string{"a", "bb"}, want: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Filter(tt.in, func(s string) bool { return len(s) > 2 })
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Filter() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want int
	}{
		{name: "sum ints", in: []int{1, 2, 3}, want: 6},
		{name: "empty", in: []int{}, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.in); got != tt.want {
				t.Fatalf("Sum() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestDemo(t *testing.T) {
	tests := []struct {
		name          string
		wantDoubled   []int
		wantLongNames []string
		wantTotal     int
	}{
		{name: "demo values", wantDoubled: []int{2, 4, 6}, wantLongNames: []string{"error", "generic"}, wantTotal: 60},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doubled, longNames, total := Demo()
			if !reflect.DeepEqual(doubled, tt.wantDoubled) {
				t.Fatalf("Demo() doubled = %#v, want %#v", doubled, tt.wantDoubled)
			}
			if !reflect.DeepEqual(longNames, tt.wantLongNames) {
				t.Fatalf("Demo() longNames = %#v, want %#v", longNames, tt.wantLongNames)
			}
			if total != tt.wantTotal {
				t.Fatalf("Demo() total = %d, want %d", total, tt.wantTotal)
			}
		})
	}
}
