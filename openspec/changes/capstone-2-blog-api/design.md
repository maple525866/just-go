# Design: Capstone 2 Blog API

## Overview

The capstone is a production-shaped monolith that stays runnable without external services. It uses in-memory repositories and JWT-like HMAC tokens so tests remain deterministic, while the package boundaries mirror a real service.

## Packages

- `model`: users, articles, comments, pagination, and API request/response types.
- `auth`: password hashing and signed bearer tokens.
- `store`: concurrency-safe in-memory repository for users, articles, tags, and nested comments.
- `cache`: small TTL article cache plus invalidation hooks for list/detail reads.
- `observability`: glue for request metrics and health checks.
- `server`: HTTP API, middleware, validation, auth, and routing.
