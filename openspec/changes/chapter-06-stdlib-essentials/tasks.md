## 1. Chapter Structure

- [x] 1.1 Remove the placeholder-only state for `stage-1-syntax/06-stdlib-essentials/` and add `main.go` that imports at least four topic subpackages and prints a standard-library report.
- [x] 1.2 Add package directories `format/`, `stream/`, `system/`, `web/`, `codec/`, `clock/`, and `inspect/` with focused exported functions whose behavior can be asserted by tests.

## 2. Formatting and Stream Examples

- [x] 2.1 Implement `format` examples for `fmt.Sprintf` and `fmt.Fprintf` formatting.
- [x] 2.2 Implement `stream` examples for `io.Reader` / `io.Writer`, `io.Copy`, and `bufio.Scanner`.
- [x] 2.3 Add table-driven `format` and `stream` tests using `t.Run`.

## 3. System and HTTP Examples

- [x] 3.1 Implement `system` examples for temporary file read/write and environment variable reading.
- [x] 3.2 Implement `system` examples for safe `os/exec` command execution.
- [x] 3.3 Implement `web` examples for `net/http` handler and client using local test servers.
- [x] 3.4 Add table-driven `system` and `web` tests using `t.Run`; web tests must use `httptest` and not depend on external internet.

## 4. Encoding, Time, and Reflect Examples

- [x] 4.1 Implement `codec` examples for JSON and XML round trips.
- [x] 4.2 Implement `clock` examples for time formatting, duration, and finite ticker usage.
- [x] 4.3 Implement `inspect` examples for read-only reflect type/field/tag inspection.
- [x] 4.4 Add table-driven `codec`, `clock`, and `inspect` tests using `t.Run`.

## 5. Learning Materials and Verification

- [x] 5.1 Update `stage-1-syntax/06-stdlib-essentials/README.md` to replace all placeholder content with actual file list, run commands, and a knowledge-aligned checklist.
- [x] 5.2 Add `stage-1-syntax/06-stdlib-essentials/EXERCISES.md` with 3 to 5 exercises, each including explicit acceptance criteria.
- [x] 5.3 Run `go test ./stage-1-syntax/06-stdlib-essentials/...`, `go run ./stage-1-syntax/06-stdlib-essentials`, `go test ./...`, and `go build ./...`; fix any failures.
