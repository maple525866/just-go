package calc

import (
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{name: "positive", a: 2, b: 3, want: 5},
		{name: "negative", a: -2, b: 1, want: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.a, tt.b); got != tt.want {
				t.Fatalf("Add() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestFibonacci(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{name: "zero", n: 0, want: 0},
		{name: "one", n: 1, want: 1},
		{name: "ten", n: 10, want: 55},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Fibonacci(tt.n); got != tt.want {
				t.Fatalf("Fibonacci() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestNormalizeWords(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{name: "lowercase and split", input: " Go  Test BENCH ", want: []string{"go", "test", "bench"}},
		{name: "empty", input: " ", want: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeWords(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("NormalizeWords() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func BenchmarkNormalizeWords(b *testing.B) {
	input := "Go testing benchmark pprof engineering quality gate"
	for i := 0; i < b.N; i++ {
		_ = NormalizeWords(input)
	}
}
