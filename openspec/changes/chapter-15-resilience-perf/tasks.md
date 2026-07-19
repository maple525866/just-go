## 1. Upstream and OpenSpec Baseline

- [x] 1.1 Add OpenSpec proposal, design, tasks, and metadata for `chapter-15-resilience-perf`.
- [x] 1.2 Add the upstream product DTO, error classification helpers, HTTP client, and scripted fake upstream.
- [x] 1.3 Verify upstream tests pass without external services.

## 2. Local Overload Protection

- [x] 2.1 Add a time-injectable token bucket limiter with deterministic tests.
- [x] 2.2 Add a semaphore bulkhead with context-aware acquire/release tests.

## 3. Dependency Failure Handling

- [x] 3.1 Add retry with exponential backoff, jitter injection, cancellation, and error classification tests.
- [x] 3.2 Add a `gobreaker/v2` wrapper and tests for opening and rejecting calls.

## 4. Resilience Gateway

- [x] 4.1 Add Gateway success, rate-limit, bulkhead, retry, breaker, timeout, fallback, and error-mapping tests.
- [x] 4.2 Implement Gateway orchestration in the order limit → bulkhead → timeout → retry → breaker → upstream → fallback.

## 5. Profiling and Runnable Example

- [x] 5.1 Add pprof registration and heap workload tests.
- [x] 5.2 Add a runnable demo that starts fake upstream and Gateway servers, performs representative requests, and shuts down.

## 6. Learning Materials and Verification

- [x] 6.1 Replace the Chapter 15 README with Chinese learning materials, diagrams, commands, pprof workflow, and production alternatives.
- [x] 6.2 Add measurable exercises for distributed limiting, breaker tuning, fallback/caching, pprof, load testing, and service mesh migration.
- [x] 6.3 Update `ROADMAP.md` output and progress while leaving Capstone 3 incomplete.
- [x] 6.4 Run gofmt, Chapter 15 tests, full repository tests, race tests, vet, build, golangci-lint if available, OpenSpec validation if available, and independent code review.
