## 1. Chapter Structure

- [x] 1.1 Add `main.go` that prints a cache and messaging learning report through the chapter packages.
- [x] 1.2 Create focused subpackages `cachex/` and `mqdemo/`.

## 2. Cache Store and Cache Patterns

- [x] 2.1 Implement a TTL key-value store with Get, Set, Delete, SetNX, and CompareAndDelete.
- [x] 2.2 Implement Cache-Aside with tests for miss, cache fill, and second-read hit.
- [x] 2.3 Implement Read-Through cache with tests for internal loader behavior.
- [x] 2.4 Implement Write-Through cache with tests for source and cache synchronization.

## 3. Cache Risk Countermeasures and Locks

- [x] 3.1 Implement negative cache behavior and tests that prevent repeated penetration.
- [x] 3.2 Implement deterministic TTL jitter calculation and tests for spread within bounds.
- [x] 3.3 Implement singleflight-style mutex loading and tests proving one loader call under concurrency.
- [x] 3.4 Implement Redis-style lock with token, TTL, holder-only release, and tests.

## 4. Message Queue Examples

- [x] 4.1 Implement an in-memory broker with Publish and Fetch.
- [x] 4.2 Implement Ack and RequeueExpired behavior.
- [x] 4.3 Add tests for publish/consume, ack removal, and unacked redelivery.

## 5. Learning Materials and Verification

- [x] 5.1 Update `stage-2-business/10-cache-and-mq/README.md` to replace placeholder content with package list, real Redis/NATS/Kafka notes, run commands, and knowledge-aligned checklist.
- [x] 5.2 Add `stage-2-business/10-cache-and-mq/EXERCISES.md` with 3 to 5 exercises, each including explicit acceptance criteria.
- [x] 5.3 Run `go test ./stage-2-business/10-cache-and-mq/...`, `go run ./stage-2-business/10-cache-and-mq`, `go test ./...`, and `go build ./...`; fix any failures.
