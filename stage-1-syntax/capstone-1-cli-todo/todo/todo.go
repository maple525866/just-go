// Package todo contains the domain model and operations for the CLI Todo capstone.
//
// import 路径：just-go/stage-1-syntax/capstone-1-cli-todo/todo
package todo

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	// ErrMissingTitle means an add command did not provide a usable title.
	ErrMissingTitle = errors.New("todo title is required")
	// ErrTaskNotFound means no task exists with the requested id.
	ErrTaskNotFound = errors.New("todo task not found")
	// ErrInvalidID means a command provided an invalid task id.
	ErrInvalidID = errors.New("todo id is invalid")
)

// Task is one todo item persisted to JSON.
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// List is an ordered collection of tasks.
type List struct {
	Tasks  []Task `json:"tasks"`
	NextID int    `json:"next_id"`
}

// NewList creates an empty list with ids starting at 1.
func NewList() List {
	return List{Tasks: []Task{}, NextID: 1}
}

// Normalize ensures loaded lists have usable zero-value defaults.
func (l *List) Normalize() {
	if l.Tasks == nil {
		l.Tasks = []Task{}
	}
	maxID := 0
	for _, task := range l.Tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	if l.NextID <= maxID {
		l.NextID = maxID + 1
	}
	if l.NextID == 0 {
		l.NextID = 1
	}
}

// Add appends a new task with a stable id.
func (l *List) Add(title string, now time.Time) (Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Task{}, ErrMissingTitle
	}
	l.Normalize()
	task := Task{ID: l.NextID, Title: title, CreatedAt: now, UpdatedAt: now}
	l.NextID++
	l.Tasks = append(l.Tasks, task)
	return task, nil
}

// MarkDone marks an existing task as complete.
func (l *List) MarkDone(id int, now time.Time) (Task, error) {
	if id <= 0 {
		return Task{}, ErrInvalidID
	}
	for i := range l.Tasks {
		if l.Tasks[i].ID == id {
			l.Tasks[i].Done = true
			l.Tasks[i].UpdatedAt = now
			return l.Tasks[i], nil
		}
	}
	return Task{}, fmt.Errorf("%w: %d", ErrTaskNotFound, id)
}

// Delete removes an existing task.
func (l *List) Delete(id int) (Task, error) {
	if id <= 0 {
		return Task{}, ErrInvalidID
	}
	for i, task := range l.Tasks {
		if task.ID == id {
			l.Tasks = append(l.Tasks[:i], l.Tasks[i+1:]...)
			return task, nil
		}
	}
	return Task{}, fmt.Errorf("%w: %d", ErrTaskNotFound, id)
}

// Clear removes all tasks and resets the next id.
func (l *List) Clear() int {
	count := len(l.Tasks)
	l.Tasks = []Task{}
	l.NextID = 1
	return count
}

// Active returns unfinished tasks.
func (l List) Active() []Task {
	out := make([]Task, 0, len(l.Tasks))
	for _, task := range l.Tasks {
		if !task.Done {
			out = append(out, task)
		}
	}
	return out
}
