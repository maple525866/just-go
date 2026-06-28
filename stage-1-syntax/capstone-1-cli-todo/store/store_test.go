package store

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"just-go/stage-1-syntax/capstone-1-cli-todo/todo"
)

func TestStoreLoadSave(t *testing.T) {
	now := time.Date(2026, 6, 28, 1, 2, 3, 0, time.UTC)
	tests := []struct {
		name string
		seed bool
	}{
		{name: "missing file loads empty"},
		{name: "round trip", seed: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := Store{Path: filepath.Join(t.TempDir(), "todos.json")}
			if tt.seed {
				list := todo.NewList()
				_, _ = list.Add("persist", now)
				if err := store.Save(list); err != nil {
					t.Fatalf("Save() unexpected error: %v", err)
				}
			}
			loaded, err := store.Load()
			if err != nil {
				t.Fatalf("Load() unexpected error: %v", err)
			}
			if tt.seed {
				if len(loaded.Tasks) != 1 || loaded.Tasks[0].Title != "persist" || loaded.Tasks[0].CreatedAt.IsZero() {
					t.Fatalf("Load() = %#v", loaded)
				}
			} else if len(loaded.Tasks) != 0 || loaded.NextID != 1 {
				t.Fatalf("Load() missing = %#v, want empty list", loaded)
			}
		})
	}
}

func TestStoreLoadInvalidJSON(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{name: "invalid json", data: "not-json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "todos.json")
			if err := osWriteFile(path, []byte(tt.data)); err != nil {
				t.Fatalf("write fixture: %v", err)
			}
			_, err := (Store{Path: path}).Load()
			if !errors.Is(err, ErrStore) {
				t.Fatalf("Load() error = %v, want ErrStore", err)
			}
		})
	}
}

var osWriteFile = func(path string, data []byte) error {
	return os.WriteFile(path, data, 0o644)
}
