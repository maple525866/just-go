## 1. Chapter Structure and Dependencies

- [x] 1.1 Add GORM and SQLite dependencies with `go get gorm.io/gorm gorm.io/driver/sqlite github.com/mattn/go-sqlite3`.
- [x] 1.2 Add `main.go` that prints a persistence learning report through the chapter packages.
- [x] 1.3 Create focused subpackages `dbx/`, `sqlcrud/`, `gormdemo/`, `txdemo/`, and `migrations/`.

## 2. database/sql and Migration Examples

- [x] 2.1 Implement `dbx` helpers for opening SQLite memory databases, applying migration SQL, and reporting connection pool settings.
- [x] 2.2 Add migration SQL that creates an `articles` table.
- [x] 2.3 Implement `sqlcrud` Article CRUD with parameter binding and table-driven tests.
- [x] 2.4 Add a SQL injection regression test proving user input is bound as data.

## 3. GORM and Association Examples

- [x] 3.1 Implement GORM models for users, posts, and tags.
- [x] 3.2 Implement GORM AutoMigrate and seed helpers.
- [x] 3.3 Add tests for GORM CRUD and one-to-many `Preload("Posts")`.
- [x] 3.4 Add tests for many-to-many `Preload("Tags")`.

## 4. Transaction Examples

- [x] 4.1 Implement a successful transaction example that commits two article inserts.
- [x] 4.2 Implement a failing transaction example that rolls back inserted data.
- [x] 4.3 Add tests that verify commit and rollback outcomes with real database state.

## 5. Learning Materials and Verification

- [x] 5.1 Update `stage-2-business/09-data-persistence/README.md` to replace placeholder content with package list, SQLite/MySQL notes, run commands, and knowledge-aligned checklist.
- [x] 5.2 Add `stage-2-business/09-data-persistence/EXERCISES.md` with 3 to 5 exercises, each including explicit acceptance criteria.
- [x] 5.3 Run `go test ./stage-2-business/09-data-persistence/...`, `go run ./stage-2-business/09-data-persistence`, `go test ./...`, and `go build ./...`; fix any failures.
