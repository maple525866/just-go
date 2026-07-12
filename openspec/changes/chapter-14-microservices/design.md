## Context

Chapter 14 is currently a README-only placeholder. Learners have completed clean architecture and DDD examples and now need to cross a real process boundary using contract-first RPC, discovery, configuration, and a Gateway. The repository is one Go module and CI must remain deterministic, so the chapter must use real network transports without requiring Docker or external infrastructure.

The local environment does not currently provide `protoc`; generated Go files therefore need a documented, pinned generation path and must be committed so learners can run the chapter without generator tooling.

## Goals / Non-Goals

**Goals:**

- Demonstrate protobuf contracts plus unary, server-streaming, and bidirectional-streaming gRPC in one coherent product-and-inventory example.
- Demonstrate service registration, resolution, watching, dynamic configuration, deterministic rollout, Gateway routing, authentication, local rate limiting, and aggregation.
- Preserve real network boundaries while keeping default execution free of external services.
- Make concurrency ownership, cancellation, status mapping, startup, and shutdown explicit and testable.
- Provide focused tests, a runnable example, Chinese documentation, diagrams, and acceptance-driven exercises.

**Non-Goals:**

- Implement production Consul, etcd, Nacos, Kubernetes, TLS/mTLS, distributed rate limiting, or persistent configuration.
- Implement a message broker or the complete Capstone 3 blog services.
- Present the in-memory adapters as production-ready infrastructure.

## Decisions

### Use product and inventory services behind a Gateway

`ProductService.GetProduct` supplies product metadata. `InventoryService` supplies unary `GetStock`, server-streaming `WatchStock`, and bidirectional-streaming `SyncStock`. The Gateway concurrently calls the unary product and inventory operations to build one HTTP product-details response. This keeps every infrastructure concept connected to a visible request without prematurely implementing the blog capstone. A single greeting service was rejected because it cannot meaningfully demonstrate aggregation or independent service discovery.

### Treat protobuf as the contract source and commit generated Go code

Hand-written `.proto` files own field numbers and service definitions. Generated `.pb.go` and `_grpc.pb.go` files are committed and never edited manually. README generation commands pin `protoc-gen-go` and `protoc-gen-go-grpc` versions. Committing generated files lets learners run tests without installing generators; omitting the source was rejected because it would make protocol evolution opaque and irreproducible.

### Use real gRPC transport with self-contained infrastructure adapters

Services listen on random local TCP addresses in the runnable example, and clients connect through gRPC. Unit integration tests may use `bufconn`, which still exercises serialization and gRPC transport. Discovery stores service names, instance metadata, and addresses rather than Go service objects. This preserves the boundary being taught while avoiding external processes.

### Define narrow discovery and configuration contracts

Discovery exposes registration, deregistration, deterministic resolution, and snapshot watching. Configuration exposes versioned reads, validated updates, and snapshot subscriptions. Both in-memory adapters are concurrency-safe, context-aware, explicitly closable, and never block a producer indefinitely on a slow subscriber. Narrow contracts allow Consul or etcd adapters to be exercises rather than required dependencies.

### Make rollout deterministic

Gateway configuration contains a route switch, request timeout, local rate-limit values, and rollout percentage. A stable request key is hashed into a fixed bucket so the same key receives the same rollout decision while a configuration version is active. Random selection was rejected because it makes tests flaky and user behavior inconsistent.

### Keep Gateway policy at the transport boundary

The standard-library HTTP Gateway applies Bearer authentication, local rate limiting, dynamic route checks, discovery, gRPC connection management, concurrent downstream calls, and JSON/status mapping. Unknown internal errors are hidden. Required downstream failure fails the whole aggregate response. Partial responses were rejected because they would need an explicit partial-data contract that distracts from the chapter.

### Make lifecycle ownership explicit

The composition root owns listeners, registration handles, subscriptions, gRPC connections, HTTP and gRPC servers, and worker goroutines. Streaming loops and watchers observe context cancellation. Shutdown stops ingress, cancels work, closes connections, stops servers, and unregisters instances in an explicit order.

## Risks / Trade-offs

- **[Risk] In-memory discovery and configuration may look production-ready** → Document their guarantees and missing distributed-system behavior, and provide external-adapter exercises.
- **[Risk] Generated files enlarge the change** → Keep the protocol small, commit its source, pin generator versions, and verify regeneration instructions.
- **[Risk] Several components can obscure the core request flow** → Organize every component around one product-details request and keep packages focused.
- **[Risk] Streaming and subscriptions can leak goroutines** → Require context-aware loops, explicit close methods, cancellation tests, and race-detector verification.
- **[Risk] A single process cannot demonstrate deployment topology** → Use real listeners and gRPC connections, and state that deployment orchestration is out of scope.
- **[Risk] Local rate limiting differs from distributed enforcement** → Name and document it as a single-process teaching implementation.

## Migration Plan

1. Add protocol sources, pinned generation instructions, generated code, and runtime dependencies.
2. Implement product and inventory servers test-first, including all three RPC forms.
3. Implement discovery and dynamic configuration contracts and adapters test-first.
4. Implement Gateway policy, connection management, aggregation, and error mapping test-first.
5. Add the runnable composition root, integration tests, README, exercises, and roadmap updates.
6. Run chapter and repository quality gates, review the change, and fix all critical or important findings.

Rollback is isolated to the Chapter 14 directory, Go module dependency changes, roadmap entry, and this change's OpenSpec artifacts.

## Open Questions

None. The self-contained real-gRPC boundary and the product-and-inventory example are approved.
