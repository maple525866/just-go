package app

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"just-go/stage-1-syntax/capstone-1-cli-todo/todo"
)

func TestRunCommands(t *testing.T) {
	tests := []struct {
		name       string
		commands   [][]string
		wantOut    []string
		wantErr    error
		wantErrOut string
		seedStore  string
	}{
		{
			name:     "add list done delete clear",
			commands: [][]string{{"add", "write", "tests"}, {"list"}, {"done", "1"}, {"delete", "1"}, {"clear"}},
			wantOut:  []string{"added #1 write tests", "[ ] #1 write tests", "done #1 write tests", "deleted #1 write tests", "cleared 0 tasks"},
		},
		{
			name:     "help",
			commands: [][]string{{"--help"}},
			wantOut:  []string{"Usage:", "add <title>", "list"},
		},
		{
			name:       "invalid command",
			commands:   [][]string{{"wat"}},
			wantErr:    ErrInvalidCommand,
			wantErrOut: "Usage:",
		},
		{
			name:       "invalid command ignores corrupt store",
			commands:   [][]string{{"wat"}},
			wantErr:    ErrInvalidCommand,
			wantErrOut: "Usage:",
			seedStore:  "not-json",
		},
		{
			name:     "missing title",
			commands: [][]string{{"add"}},
			wantErr:  todo.ErrMissingTitle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "todos.json")
			t.Setenv(dataFileEnv, path)
			if tt.seedStore != "" {
				if err := os.WriteFile(path, []byte(tt.seedStore), 0o644); err != nil {
					t.Fatalf("seed store: %v", err)
				}
			}
			var out, errOut bytes.Buffer
			var err error
			for _, command := range tt.commands {
				err = Run(command, &out, &errOut)
				if err != nil {
					break
				}
			}
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Run() error = %v, want %v", err, tt.wantErr)
				}
				if tt.wantErrOut != "" && !strings.Contains(errOut.String(), tt.wantErrOut) {
					t.Fatalf("stderr = %q, want containing %q", errOut.String(), tt.wantErrOut)
				}
				return
			}
			if err != nil {
				t.Fatalf("Run() unexpected error: %v", err)
			}
			for _, part := range tt.wantOut {
				if !strings.Contains(out.String(), part) {
					t.Fatalf("stdout = %q, want containing %q", out.String(), part)
				}
			}
		})
	}
}

func TestRender(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "empty", want: "No todos.\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Render(todo.NewList()); got != tt.want {
				t.Fatalf("Render() = %q, want %q", got, tt.want)
			}
		})
	}
}

func BenchmarkRender(b *testing.B) {
	list := todo.NewList()
	for i := 0; i < 100; i++ {
		_, _ = list.Add("benchmark task", time.Unix(int64(i), 0))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Render(list)
	}
}
