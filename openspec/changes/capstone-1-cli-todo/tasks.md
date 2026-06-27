## 1. Project Structure and Domain

- [x] 1.1 Add `main.go` and packages `app/`, `todo/`, `store/`, and `asyncsave/` under `stage-1-syntax/capstone-1-cli-todo/`.
- [x] 1.2 Implement `todo` task/list types with id, title, done state, created/updated timestamps, and operations for add/list/done/delete/clear.
- [x] 1.3 Implement custom or sentinel errors for invalid command, missing title, invalid id, task not found, and storage failures.

## 2. Persistence and Async Save

- [x] 2.1 Implement `store` JSON load/save with configurable file path, missing-file-as-empty behavior, and time-preserving round trip.
- [x] 2.2 Implement `asyncsave` goroutine + channel save worker with submit and close behavior and no goroutine leak.
- [x] 2.3 Add table-driven `todo`, `store`, and `asyncsave` tests using `t.Run`.

## 3. CLI App Behavior

- [x] 3.1 Implement `app.Run(args, stdout, stderr)` command parsing for `add`, `list`, `done`, `delete`, `clear`, and `help` / `--help`.
- [x] 3.2 Implement human-readable list output and command success/error messages.
- [x] 3.3 Add table-driven CLI tests using temp JSON files and in-memory stdout/stderr buffers.

## 4. Benchmark and Verification

- [x] 4.1 Add at least one benchmark for a meaningful operation such as list rendering or adding tasks.
- [x] 4.2 Update `stage-1-syntax/capstone-1-cli-todo/README.md` with implemented functionality, file layout, example commands, data file behavior, and verification commands.
- [x] 4.3 Add `stage-1-syntax/capstone-1-cli-todo/EXERCISES.md` with 3 to 5 extension exercises, each including explicit acceptance criteria.
- [x] 4.4 Run `go test ./stage-1-syntax/capstone-1-cli-todo/...`, `go test -bench=. ./stage-1-syntax/capstone-1-cli-todo/...`, `go run ./stage-1-syntax/capstone-1-cli-todo --help`, `go test ./...`, and `go build ./...`; fix any failures.
