// Package asyncsave demonstrates controlled asynchronous JSON persistence.
//
// import 路径：just-go/stage-1-syntax/capstone-1-cli-todo/asyncsave
package asyncsave

import (
	"errors"
	"sync"

	"just-go/stage-1-syntax/capstone-1-cli-todo/store"
	"just-go/stage-1-syntax/capstone-1-cli-todo/todo"
)

var ErrClosed = errors.New("async save worker is closed")

type request struct {
	list todo.List
	ack  chan error
}

// Worker serializes save requests through one goroutine.
type Worker struct {
	store  store.Store
	jobs   chan request
	done   chan struct{}
	closed bool
	mu     sync.Mutex
}

// New starts an async save worker.
func New(s store.Store) *Worker {
	w := &Worker{store: s, jobs: make(chan request), done: make(chan struct{})}
	go w.loop()
	return w
}

func (w *Worker) loop() {
	defer close(w.done)
	for req := range w.jobs {
		req.ack <- w.store.Save(req.list)
		close(req.ack)
	}
}

// Submit queues one save request and waits for that save to finish.
func (w *Worker) Submit(list todo.List) error {
	w.mu.Lock()
	closed := w.closed
	w.mu.Unlock()
	if closed {
		return ErrClosed
	}
	ack := make(chan error, 1)
	w.jobs <- request{list: list, ack: ack}
	return <-ack
}

// Close stops the worker after queued work is processed.
func (w *Worker) Close() error {
	w.mu.Lock()
	if w.closed {
		w.mu.Unlock()
		<-w.done
		return nil
	}
	w.closed = true
	close(w.jobs)
	w.mu.Unlock()
	<-w.done
	return nil
}
