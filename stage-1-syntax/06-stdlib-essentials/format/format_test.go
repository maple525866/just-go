package format

import "testing"

func TestProgress(t *testing.T) {
	tests := []struct {
		name    string
		chapter string
		done    int
		total   int
		want    string
	}{
		{name: "pads counts", chapter: "stdlib", done: 3, total: 12, want: "stdlib: 03/12 tasks"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Progress(tt.chapter, tt.done, tt.total); got != tt.want {
				t.Fatalf("Progress() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWriteSummary(t *testing.T) {
	tests := []struct {
		name  string
		topic string
		count int
		want  string
	}{
		{name: "writes to buffer", topic: "fmt", count: 2, want: "topic=fmt examples=2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WriteSummary(tt.topic, tt.count); got != tt.want {
				t.Fatalf("WriteSummary() = %q, want %q", got, tt.want)
			}
		})
	}
}
