package quality

import (
	"reflect"
	"testing"
)

func TestCommands(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{name: "ci aligned commands", want: []string{"go vet ./...", "go test -race -count=1 ./...", "go build ./...", "golangci-lint run"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Commands(); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Commands() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestChecks(t *testing.T) {
	tests := []struct {
		name    string
		wantLen int
	}{
		{name: "has four checks", wantLen: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Checks(); len(got) != tt.wantLen {
				t.Fatalf("Checks() len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}
