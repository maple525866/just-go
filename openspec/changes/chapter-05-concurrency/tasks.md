## 1. Chapter Structure

- [x] 1.1 Remove the placeholder-only state for `stage-1-syntax/05-concurrency/` and add `main.go` that imports at least three topic subpackages and prints a concurrency report.
- [x] 1.2 Add package directories `goroutine/`, `channel/`, `syncx/`, `ctx/`, and `pitfall/` with focused exported functions whose behavior can be asserted by tests.

## 2. Goroutine Examples

- [x] 2.1 Implement `goroutine` examples for starting multiple goroutines, waiting for completion, and collecting deterministic results.
- [x] 2.2 Add table-driven `goroutine` tests using `t.Run` that verify task completion without leaking background goroutines.

## 3. Channel Examples

- [x] 3.1 Implement `channel` examples for unbuffered send/receive and buffered channel behavior.
- [x] 3.2 Implement `channel` examples for close + range and `select` timeout.
- [x] 3.3 Add table-driven `channel` tests using `t.Run` for unbuffered communication, close/range, and timeout behavior.

## 4. sync Examples

- [x] 4.1 Implement `syncx` examples for `sync.Mutex` protected shared state and `sync.WaitGroup` coordination.
- [x] 4.2 Implement `syncx` examples for `sync.RWMutex` reads and `sync.Once` one-time initialization.
- [x] 4.3 Add table-driven `syncx` tests using `t.Run` for Mutex/RWMutex/WaitGroup/Once behavior.

## 5. context and Pitfall Examples

- [x] 5.1 Implement `ctx` examples for cancellation, timeout, and cooperative worker exit.
- [x] 5.2 Implement `pitfall` examples that safely summarize data race, goroutine leak, and channel deadlock risks without introducing race or hangs.
- [x] 5.3 Add table-driven `ctx` and `pitfall` tests using `t.Run` for cancellation, timeout, and pitfall summaries.

## 6. Learning Materials and Verification

- [x] 6.1 Update `stage-1-syntax/05-concurrency/README.md` to replace all placeholder content with actual file list, run commands, and a knowledge-aligned checklist.
- [x] 6.2 Add `stage-1-syntax/05-concurrency/EXERCISES.md` with 3 to 5 exercises, each including explicit acceptance criteria.
- [x] 6.3 Run `go test -race ./stage-1-syntax/05-concurrency/...`, `go run ./stage-1-syntax/05-concurrency`, `go test ./...`, and `go build ./...`; fix any failures.
