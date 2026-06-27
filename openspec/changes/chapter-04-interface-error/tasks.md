## 1. Chapter Structure

- [x] 1.1 Remove the placeholder-only state for `stage-1-syntax/04-interface-error/` and add `main.go` that imports at least two topic subpackages and prints an interface/error/generic report.
- [x] 1.2 Add package directories `iface/`, `apperr/`, and `generic/` with focused exported functions whose behavior can be asserted by tests.

## 2. Interface Examples

- [x] 2.1 Implement `iface` examples for implicit interface satisfaction and small-interface design using an exported function that accepts an interface and returns a concrete result.
- [x] 2.2 Implement `iface` examples for `any`, type assertion, and type switch with at least three distinguishable input categories.
- [x] 2.3 Add table-driven `iface` tests using `t.Run` for implicit implementation, small-interface behavior, and `any` classification.

## 3. Error Examples

- [x] 3.1 Implement `apperr` examples for sentinel errors, `%w` wrapping, and `errors.Is`.
- [x] 3.2 Implement `apperr` examples for a custom error type and `errors.As`, plus an exported summary containing error handling key terms.
- [x] 3.3 Add table-driven `apperr` tests using `t.Run` for `errors.Is`, `errors.As`, and the summary output.

## 4. Generic Examples

- [x] 4.1 Implement `generic` examples for type-parameter based `Map` and `Filter` functions over slices.
- [x] 4.2 Implement at least one type-set constraint and a constrained exported function that works for numeric types.
- [x] 4.3 Add table-driven `generic` tests using `t.Run` for `Map`, `Filter`, and the constrained numeric function.

## 5. Learning Materials and Verification

- [x] 5.1 Update `stage-1-syntax/04-interface-error/README.md` to replace all placeholder content with actual file list, run commands, and a knowledge-aligned checklist.
- [x] 5.2 Add `stage-1-syntax/04-interface-error/EXERCISES.md` with 3 to 5 exercises, each including explicit acceptance criteria.
- [x] 5.3 Run `go test ./stage-1-syntax/04-interface-error/...`, `go run ./stage-1-syntax/04-interface-error`, `go test ./...`, and `go build ./...`; fix any failures.
