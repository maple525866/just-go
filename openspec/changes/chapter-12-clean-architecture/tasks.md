## 1. Chapter Structure and Domain

- [x] 1.1 Create the four-layer chapter package structure and a minimal composition-root application.
- [x] 1.2 Implement the article entity, validation errors, and publish state transition in `domain`.
- [x] 1.3 Add domain tests for creation, validation, and publication invariants.

## 2. Use Cases and Dependency Inversion

- [x] 2.1 Define the article repository port and article application service in `usecase`.
- [x] 2.2 Implement create, get, and publish workflows with stable application errors.
- [x] 2.3 Add mock-based use-case tests that do not import infrastructure or HTTP packages.

## 3. Outer-Layer Adapters

- [x] 3.1 Implement a concurrency-safe in-memory article repository with copy isolation.
- [x] 3.2 Add repository adapter tests for save, lookup, duplicate, and copy behavior.
- [x] 3.3 Implement JSON HTTP routes for create, get, and publish workflows with error mapping.
- [x] 3.4 Add `httptest` coverage for successful workflows, validation errors, missing articles, and malformed input.

## 4. Composition and Architecture Enforcement

- [x] 4.1 Add Wire provider declarations and checked-in generated dependency wiring.
- [x] 4.2 Add an architecture test that prevents `domain` and `usecase` from importing outer chapter layers.
- [x] 4.3 Add a runnable chapter entry point and smoke-test the composed application.

## 5. Learning Materials and Verification

- [x] 5.1 Replace the Chapter 12 placeholder README with architecture diagrams, package responsibilities, commands, and self-checks.
- [x] 5.2 Add exercises with explicit acceptance criteria and update the Chapter 12 roadmap output/progress entry.
- [x] 5.3 Run chapter tests, full tests, race tests, vet, build, and OpenSpec validation; fix failures.
