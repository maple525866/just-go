## ADDED Requirements

### Requirement: Chapter 14 provides an executable microservices foundations example

The curriculum SHALL provide a runnable Chapter 14 application under `stage-3-architecture/14-microservices/` that composes product and inventory gRPC services, service discovery, dynamic configuration, and an HTTP API Gateway without requiring an external service.

#### Scenario: Learner runs the chapter application

- **WHEN** a learner runs `go run ./stage-3-architecture/14-microservices`
- **THEN** the application MUST start the composed services, execute a product-details request through the Gateway, display the aggregate result, and exit cleanly

#### Scenario: Learner verifies the complete chapter

- **WHEN** a learner runs `go test ./stage-3-architecture/14-microservices/...`
- **THEN** all protocol-service, discovery, configuration, Gateway, streaming, lifecycle, and integration tests MUST pass without Docker or an external service

### Requirement: Protobuf contracts define real gRPC communication

The chapter SHALL provide hand-written protobuf definitions and generated Go/gRPC code for product and inventory services, and clients MUST communicate with servers through gRPC transport rather than direct service-object calls.

#### Scenario: Unary RPC returns product and stock data

- **WHEN** a client calls `ProductService.GetProduct` or `InventoryService.GetStock` with a known SKU
- **THEN** the corresponding service MUST return its data through a unary gRPC response

#### Scenario: Server stream reports inventory changes

- **WHEN** a client subscribes to `InventoryService.WatchStock` for a known SKU and stock changes
- **THEN** the server MUST stream ordered stock snapshots until the client cancels or the stream ends

#### Scenario: Bidirectional stream synchronizes adjustments

- **WHEN** a client sends multiple valid stock adjustments through `InventoryService.SyncStock`
- **THEN** the server MUST apply each adjustment in order and stream one updated stock result for each accepted request

#### Scenario: RPC validation uses stable status codes

- **WHEN** a request is malformed or names an unknown product or SKU
- **THEN** the service MUST return `InvalidArgument` for invalid input or `NotFound` for a missing resource without exposing internal errors

### Requirement: Service discovery resolves network instances

The chapter SHALL define a replaceable service-discovery contract and a concurrency-safe in-memory adapter that registers, deregisters, deterministically resolves, and watches immutable snapshots of service instances containing network addresses.

#### Scenario: Registered service can be resolved

- **WHEN** an instance with a valid service name, instance ID, and address is registered
- **THEN** a resolver MUST return that instance as an available network endpoint

#### Scenario: Deregistration updates resolvers and watchers

- **WHEN** a registration is explicitly removed
- **THEN** subsequent resolution and watcher snapshots MUST exclude that instance

#### Scenario: Resolution reports no available instance

- **WHEN** a caller resolves a service with no registered instance
- **THEN** discovery MUST return a stable unavailable error

#### Scenario: Watch stops on cancellation or close

- **WHEN** a watcher context is canceled or the registry is closed
- **THEN** the watcher MUST terminate without leaking a goroutine and new operations after close MUST fail

### Requirement: Dynamic configuration is versioned and observable

The chapter SHALL define a replaceable configuration contract and a concurrency-safe in-memory adapter that validates Gateway settings, assigns monotonically increasing versions, publishes immutable snapshots, and makes rollout decisions deterministic for a stable request key.

#### Scenario: Valid update advances the version

- **WHEN** a valid Gateway configuration is stored
- **THEN** its version MUST be greater than the previous version and subscribers MUST observe the new immutable snapshot

#### Scenario: Invalid update preserves current configuration

- **WHEN** an update contains an invalid timeout, rate limit, or rollout percentage
- **THEN** the update MUST fail and the previously active version MUST remain unchanged

#### Scenario: Rollout decision is stable and bounded

- **WHEN** the same request key is evaluated repeatedly against one rollout percentage
- **THEN** every evaluation MUST return the same decision, with zero percent always disabled and one hundred percent always enabled

#### Scenario: Subscription stops on cancellation or close

- **WHEN** a subscriber context is canceled or the configuration store is closed
- **THEN** the subscription MUST terminate without leaking a goroutine and new operations after close MUST fail

### Requirement: API Gateway applies edge policy and aggregates services

The chapter SHALL provide an HTTP API Gateway that applies Bearer authentication, single-process rate limiting, dynamic route and rollout policy, discovery-based gRPC calls, and product-plus-stock aggregation with a configured deadline.

#### Scenario: Authorized request returns an aggregate response

- **WHEN** an authorized, rate-permitted request targets an enabled product-details route and both services succeed
- **THEN** the Gateway MUST return HTTP 200 with one JSON response containing product and stock data

#### Scenario: Authentication and rate limits reject requests

- **WHEN** a request lacks the configured Bearer token or exceeds the active local rate limit
- **THEN** the Gateway MUST return HTTP 401 or HTTP 429 respectively without calling downstream services

#### Scenario: Dynamic policy disables a route

- **WHEN** the route is disabled or the stable request key is outside the active rollout percentage
- **THEN** the Gateway MUST reject the request without calling downstream services

#### Scenario: Downstream status maps to stable HTTP status

- **WHEN** a required gRPC call returns `InvalidArgument`, `NotFound`, `Unavailable`, or `DeadlineExceeded`
- **THEN** the Gateway MUST return HTTP 400, 404, 503, or 504 respectively and MUST NOT expose internal error text

#### Scenario: Required downstream failure rejects partial data

- **WHEN** either the product or inventory result fails
- **THEN** the Gateway MUST return an error response and MUST NOT return a partial aggregate

### Requirement: Component lifecycle is explicit and cancelable

The chapter SHALL make ownership and shutdown of listeners, registrations, subscriptions, gRPC connections, HTTP and gRPC servers, and worker goroutines explicit.

#### Scenario: Composition shuts down cleanly

- **WHEN** the runnable example completes or its root context is canceled
- **THEN** ingress MUST stop, outstanding work MUST be canceled, owned connections and servers MUST close, and service registrations MUST be removed

#### Scenario: Streaming client cancellation stops server work

- **WHEN** a streaming client cancels its context
- **THEN** the corresponding server loop MUST stop without continuing to publish or leaking a goroutine

### Requirement: Chapter 14 includes guided learning materials

The Chapter 14 README and exercises SHALL explain protobuf compatibility, RPC forms, discovery, dynamic configuration, Gateway responsibilities, synchronous gRPC versus asynchronous MQ, run and generation commands, limitations, and measurable extension tasks.

#### Scenario: Learner follows the chapter documentation

- **WHEN** a learner opens the Chapter 14 README and `EXERCISES.md`
- **THEN** they MUST be able to trace the aggregate request, run and test the example, regenerate protocol code with pinned tools, and complete exercises with explicit acceptance criteria

#### Scenario: Production limitations are explicit

- **WHEN** a learner reads the infrastructure sections
- **THEN** the materials MUST state that the in-memory discovery, configuration, and local rate limiter are teaching adapters rather than production distributed implementations
