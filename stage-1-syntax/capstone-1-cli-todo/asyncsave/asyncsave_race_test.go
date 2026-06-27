package asyncsave

import (
	"errors"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"just-go/stage-1-syntax/capstone-1-cli-todo/store"
	"just-go/stage-1-syntax/capstone-1-cli-todo/todo"
)

func TestWorkerSubmitCloseRaceDoesNotPanic(t *testing.T) {
	now := time.Date(2026, 6, 28, 1, 2, 3, 0, time.UTC)
	tests := []struct {
		name string
	}{
		{name: "concurrent submit and close"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := store.Store{Path: filepath.Join(t.TempDir(), "todos.json")}
			list := todo.NewList()
			_, _ = list.Add("race", now)
			worker := New(st)

			var wg sync.WaitGroup
			errs := make(chan error, 2)
			wg.Add(2)
			go func() {
				defer wg.Done()
				errs <- worker.Submit(list)
			}()
			go func() {
				defer wg.Done()
				errs <- worker.Close()
			}()
			wg.Wait()
			close(errs)

			for err := range errs {
				if err != nil && !errors.Is(err, ErrClosed) {
					t.Fatalf("concurrent Submit/Close error = %v, want nil or ErrClosed", err)
				}
			}
		})
	}
}
