## 1. Value Objects and Domain Primitives

- [x] 1.1 Create the Chapter 13 package structure and shared domain errors and event primitives.
- [x] 1.2 Implement immutable validated `Money` and `Address` value objects.
- [x] 1.3 Add value-object tests for validation, equality, arithmetic, immutability, and currency mismatch.

## 2. Entities and Aggregate Root

- [x] 2.1 Implement order-line entities with stable identity and defensive value access.
- [x] 2.2 Implement the versioned `Order` aggregate root with line-management and confirmation invariants.
- [x] 2.3 Record immutable `OrderConfirmed` events and expose pull-and-clear event behavior.
- [x] 2.4 Add aggregate tests for creation, line behavior, lifecycle rules, totals, versions, and event recording.

## 3. Domain and Application Services

- [x] 3.1 Implement the pricing domain service and discount policy contract with focused tests.
- [x] 3.2 Define repository and event-publisher ports plus stable application errors.
- [x] 3.3 Implement create, add-line, get, and confirm application workflows with correct save-before-publish ordering.
- [x] 3.4 Add mock-based application tests for orchestration, conflicts, failures, and event publication.

## 4. Infrastructure Adapters and Composition

- [x] 4.1 Implement a concurrency-safe in-memory repository with deep-copy isolation and optimistic concurrency.
- [x] 4.2 Add repository tests for creation, lookup, isolation, updates, and stale-version conflicts.
- [x] 4.3 Implement a synchronous event bus and inventory projection handler.
- [x] 4.4 Add event bus tests for routing, error propagation, unknown events, and projection updates.
- [x] 4.5 Add a runnable composition-root example and an integration test for the complete confirmation flow.

## 5. Learning Materials and Verification

- [x] 5.1 Replace the Chapter 13 placeholder README with diagrams, concept mapping, workflow, commands, limitations, and self-checks.
- [x] 5.2 Add exercises with explicit acceptance criteria and update the Chapter 13 roadmap output/progress entry.
- [x] 5.3 Run chapter tests, full tests, race tests, vet, build, lint, and OpenSpec validation; fix failures.
