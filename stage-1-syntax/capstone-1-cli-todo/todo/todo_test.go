package todo

import (
	"errors"
	"testing"
	"time"
)

func TestListAdd(t *testing.T) {
	now := time.Date(2026, 6, 28, 1, 2, 3, 0, time.UTC)
	tests := []struct {
		name    string
		title   string
		wantErr error
	}{
		{name: "adds task", title: "write tests"},
		{name: "missing title", title: "  ", wantErr: ErrMissingTitle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewList()
			task, err := list.Add(tt.title, now)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Add() err = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Add() unexpected error: %v", err)
			}
			if task.ID != 1 || task.Title != tt.title || task.Done || task.CreatedAt.IsZero() || task.UpdatedAt.IsZero() {
				t.Fatalf("Add() task = %#v", task)
			}
		})
	}
}

func TestListMutations(t *testing.T) {
	now := time.Date(2026, 6, 28, 1, 2, 3, 0, time.UTC)
	later := now.Add(time.Hour)
	tests := []struct {
		name string
		act  func(*List) error
	}{
		{
			name: "mark done",
			act: func(l *List) error {
				_, _ = l.Add("ship", now)
				task, err := l.MarkDone(1, later)
				if err != nil {
					return err
				}
				if !task.Done || !task.UpdatedAt.Equal(later) {
					t.Fatalf("MarkDone() task = %#v", task)
				}
				return nil
			},
		},
		{
			name: "delete",
			act: func(l *List) error {
				_, _ = l.Add("ship", now)
				deleted, err := l.Delete(1)
				if err != nil {
					return err
				}
				if deleted.ID != 1 || len(l.Tasks) != 0 {
					t.Fatalf("Delete() = %#v len=%d", deleted, len(l.Tasks))
				}
				return nil
			},
		},
		{
			name: "clear",
			act: func(l *List) error {
				_, _ = l.Add("a", now)
				_, _ = l.Add("b", now)
				if got := l.Clear(); got != 2 || len(l.Tasks) != 0 || l.NextID != 1 {
					t.Fatalf("Clear() = %d len=%d next=%d", got, len(l.Tasks), l.NextID)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewList()
			if err := tt.act(&list); err != nil {
				t.Fatalf("mutation returned error: %v", err)
			}
		})
	}
}

func TestListErrors(t *testing.T) {
	tests := []struct {
		name    string
		act     func(*List) error
		wantErr error
	}{
		{name: "invalid done id", act: func(l *List) error { _, err := l.MarkDone(0, time.Now()); return err }, wantErr: ErrInvalidID},
		{name: "missing done id", act: func(l *List) error { _, err := l.MarkDone(99, time.Now()); return err }, wantErr: ErrTaskNotFound},
		{name: "invalid delete id", act: func(l *List) error { _, err := l.Delete(-1); return err }, wantErr: ErrInvalidID},
		{name: "missing delete id", act: func(l *List) error { _, err := l.Delete(99); return err }, wantErr: ErrTaskNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewList()
			if err := tt.act(&list); !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
