## 1. Contracts and Shared Infrastructure

- [x] 1.1 Add the blog protobuf contract and generated Go/gRPC code for user, post, and comment services.
- [x] 1.2 Add shared domain errors, W3C trace propagation, structured span export, and gRPC interceptors with tests.
- [x] 1.3 Add reusable limiter, retry, and circuit-breaker support for Gateway aggregation with tests.

## 2. Bounded Context Services

- [x] 2.1 Implement user-svc registration, login, token validation, lookup, independent storage, and tests.
- [x] 2.2 Implement post-svc article CRUD/list, ownership checks, independent storage, and tests.
- [x] 2.3 Implement comment-svc nested comments, soft deletion, independent storage, and tests.

## 3. API Gateway and End-to-End Flow

- [x] 3.1 Implement Gateway HTTP routes, authentication, error mapping, and backend clients.
- [x] 3.2 Implement article detail aggregation with timeout, retry, breaker, limiter, and explicit comment fallback.
- [x] 3.3 Add in-process real gRPC/HTTP end-to-end tests covering the complete blog flow and trace propagation.

## 4. Runnable Services and Delivery

- [x] 4.1 Add service and Gateway command entrypoints with environment-based addresses and graceful shutdown.
- [x] 4.2 Add Dockerfile, docker-compose, health endpoints, and optional Jaeger wiring.
- [x] 4.3 Replace the placeholder README, add measurable exercises, and document architecture/CAP/failure trade-offs.
- [x] 4.4 Update `ROADMAP.md` and Capstone 3 completion checklist.

## 5. Verification

- [x] 5.1 Run gofmt, Capstone 3 tests, full tests, race tests, vet, build, golangci-lint, and strict OpenSpec validation where available.
- [x] 5.2 Complete independent subagent code review and resolve all blocking findings.
