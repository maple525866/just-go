## ADDED Requirements

### Requirement: Stage 2 capstone blog API
The curriculum SHALL provide an executable Capstone 2 monolith under `stage-2-business/capstone-2-blog-api/` that integrates Web routing, persistence, cache, messaging-oriented invalidation, and observability concepts.

#### Scenario: Learner runs the capstone demo
- **WHEN** a learner runs `go run ./stage-2-business/capstone-2-blog-api`
- **THEN** the program prints a concise service report including available routes and observability endpoints

#### Scenario: Learner verifies capstone behavior
- **WHEN** a learner runs `go test ./stage-2-business/capstone-2-blog-api/...`
- **THEN** tests cover auth, article CRUD, pagination, tags, comments, cache invalidation, HTTP smoke flow, health, and metrics
