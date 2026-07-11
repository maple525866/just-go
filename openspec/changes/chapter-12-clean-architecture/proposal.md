## Why

Stage 3 needs an executable clean-architecture chapter that teaches learners how to keep business rules independent from delivery and persistence details. Chapter 12 is the foundation for the later DDD and microservices chapters, so it must turn the existing placeholder into a tested, layered Go example before those chapters are implemented.

## What Changes

- Replace the Chapter 12 placeholder with a runnable four-layer example organized around `domain`, `usecase`, `interface`, and `infrastructure` packages.
- Demonstrate dependency inversion through an inward-facing use-case port and infrastructure adapters without allowing the business core to import external details.
- Add compile-time dependency injection with Wire, while keeping generated wiring buildable without requiring learners to install the Wire CLI.
- Add mock-based use-case tests, architecture boundary tests, focused exercises, and chapter documentation.

## Capabilities

### New Capabilities

- `clean-architecture-tutorial`: Defines the runnable layered example, dependency rules, dependency-injection wiring, tests, and learning materials for Chapter 12.

### Modified Capabilities

- `learning-curriculum`: Marks Chapter 12 as implemented and requires its README and roadmap entry to describe concrete outputs instead of placeholder content.

## Impact

- Adds Go code and tests under `stage-3-architecture/12-clean-architecture/`.
- Updates the Chapter 12 README and `ROADMAP.md` progress/output text.
- Adds Wire as a compile-time-only wiring tool dependency if required by the chosen implementation.
- Adds OpenSpec tracking artifacts for the chapter.
