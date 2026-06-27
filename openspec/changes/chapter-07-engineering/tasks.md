## 1. Chapter Structure

- [x] 1.1 Remove the placeholder-only state for `stage-1-syntax/07-engineering/` and add `main.go` that imports at least three topic subpackages and prints an engineering report.
- [x] 1.2 Add package directories `moduleinfo/`, `calc/`, `quality/`, `debugx/`, and `profile/` with focused exported functions whose behavior can be asserted by tests.

## 2. Module and Testing Examples

- [x] 2.1 Implement `moduleinfo` examples for module, go.work, and semantic version summaries.
- [x] 2.2 Implement `calc` pure functions suitable for table-driven tests and subtests.
- [x] 2.3 Add table-driven `moduleinfo` and `calc` tests using `t.Run`.

## 3. Benchmark and Quality Examples

- [x] 3.1 Add at least one `BenchmarkXxx` function for a `calc` pure function.
- [x] 3.2 Implement `quality` examples that return local verification commands for `go vet`, `go test -race`, `go build`, and `golangci-lint`.
- [x] 3.3 Add table-driven `quality` tests using `t.Run`.

## 4. Debug and Profile Examples

- [x] 4.1 Implement `debugx` examples for `log/slog` output and dlv/IDE debugging command summaries.
- [x] 4.2 Implement `profile` examples for pprof CPU, memory, and blocking profile summaries plus command hints.
- [x] 4.3 Add table-driven `debugx` and `profile` tests using `t.Run`.

## 5. Learning Materials and Verification

- [x] 5.1 Update `stage-1-syntax/07-engineering/README.md` to replace all placeholder content with actual file list, run commands, benchmark commands, and a knowledge-aligned checklist.
- [x] 5.2 Add `stage-1-syntax/07-engineering/EXERCISES.md` with 3 to 5 exercises, each including explicit acceptance criteria.
- [x] 5.3 Run `go test ./stage-1-syntax/07-engineering/...`, `go test -bench=. ./stage-1-syntax/07-engineering/...`, `go run ./stage-1-syntax/07-engineering`, `go test ./...`, and `go build ./...`; fix any failures.
