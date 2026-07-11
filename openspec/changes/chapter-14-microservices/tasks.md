## 1. Protocol Contracts and Dependencies

- [x] 1.1 Define product and inventory protobuf contracts with unary, server-streaming, and bidirectional-streaming RPCs.
- [x] 1.2 Add pinned protobuf generation configuration or commands and generate the Go/gRPC sources from the protocol files.
- [x] 1.3 Add the minimal gRPC and protobuf runtime dependencies, run `go mod tidy`, and verify generated packages build without local generator tools.

## 2. Product and Inventory gRPC Services

- [x] 2.1 Write failing tests and implement a concurrency-safe product catalog with validation and immutable reads.
- [x] 2.2 Write failing tests and implement `ProductService.GetProduct` with `InvalidArgument` and `NotFound` status mapping.
- [x] 2.3 Write failing tests and implement inventory reads and adjustments with validation, ordered updates, and concurrent access safety.
- [ ] 2.4 Write failing tests and implement `InventoryService.GetStock` and `WatchStock`, including cancellation and ordered snapshot delivery.
- [ ] 2.5 Write failing tests and implement `InventoryService.SyncStock`, including per-request responses, validation failure, EOF, and cancellation.
- [ ] 2.6 Add a gRPC transport integration test proving unary and streaming clients communicate through serialization rather than direct Go calls.

## 3. Service Discovery and Dynamic Configuration

- [ ] 3.1 Define service instance, registrar, resolver, watcher, unavailable, validation, and closed semantics in a focused discovery contract.
- [ ] 3.2 Write failing tests and implement concurrency-safe in-memory registration, deterministic resolution, immutable snapshots, deregistration, cancellation, and close behavior.
- [ ] 3.3 Define validated Gateway configuration, versioned snapshot, subscription, and stable rollout-selection contracts.
- [ ] 3.4 Write failing tests and implement concurrency-safe configuration reads, monotonic updates, non-blocking immutable subscriptions, cancellation, close behavior, and deterministic 0–100 percent rollout.

## 4. HTTP API Gateway

- [ ] 4.1 Write failing tests and implement Bearer authentication plus a documented single-process rate limiter that rejects before downstream calls.
- [ ] 4.2 Write failing tests and implement dynamic route and rollout middleware using the active configuration snapshot and stable request keys.
- [ ] 4.3 Write failing tests and implement discovery-backed gRPC connection ownership for product and inventory clients with explicit close behavior.
- [ ] 4.4 Write failing tests and implement concurrent product-plus-stock aggregation under the configured deadline, rejecting partial results.
- [ ] 4.5 Write failing tests and implement stable JSON responses plus gRPC-to-HTTP mappings for invalid, missing, unavailable, deadline, and unknown errors without leaking internal text.
- [ ] 4.6 Add HTTP integration tests for success, authentication, rate limiting, dynamic route changes, rollout decisions, discovery failure, downstream failure, and timeout.

## 5. Composition and Learning Materials

- [ ] 5.1 Write a failing lifecycle test and implement the composition root with random local listeners, registrations, Gateway startup, one aggregate request, and ordered shutdown.
- [ ] 5.2 Replace the Chapter 14 placeholder README with architecture and sequence diagrams, protocol compatibility rules, RPC-form guidance, discovery/configuration/Gateway boundaries, gRPC-versus-MQ guidance, commands, and explicit production limitations.
- [ ] 5.3 Add `EXERCISES.md` with measurable acceptance criteria for Consul/etcd adapters, health checking, client streaming, TLS, persistent configuration, distributed limiting, and MQ decoupling.
- [ ] 5.4 Update the Chapter 14 roadmap output and progress checkbox while leaving Chapter 15 and Capstone 3 incomplete.

## 6. Verification and Review

- [ ] 6.1 Run `gofmt`, Chapter 14 tests, full repository tests, `go vet`, `go build`, race tests, `golangci-lint`, and OpenSpec validation; fix every in-scope failure and record unavailable checks accurately.
- [ ] 6.2 Review the implementation against every OpenSpec scenario and the approved design, and confirm generated files match their protocol sources.
- [ ] 6.3 Dispatch the requested code-review subagent against the full `origin/main...HEAD` change and resolve every Critical or Important finding.
- [ ] 6.4 Create the requested GitHub issue, link it from the change, push `codex/chapter-14-microservices`, and create a Pull Request against `main` with verification evidence.
