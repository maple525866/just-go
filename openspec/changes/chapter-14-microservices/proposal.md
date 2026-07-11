## Why

Stage 3 needs an executable microservices chapter that moves learners from in-process domain boundaries to contract-first network communication. Chapter 14 must replace its placeholder with a tested, self-contained example before the resilience and microservice capstone chapters can build on gRPC, service discovery, API Gateway, and dynamic configuration.

## What Changes

- Replace the Chapter 14 placeholder with runnable product and inventory gRPC services using generated protobuf contracts.
- Demonstrate unary, server-streaming, and bidirectional-streaming RPC with stable validation and status semantics.
- Add replaceable service-discovery and dynamic-configuration contracts with deterministic, concurrency-safe in-memory adapters.
- Add an HTTP API Gateway that demonstrates Bearer authentication, local rate limiting, dynamic route control, service resolution, concurrent downstream aggregation, and stable error mapping.
- Add focused tests, a complete runnable composition example, Chinese learning materials, diagrams, measurable exercises, and updated curriculum progress.
- Add the gRPC and protobuf runtime dependencies required by the chapter while keeping default execution free of external services.

## Capabilities

### New Capabilities

- `microservices-foundations-tutorial`: Defines the executable product-and-inventory example, protobuf/gRPC communication, service discovery, dynamic configuration, API Gateway behavior, tests, and learning materials for Chapter 14.

### Modified Capabilities

- `learning-curriculum`: Marks Chapter 14 as implemented and requires its README and roadmap entry to describe concrete outputs instead of placeholder content.

## Impact

- Adds Go code, generated protobuf code, tests, and exercises under `stage-3-architecture/14-microservices/`.
- Updates the Chapter 14 README and `ROADMAP.md` output/progress text.
- Adds `google.golang.org/grpc` and `google.golang.org/protobuf` runtime dependencies plus reproducible protobuf generation instructions.
- Adds no required Docker, Consul, etcd, Nacos, Kubernetes, TLS, or message-broker runtime.
- Adds OpenSpec tracking artifacts for the chapter.
