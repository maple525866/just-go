## 1. Chapter Structure and Dependencies

- [x] 1.1 Add chi and validator dependencies with `go get github.com/go-chi/chi/v5 github.com/go-playground/validator/v10`.
- [x] 1.2 Add `main.go` that starts the chapter HTTP server through the `server` package and respects `ADDR` with `:8080` default.
- [x] 1.3 Create focused subpackages `model/`, `store/`, `response/`, `validation/`, `middleware/`, and `server/`.

## 2. Models, Store, Validation, and Responses

- [x] 2.1 Implement article request/response and error response types in `model`.
- [x] 2.2 Implement an in-memory article store with deterministic seed data and table-driven tests.
- [x] 2.3 Implement JSON success/error helpers that set status codes and `Content-Type: application/json`.
- [x] 2.4 Implement validator-backed request validation with field-level validation error formatting and tests.

## 3. Middleware Examples

- [x] 3.1 Implement request ID context propagation middleware and tests for context value plus response header.
- [x] 3.2 Implement Recover middleware that converts panic into a JSON 500 response and tests it with `httptest`.
- [x] 3.3 Implement logging and CORS middleware and tests for CORS response headers.
- [x] 3.4 Implement a simple in-memory limiter middleware and tests for 429 responses after capacity is exhausted.

## 4. Router and Handler Examples

- [x] 4.1 Implement `server.NewStdMux` with at least `GET /healthz` for standard-library ServeMux demonstration.
- [x] 4.2 Implement `server.NewRouter` with chi routes for health, article list, article create, and article detail.
- [x] 4.3 Add `httptest` coverage for successful health/list/create/detail requests.
- [x] 4.4 Add `httptest` coverage for invalid JSON, validation failure, and not-found REST error responses.

## 5. Learning Materials and Verification

- [x] 5.1 Update `stage-2-business/08-web-foundations/README.md` to replace placeholder content with actual package list, API examples, run commands, and knowledge-aligned checklist.
- [x] 5.2 Add `stage-2-business/08-web-foundations/EXERCISES.md` with 3 to 5 exercises, each including explicit acceptance criteria.
- [x] 5.3 Run `go test ./stage-2-business/08-web-foundations/...`, `go run ./stage-2-business/08-web-foundations`, `go test ./...`, and `go build ./...`; fix any failures.
