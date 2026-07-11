## ADDED Requirements

### Requirement: Chapter 13 provides an executable DDD tactical-patterns example
The curriculum SHALL provide a runnable Chapter 13 application under `stage-3-architecture/13-ddd-patterns/` that demonstrates entities, value objects, an aggregate root, a repository, a domain service, an application service, and domain events in one coherent ordering example.

#### Scenario: Learner runs the chapter application
- **WHEN** a learner runs `go run ./stage-3-architecture/13-ddd-patterns`
- **THEN** the application creates and confirms an order and displays the resulting domain-event collaboration

#### Scenario: Learner verifies the complete chapter
- **WHEN** a learner runs `go test ./stage-3-architecture/13-ddd-patterns/...`
- **THEN** all domain, application, repository, event-publisher, and integration tests pass

### Requirement: Value objects are valid and immutable
The domain SHALL provide `Money` and `Address` value objects that validate construction, expose behavior without mutable fields, compare by value, and reject invalid arithmetic.

#### Scenario: Invalid value object is rejected
- **WHEN** a caller constructs money with an unsupported currency or an address missing required components
- **THEN** construction returns a stable domain validation error

#### Scenario: Money arithmetic preserves currency invariants
- **WHEN** money values with the same currency are added or multiplied
- **THEN** a new value is returned without modifying either operand

- **WHEN** money values with different currencies are added
- **THEN** the operation returns a currency-mismatch error

### Requirement: Order aggregate protects business invariants
The `Order` aggregate root SHALL own its order-line entities and SHALL permit mutation only through behavior methods that enforce quantity, uniqueness, lifecycle, and confirmation invariants.

#### Scenario: Draft order manages lines through the aggregate root
- **WHEN** a valid product and quantity are added to a draft order
- **THEN** the aggregate updates its lines and total while callers cannot mutate internal line state directly

#### Scenario: Confirmed order rejects further mutation
- **WHEN** a caller confirms a non-empty order and then attempts to add or remove a line
- **THEN** confirmation succeeds once and subsequent mutation returns a stable invalid-state error

### Requirement: Domain and application services have distinct responsibilities
The chapter SHALL provide a stateless pricing domain service for domain calculation and an application service that orchestrates repository access, aggregate behavior, persistence, and event publication.

#### Scenario: Domain service calculates an order price
- **WHEN** the pricing service receives order lines and an applicable discount policy
- **THEN** it calculates a valid total without loading data, saving aggregates, or publishing events

#### Scenario: Application service confirms an order
- **WHEN** the application service is asked to confirm an existing valid order
- **THEN** it loads the aggregate, invokes its behavior, saves it, and publishes the resulting events in that order

### Requirement: Repository isolates aggregate persistence
The application layer SHALL express an order repository contract, and the chapter SHALL provide a concurrency-safe in-memory adapter that preserves aggregate isolation through deep copies.

#### Scenario: Stored aggregate is isolated from callers
- **WHEN** an order is saved and a caller later mutates either the original or a loaded instance
- **THEN** the repository's stored aggregate changes only after an explicit successful save

#### Scenario: Concurrent saves detect stale versions
- **WHEN** two callers save changes derived from the same aggregate version
- **THEN** the first valid save succeeds and the stale save returns a conflict error

### Requirement: Domain events enable cross-aggregate collaboration
The order aggregate SHALL record an immutable `OrderConfirmed` domain event, and the application SHALL publish collected events after persistence through a replaceable event-publisher contract.

#### Scenario: Confirmation publishes an event after saving
- **WHEN** an order is confirmed successfully
- **THEN** exactly one `OrderConfirmed` event containing order identity and line data is published after the confirmed aggregate is stored

#### Scenario: Event handler updates a separate projection
- **WHEN** the synchronous event bus dispatches `OrderConfirmed`
- **THEN** an inventory-oriented handler records the confirmed quantities without the order aggregate directly depending on that handler

### Requirement: Chapter 13 includes guided learning materials
The Chapter 13 README and exercises SHALL explain tactical-pattern boundaries, aggregate invariants, repository isolation, service responsibilities, event flow, run and test commands, limitations, and measurable extension tasks.

#### Scenario: Learner follows the chapter documentation
- **WHEN** a learner opens the Chapter 13 README and `EXERCISES.md`
- **THEN** they can trace the ordering workflow, run the example, verify it, and complete exercises with explicit acceptance criteria
