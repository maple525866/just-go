package iface

import "testing"

func TestBuildReport(t *testing.T) {
	tests := []struct {
		name string
		in   Describer
		want Report
	}{
		{
			name: "book implicitly implements describer",
			in:   Book{Title: "Go", Author: "Gopher"},
			want: Report{Kind: "describer", Text: "Go by Gopher"},
		},
		{
			name: "lesson implicitly implements describer",
			in:   Lesson{Name: "接口", Day: 4},
			want: Report{Kind: "describer", Text: "Day 4: 接口"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildReport(tt.in)
			if got != tt.want {
				t.Fatalf("BuildReport() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestClassifyAny(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want string
	}{
		{name: "string", in: "hello", want: "string:hello"},
		{name: "int", in: 7, want: "int:7"},
		{name: "describer", in: Book{Title: "Go"}, want: "describer:Go"},
		{name: "unknown", in: true, want: "unknown:bool"},
		{name: "nil", in: nil, want: "nil"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClassifyAny(tt.in); got != tt.want {
				t.Fatalf("ClassifyAny() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsDescriber(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want bool
	}{
		{name: "book is describer", in: Book{Title: "Go"}, want: true},
		{name: "lesson is describer", in: Lesson{Name: "泛型", Day: 4}, want: true},
		{name: "plain string is not describer", in: "Go", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDescriber(tt.in); got != tt.want {
				t.Fatalf("IsDescriber() = %t, want %t", got, tt.want)
			}
		})
	}
}
