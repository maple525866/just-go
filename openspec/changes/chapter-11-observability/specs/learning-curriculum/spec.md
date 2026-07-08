## ADDED Requirements

### Requirement: Chapter 11 observability learning unit
The curriculum SHALL provide an executable Chapter 11 under `stage-2-business/11-observability/` that teaches structured logging, metrics, traces, and health checks.

#### Scenario: Learner runs the observability demo
- **WHEN** a learner runs `go run ./stage-2-business/11-observability`
- **THEN** the program prints a concise report describing log, metric, trace, and health examples

#### Scenario: Learner verifies observability packages
- **WHEN** a learner runs `go test ./stage-2-business/11-observability/...`
- **THEN** tests cover trace propagation, structured logging, Prometheus-style metrics, health checks, and HTTP observability routes
