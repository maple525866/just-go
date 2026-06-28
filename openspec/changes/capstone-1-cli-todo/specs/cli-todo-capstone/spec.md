## ADDED Requirements

### Requirement: Capstone 1 SHALL provide a runnable CLI Todo program

`stage-1-syntax/capstone-1-cli-todo/` SHALL contain a runnable `package main` CLI program for managing local todo tasks. The program MUST support `add`, `list`, `done`, `delete`, `clear`, and `help` behavior.

#### Scenario: help output is available
- **WHEN** a learner executes `go run ./stage-1-syntax/capstone-1-cli-todo --help`
- **THEN** the command MUST exit successfully and print usage text listing the supported subcommands

#### Scenario: CLI entry delegates to packages
- **WHEN** a reader inspects `main.go`
- **THEN** it MUST delegate command handling to an internal package instead of implementing all todo logic directly in `main.go`

### Requirement: Capstone 1 SHALL model todos with structs and slices

The project SHALL define task and list types using Go structs and slices. Each task MUST include an id, title, done state, created time, and updated time.

#### Scenario: add creates a task
- **WHEN** a learner adds a todo title to an empty list
- **THEN** the list MUST contain one task with a stable id, the given title, `Done == false`, and non-zero timestamps

#### Scenario: done marks a task complete
- **WHEN** a learner marks an existing task id as done
- **THEN** that task MUST have `Done == true` and an updated timestamp

#### Scenario: delete removes a task
- **WHEN** a learner deletes an existing task id
- **THEN** that task MUST no longer appear in the task list

### Requirement: Capstone 1 SHALL persist todos as JSON

The project SHALL persist todo data in a local JSON file using `os`, `encoding/json`, and `time`. The storage path MUST be configurable for tests and CLI use.

#### Scenario: JSON store round trip
- **WHEN** a learner saves a todo list to a JSON file and loads it again
- **THEN** the loaded list MUST preserve task ids, titles, done states, and timestamps

#### Scenario: missing storage file loads empty list
- **WHEN** the configured JSON file does not exist
- **THEN** loading MUST return an empty todo list rather than failing with an unhandled error

### Requirement: Capstone 1 SHALL use custom errors

The project SHALL define custom or sentinel errors for expected domain and command failures, including invalid commands, missing titles, invalid ids, and task not found.

#### Scenario: unknown command returns command error
- **WHEN** a learner invokes an unsupported subcommand
- **THEN** the app MUST return an error that callers can classify as invalid command

#### Scenario: missing task returns not found error
- **WHEN** a learner marks or deletes an id that does not exist
- **THEN** the app MUST return an error that callers can classify as task not found

### Requirement: Capstone 1 SHALL demonstrate controlled asynchronous persistence

The project SHALL include a goroutine + channel based asynchronous save worker that accepts save requests and can be closed without leaking goroutines.

#### Scenario: async save writes data
- **WHEN** a learner submits a todo list to the async save worker and closes it
- **THEN** the configured JSON file MUST contain the submitted data

#### Scenario: async save exits cleanly
- **WHEN** the async save worker is closed
- **THEN** it MUST stop accepting new work and return without indefinite blocking

### Requirement: Capstone 1 SHALL provide tests and benchmark

The project SHALL include table-driven tests for domain logic, storage, CLI command behavior, and async save behavior; it MUST also include at least one benchmark for a meaningful operation such as rendering or list mutation.

#### Scenario: go test passes
- **WHEN** a learner runs `go test ./stage-1-syntax/capstone-1-cli-todo/...`
- **THEN** all tests MUST pass with exit code 0

#### Scenario: benchmark runs
- **WHEN** a learner runs `go test -bench=. ./stage-1-syntax/capstone-1-cli-todo/...`
- **THEN** at least one benchmark MUST execute and the command MUST exit successfully

### Requirement: Capstone 1 SHALL provide README and exercises

The capstone README MUST replace placeholder content with implemented functionality, commands, file layout, and completion checklist; the project MUST include `EXERCISES.md` with 3 to 5 extension exercises and explicit acceptance criteria.

#### Scenario: README describes runnable CLI
- **WHEN** a reader opens `stage-1-syntax/capstone-1-cli-todo/README.md`
- **THEN** it MUST document supported subcommands, example invocations, data file behavior, and verification commands

#### Scenario: exercises include acceptance criteria
- **WHEN** a reader opens `stage-1-syntax/capstone-1-cli-todo/EXERCISES.md`
- **THEN** it MUST include 3 to 5 exercises and each exercise MUST include concrete acceptance criteria
