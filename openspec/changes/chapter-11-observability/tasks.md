## 1. Chapter Structure

- [x] 1.1 Add `main.go` that prints an observability learning report through the chapter packages.
- [x] 1.2 Create focused subpackages `loggingx/`, `metricsx/`, `tracex/`, `healthx/`, and `server/`.

## 2. Structured Logging and Trace Context

- [x] 2.1 Implement trace ID context helpers and span creation in `tracex`.
- [x] 2.2 Implement `slog` helpers that include trace ID and stable fields in JSON logs.
- [x] 2.3 Add tests for trace propagation and trace-aware structured log output.

## 3. Metrics and Health Checks

- [x] 3.1 Implement a Prometheus-style registry with Counter, Gauge, and Histogram.
- [x] 3.2 Add tests for metric mutation and text exposition.
- [x] 3.3 Implement liveness/readiness checks with aggregate reports and tests.

## 4. HTTP Integration

- [x] 4.1 Implement routes for `/livez`, `/readyz`, `/metrics`, and `/work`.
- [x] 4.2 Add middleware or handler logic that records request metrics, logs trace ID, and opens a span.
- [x] 4.3 Add `httptest` coverage for health, metrics, successful work, and readiness failure.

## 5. Learning Materials and Verification

- [x] 5.1 Update README with package list, run commands, sample endpoints, and a filled checklist.
- [x] 5.2 Add EXERCISES with 3 to 5 exercises and explicit acceptance criteria.
- [x] 5.3 Run chapter tests, chapter demo, full tests, vet, and build; fix failures.
