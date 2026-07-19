## Why

Stage 3 needs a final architecture chapter that moves learners from microservices that can run to services that remain controlled under overload, latency, and dependency failure. Chapter 15 replaces the resilience/performance placeholder with a tested, self-contained teaching slice that builds on the Chapter 14 product/inventory concepts without modifying Chapter 14.

Tracking issue: [#35](https://github.com/maple525866/just-go/issues/35)

## What Changes

- Add a runnable Chapter 15 product-details Gateway with token bucket limiting, bulkhead isolation, retry with exponential backoff and jitter, circuit breaking, fallback, timeouts, and pprof hooks.
- Add deterministic fake upstream behavior for success, slow responses, intermittent failures, sustained failures, and client errors.
- Add tests, Chinese learning materials, measurable exercises, and roadmap progress updates.
- Document that the hand-written components are for practice and compare them with production-grade libraries, gateways, service meshes, distributed limiters, observability, and load-testing tools.

## Capabilities

### New Capabilities

- `resilience-performance-tutorial`: Defines the executable Chapter 15 example, resilience policies, performance profiling entry points, tests, and learning materials.

### Modified Capabilities

- `learning-curriculum`: Marks Chapter 15 as implemented and replaces its placeholder output with concrete learning deliverables.

## Impact

- Adds Go code, tests, README, and exercises under `stage-3-architecture/15-resilience-perf/`.
- Adds `github.com/sony/gobreaker/v2` as the circuit-breaker dependency.
- Updates `ROADMAP.md` Chapter 15 output and progress checkbox.
- Requires no external services or load-testing tools for default test execution.
