# Design: Chapter 11 Observability

## Overview

The chapter uses lightweight, in-process components so learners can run tests without external Prometheus or OpenTelemetry collectors. The API mirrors production concepts: structured log records contain trace IDs, metrics expose Prometheus text format, spans are propagated through context, and health probes separate liveness from readiness.

## Packages

- `loggingx`: helpers around `log/slog` for trace-aware structured logs.
- `metricsx`: concurrency-safe counter, gauge, and histogram registry with Prometheus exposition.
- `tracex`: minimal span context propagation model for teaching trace/span IDs.
- `healthx`: liveness/readiness checks and aggregate reports.
- `server`: HTTP routes that combine all observability pieces.
