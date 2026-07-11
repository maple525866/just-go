## Context

Chapter 13 is currently a README-only placeholder. Learners have just completed a layered article application in Chapter 12 and now need to learn how behavior and invariants are modeled inside the domain core. The repository is a single Go module, so the example must stay independently understandable and pass repository-wide test, race, vet, and build commands without external infrastructure.

## Goals / Non-Goals

**Goals:**

- Demonstrate entity identity, immutable value objects, aggregate ownership, repository boundaries, domain services, application services, and domain events in idiomatic Go.
- Use an order aggregate with realistic invariants that can only be changed through aggregate behavior.
- Demonstrate a cross-aggregate reaction through a synchronous event publisher while keeping dispatch outside the aggregate.
- Make concepts executable through focused tests, a runnable example, diagrams, and acceptance-driven exercises.

**Non-Goals:**

- Build an HTTP API, production database, transactional outbox, event broker, or event-sourced system.
- Provide a generic DDD framework or base entity hierarchy.
- Claim that every application needs aggregates, repositories, or domain events.

## Decisions

### Model an ordering bounded context

The core example uses an `Order` aggregate containing order-line entities. `Money` and `Address` are immutable value objects, and the order owns all line mutation and confirmation rules. Ordering gives each tactical pattern a concrete purpose without repeating Chapter 12's article example. A library or bank-account example was rejected because it offers a weaker demonstration of both a domain service and cross-aggregate collaboration.

### Keep the domain model free of persistence and dispatch

The `domain/order` package exposes constructors and behavior methods rather than public mutable fields. It records immutable `DomainEvent` values internally and exposes a pull-and-clear operation. Repository and publisher contracts live in the application package, where orchestration needs them. The aggregate never saves itself or calls an event bus.

### Use explicit value objects and integer minor units

`Money` stores an ISO-style currency code and integer minor units to avoid floating-point errors. Arithmetic rejects currency mismatches. `Address` validates and normalizes required components at construction. Plain structs with exported fields were rejected because they would let callers create invalid values and weaken the teaching point.

### Separate domain calculation from application orchestration

A stateless pricing domain service calculates totals from order lines and a discount policy because pricing spans multiple domain values without belonging to one entity's lifecycle. The application service loads and saves aggregates, invokes domain behavior, publishes collected events only after persistence succeeds, and maps infrastructure errors into stable application errors.

### Use deterministic in-memory adapters

A concurrency-safe in-memory repository returns deep copies to preserve aggregate boundaries. A synchronous event bus registers handlers by event name and stops on handler failure. An inventory projection reacts to `OrderConfirmed` to demonstrate cross-aggregate collaboration. These adapters keep tests deterministic and make the dispatch flow visible; production durability semantics are documented as out of scope.

## Risks / Trade-offs

- **[Risk] The example may imply one canonical Go DDD package layout** → Explain that boundaries and invariant ownership matter more than directory names.
- **[Risk] Pulling events from an aggregate before persistence could lose events on failure** → The application service only pulls events after a successful save and tests the failure path.
- **[Risk] Synchronous dispatch can partially execute handlers** → Document the limitation and use exercises to introduce idempotency and outbox design.
- **[Risk] Defensive copying adds code that distracts from the patterns** → Keep copy logic small and test it explicitly as part of aggregate isolation.

## Migration Plan

1. Add the domain types and invariant tests.
2. Add application contracts, pricing service, workflows, and mock-based tests.
3. Add in-memory repository, event bus, projection, and adapter tests.
4. Add the runnable composition example and learning materials.
5. Update roadmap progress and run all repository quality gates.

Rollback is isolated to the Chapter 13 directory, roadmap entry, and this change's OpenSpec artifacts.

## Open Questions

None. The ordering slice and in-memory delivery semantics are intentionally scoped for a tactical-patterns tutorial.
