package pitfall

import (
	"reflect"
	"strings"
	"testing"
)

func TestSummaries(t *testing.T) {
	tests := []struct {
		name     string
		contains []string
	}{
		{name: "covers common pitfalls", contains: []string{"data race", "goroutine leak", "channel deadlock"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Summaries()
			joined := ""
			for _, item := range got {
				joined += item.Name + " " + item.Risk + " " + item.Prevention + "\n"
			}
			for _, part := range tt.contains {
				if !strings.Contains(joined, part) {
					t.Fatalf("Summaries() = %#v, want text containing %q", got, part)
				}
			}
		})
	}
}

func TestNames(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{name: "pitfall names", want: []string{"data race", "goroutine leak", "channel deadlock"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Names(); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Names() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
