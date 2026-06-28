package moduleinfo

import (
	"reflect"
	"strings"
	"testing"
)

func TestConcepts(t *testing.T) {
	tests := []struct {
		name     string
		module   string
		contains []string
	}{
		{name: "engineering concepts", module: "just-go", contains: []string{"module", "go.work", "semantic version"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			concepts := Concepts(tt.module)
			joined := ""
			for _, concept := range concepts {
				joined += concept.Name + " " + concept.Summary + "\n"
			}
			for _, part := range tt.contains {
				if !strings.Contains(joined, part) {
					t.Fatalf("Concepts() = %#v, want text containing %q", concepts, part)
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
		{name: "concept names", want: []string{"module", "go.work", "semantic version"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Names(); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Names() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
