## Context

Chapter 15 teaches resilience and performance after Chapter 14 established microservice boundaries. The implementation must remain deterministic and self-contained while still showing realistic strategy ordering and failure semantics.

## Goals / Non-Goals

**Goals:**

- Demonstrate local rate limiting, bulkhead isolation, retry/backoff/jitter, circuit breaking, fallback, timeout handling, pprof registration, and load-test guidance.
- Use a controllable fake upstream so success and failure paths can be tested without external infrastructure.
- Use `github.com/sony/gobreaker/v2` behind a narrow wrapper to show third-party production-library integration.
- Document production-grade alternatives and explain why this chapter hand-writes several mechanisms for practice.

**Non-Goals:**

- Provide a production-ready distributed limiter, service mesh, API gateway, or observability platform.
- Require Redis, Envoy, Kubernetes, vegeta, wrk, hey, or Docker in default tests.
- Modify Chapter 14 code or implement Capstone 3.

## Decisions

### Keep Chapter 15 self-contained

The chapter uses the Chapter 14 product/stock concept but has its own fake upstream and Gateway. This preserves the completed Chapter 14 implementation and lets Chapter 15 focus on resilience behavior.

### Teach core mechanisms with small packages

Limiter, retry, bulkhead, upstream, gateway, and profiler packages stay focused and testable. Time and random behavior are injectable where the chapter owns the implementation.

### Use gobreaker for circuit breaking

Circuit breaking is the one resilience primitive implemented with a production common library. The wrapper isolates library types from Gateway orchestration code.

### Treat performance as observation, not fixed thresholds

Default tests verify pprof registration and workload execution, not machine-specific latency, QPS, or allocation budgets. README and exercises show how to run vegeta/wrk/hey and pprof manually.

## Risks / Trade-offs

- Hand-written examples may look production-ready; README must explicitly compare with industrial options.
- Combining many policies can obscure ordering; Gateway tests must lock in the order: limit, bulkhead, timeout, retry, breaker, fallback.
- Breaker timing comes from the library clock; tests should avoid tight timing thresholds.
