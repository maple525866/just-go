## Why

Stage 3 needs an executable DDD tactical-patterns chapter that moves learners from architectural layering to rich domain modeling. Chapter 13 must replace its placeholder with a tested example before the microservices and distributed-resilience chapters can build on aggregates, repositories, and domain events.

## What Changes

- Replace the Chapter 13 placeholder with a runnable ordering domain that demonstrates entities, value objects, aggregate roots, repositories, domain services, application services, and domain events.
- Protect aggregate invariants through behavior-oriented methods and keep persistence concerns outside the domain package.
- Add an in-memory repository and synchronous event publisher to demonstrate cross-aggregate collaboration without external infrastructure.
- Add focused tests, guided exercises, diagrams, runnable examples, and updated curriculum progress.

## Capabilities

### New Capabilities

- `ddd-tactical-patterns-tutorial`: Defines the executable ordering example, tactical DDD building blocks, event-driven collaboration, tests, and learning materials for Chapter 13.

### Modified Capabilities

- `learning-curriculum`: Marks Chapter 13 as implemented and requires its README and roadmap entry to describe concrete outputs instead of placeholder content.

## Impact

- Adds Go code and tests under `stage-3-architecture/13-ddd-patterns/`.
- Updates the Chapter 13 README, adds exercises, and updates `ROADMAP.md` output/progress text.
- Adds no production runtime dependencies or external services; examples use standard-library types and in-memory adapters.
- Adds OpenSpec tracking artifacts for the chapter.
