## 1. Project Structure and Auth

- [x] 1.1 Add focused packages `auth/`, `model/`, `store/`, `cache/`, `observability/`, and `server/` plus `main.go`.
- [x] 1.2 Implement password hashing and signed bearer tokens with tests.
- [x] 1.3 Implement user registration/login flow with validation and duplicate user handling.

## 2. Articles, Tags, Comments, and Cache

- [x] 2.1 Implement article CRUD, tag filtering, and pagination in the store.
- [x] 2.2 Implement nested comments and soft deletion.
- [x] 2.3 Implement TTL article cache and invalidation on writes.
- [x] 2.4 Add tests for CRUD, pagination, comments, and cache invalidation.

## 3. HTTP API and Observability

- [x] 3.1 Implement REST routes for auth, articles, comments, health, and metrics.
- [x] 3.2 Add auth middleware for protected write routes.
- [x] 3.3 Wire structured trace IDs, metrics counters, and health endpoints.
- [x] 3.4 Add `httptest` coverage for register/login/article/comment/metrics smoke flow.

## 4. Delivery Materials and Verification

- [x] 4.1 Update README with architecture, route list, run commands, and completion checklist.
- [x] 4.2 Add `openapi.yaml`, `Dockerfile`, `docker-compose.yml`, and `EXERCISES.md`.
- [x] 4.3 Run capstone tests, capstone demo, full tests, vet, and build; fix failures.
