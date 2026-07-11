## Context

Chapter 12 is currently a README-only placeholder. It follows the blog API capstone and precedes DDD and microservices, so the chapter needs a small but realistic domain that demonstrates architectural boundaries without duplicating the whole capstone. The repository is a single Go module and all examples must remain runnable with the standard repository-wide test and build commands.

## Goals / Non-Goals

**Goals:**

- Make dependency direction visible in the filesystem and import graph.
- Use a small article-publishing domain to demonstrate entities, ports, use cases, HTTP adapters, persistence adapters, and composition.
- Show both Wire provider declarations and generated compile-time wiring.
- Make the use-case layer testable with hand-written mocks and no database or HTTP server.
- Enforce important layer boundaries with an automated architecture test.

**Non-Goals:**

- Rebuild the complete Capstone 2 blog API.
- Add a production database, authentication, migrations, or distributed-system concerns.
- Treat the example package structure as the only valid clean-architecture layout.

## Decisions

### Use an article-publishing slice

The example will support creating and publishing an article. This is familiar from Capstone 2, has meaningful domain invariants, and is small enough that learners can trace every dependency. A synthetic greeting or calculator example was rejected because it would not demonstrate a repository port or a useful application workflow.

### Put stable contracts toward the core

Entities and domain errors live in `domain`. The repository port used by the application lives in `usecase`, beside the workflow that consumes it. `infrastructure/memory` implements that port, while `interface/httpapi` translates HTTP requests and responses. Neither `domain` nor `usecase` imports an outer layer.

The directory name `interface` mirrors the roadmap terminology; concrete Go packages below it use legal descriptive names such as `httpapi`, since `interface` is a Go keyword.

### Keep composition at the outer edge

`wire.go` declares providers behind the conventional `wireinject` build tag. A checked-in `wire_gen.go` contains the generated constructor path so normal `go build` and `go test` do not require the Wire CLI. `main.go` only starts the composed application. Manual constructor calls remain visible in tests so learners understand what Wire generates.

Repository updates use compare-and-swap semantics: the use case supplies the status it read, and the adapter rejects the write if another request has already changed that state. This keeps concurrent publication from reporting two successes without moving the domain transition into the infrastructure layer.

### Use hand-written mocks and an import-boundary test

Use-case tests use a small hand-written repository mock to keep the teaching surface explicit. An architecture test parses imports under the chapter directory and rejects imports from `domain` or `usecase` to outer-layer package paths. This demonstrates that clean architecture is a dependency rule, not merely folder naming.

## Risks / Trade-offs

- **[Risk] Wire is archived and adds a tool dependency** → Pin the dependency, keep generated wiring checked in, and explain that the underlying constructor-injection pattern is tool-independent.
- **[Risk] A small in-memory example may look production-ready** → Clearly label concurrency, persistence, and error-mapping limitations and provide exercises that extend the ports/adapters.
- **[Risk] Architecture tests can overfit directory names** → Check only the two core layers and document the intended dependency rule.
- **[Risk] The `interface` directory resembles the Go keyword** → Use it only as a grouping directory and use `httpapi` as the actual package name.

## Migration Plan

1. Replace the Chapter 12 placeholder with the layered example and tests.
2. Add and generate the Wire injector, then commit generated output.
3. Update chapter and roadmap documentation after verification.
4. Rollback is isolated to the Chapter 12 directory, roadmap entry, module dependency, and this change's artifacts.

## Open Questions

None. The implementation can proceed with the scoped article-publishing slice.
