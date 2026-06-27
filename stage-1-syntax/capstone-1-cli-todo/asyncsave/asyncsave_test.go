package asyncsave

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"just-go/stage-1-syntax/capstone-1-cli-todo/store"
	"just-go/stage-1-syntax/capstone-1-cli-todo/todo"
)

func TestWorkerSubmitAndClose(t *testing.T) {
	now := time.Date(2026, 6, 28, 1, 2, 3, 0, time.UTC)
	tests := []struct {
		name string
	}{
		{name: "writes and closes"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := store.Store{Path: filepath.Join(t.TempDir(), "todos.json")}
			list := todo.NewList()
			_, _ = list.Add("async", now)
			worker := New(st)
			if err := worker.Submit(list); err != nil {
				t.Fatalf("Submit() unexpected error: %v", err)
			}
			if err := worker.Close(); err != nil {
				t.Fatalf("Close() unexpected error: %v", err)
			}
			loaded, err := st.Load()
			if err != nil {
				t.Fatalf("Load() unexpected error: %v", err)
			}
			if len(loaded.Tasks) != 1 || loaded.Tasks[0].Title != "async" {
				t.Fatalf("loaded = %#v", loaded)
			}
			if err := worker.Submit(list); !errors.Is(err, ErrClosed) {
				t.Fatalf("Submit() after close = %v, want ErrClosed", err)
			}
		})
	}
}
