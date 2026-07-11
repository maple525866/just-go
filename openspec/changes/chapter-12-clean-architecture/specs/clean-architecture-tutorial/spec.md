## ADDED Requirements

### Requirement: Chapter 12 provides an executable layered application
The curriculum SHALL provide a runnable Chapter 12 application under `stage-3-architecture/12-clean-architecture/` organized into `domain`, `usecase`, `interface`, and `infrastructure` layers.

#### Scenario: Learner runs the chapter application
- **WHEN** a learner runs `go run ./stage-3-architecture/12-clean-architecture`
- **THEN** the application starts successfully using dependencies composed at the outermost layer

#### Scenario: Learner verifies the complete chapter
- **WHEN** a learner runs `go test ./stage-3-architecture/12-clean-architecture/...`
- **THEN** all domain, use-case, adapter, HTTP, dependency-wiring, and architecture-boundary tests pass

### Requirement: Business rules remain independent of outer layers
The `domain` and `usecase` packages MUST NOT import packages from the chapter's `interface` or `infrastructure` layers, and repository behavior required by a use case SHALL be expressed as an inward-facing Go interface.

#### Scenario: Architecture boundary is checked
- **WHEN** the architecture test examines imports in the `domain` and `usecase` packages
- **THEN** it rejects any dependency on `interface` or `infrastructure` packages

#### Scenario: Infrastructure is substituted in a use-case test
- **WHEN** the article use case is constructed with a mock repository
- **THEN** its create and publish workflows are tested without a database, HTTP server, or infrastructure adapter

### Requirement: Domain invariants are explicit and tested
The domain layer SHALL own validation and state-transition rules for the example article entity and SHALL expose stable errors for invalid operations.

#### Scenario: Invalid article data is rejected
- **WHEN** an article is created with invalid title or body data
- **THEN** the domain returns a validation error without invoking an outer layer

#### Scenario: Article publication follows domain rules
- **WHEN** a valid draft article is published
- **THEN** its status and publication timestamp change according to domain rules

### Requirement: Adapters translate external details at the boundary
The chapter SHALL include an in-memory repository adapter and an HTTP adapter that translate their respective external representations to and from the use-case and domain types.

#### Scenario: Repository adapter persists an article
- **WHEN** the use case saves and retrieves an article through the repository port
- **THEN** the in-memory adapter returns an independent copy of the stored article

#### Scenario: HTTP adapter handles an article workflow
- **WHEN** a client creates, publishes, or retrieves an article through the chapter HTTP routes
- **THEN** the adapter returns an appropriate JSON response and maps known application errors to stable HTTP status codes

### Requirement: Compile-time dependency injection is demonstrated
The application SHALL include Wire provider declarations and checked-in generated wiring that composes concrete adapters with the use case and HTTP handler at compile time.

#### Scenario: Repository builds without the Wire CLI
- **WHEN** a learner runs the repository's normal build command without installing Wire
- **THEN** checked-in generated wiring compiles and constructs the Chapter 12 application

#### Scenario: Learner inspects provider relationships
- **WHEN** a learner reads the injector declaration
- **THEN** the file shows how the concrete repository satisfies the inward-facing repository port and how providers compose the application

### Requirement: Chapter 12 includes guided learning materials
The Chapter 12 README and exercises SHALL explain the dependency rule, package responsibilities, run and test commands, Wire generation, mock testing, and measurable extension tasks.

#### Scenario: Learner follows the chapter documentation
- **WHEN** a learner opens the Chapter 12 README and `EXERCISES.md`
- **THEN** they can identify the four layers, run the example, verify it, and complete exercises with explicit acceptance criteria
