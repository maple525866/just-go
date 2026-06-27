// Package store persists todo lists as JSON files.
//
// import 路径：just-go/stage-1-syntax/capstone-1-cli-todo/store
package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"just-go/stage-1-syntax/capstone-1-cli-todo/todo"
)

var ErrStore = errors.New("todo store error")

// Store reads and writes one JSON file.
type Store struct {
	Path string
}

// Load reads a todo list. Missing files are treated as an empty list.
func (s Store) Load() (todo.List, error) {
	data, err := os.ReadFile(s.Path)
	if errors.Is(err, os.ErrNotExist) {
		return todo.NewList(), nil
	}
	if err != nil {
		return todo.List{}, fmt.Errorf("%w: read %s: %v", ErrStore, s.Path, err)
	}
	if len(data) == 0 {
		return todo.NewList(), nil
	}
	var list todo.List
	if err := json.Unmarshal(data, &list); err != nil {
		return todo.List{}, fmt.Errorf("%w: decode %s: %v", ErrStore, s.Path, err)
	}
	list.Normalize()
	return list, nil
}

// Save writes a todo list as pretty JSON.
func (s Store) Save(list todo.List) error {
	list.Normalize()
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("%w: encode: %v", ErrStore, err)
	}
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o755); err != nil {
		return fmt.Errorf("%w: mkdir: %v", ErrStore, err)
	}
	if err := os.WriteFile(s.Path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("%w: write %s: %v", ErrStore, s.Path, err)
	}
	return nil
}
