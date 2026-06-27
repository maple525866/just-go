package stream

import (
	"reflect"
	"testing"
)

func TestCopyText(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantText  string
		wantBytes int64
	}{
		{name: "copies text", input: "hello", wantText: "hello", wantBytes: 5},
		{name: "copies empty", input: "", wantText: "", wantBytes: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, n, err := CopyText(tt.input)
			if err != nil {
				t.Fatalf("CopyText() unexpected error: %v", err)
			}
			if got != tt.wantText || n != tt.wantBytes {
				t.Fatalf("CopyText() = (%q, %d), want (%q, %d)", got, n, tt.wantText, tt.wantBytes)
			}
		})
	}
}

func TestScanLines(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{name: "multiple lines", input: "a\nb\nc", want: []string{"a", "b", "c"}},
		{name: "empty", input: "", want: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ScanLines(tt.input)
			if err != nil {
				t.Fatalf("ScanLines() unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ScanLines() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
