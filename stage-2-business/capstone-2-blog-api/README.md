# 🚩 Capstone 2: 博客后端 API（单体）

> 阶段：② 业务工程 · 阶段综合项目 ｜ 难度：⭐⭐⭐⭐☆ ｜ 预计耗时：5 天

## 🎯 项目目标

做一个**生产级单体博客后端**的教学版，实现用户注册 / 登录、文章 CRUD、标签分页、嵌套评论、缓存读路径与全套可观测性端点。为了保证学习者无需安装 MySQL、Redis 或 Collector，本项目使用内存仓库、TTL 缓存和 HMAC bearer token；目录边界按真实单体服务设计，后续可以替换为 GORM、Redis、Prometheus 和 OpenTelemetry SDK。

## 🧩 综合应用的章节

- [✓ 08-web-foundations] HTTP 路由、中间件思路、JSON 编解码、REST 状态码
- [✓ 09-data-persistence] Repository 形态、事务边界意识、分页与标签过滤
- [✓ 10-cache-and-mq] TTL 缓存、写后失效、异步事件/消息队列的扩展点
- [✓ 11-observability] 健康检查、Prometheus 风格 metrics、trace ID 传播

## 📦 项目结构

```text
stage-2-business/capstone-2-blog-api/
├── main.go              # 输出 capstone 服务报告
├── auth/                # 密码哈希、HMAC bearer token 签发与校验
├── model/               # User、Article、Comment、分页与请求类型
├── store/               # 并发安全内存仓库：用户、文章、标签、嵌套评论
├── cache/               # 文章详情 TTL 缓存与写后失效
├── observability/       # metrics + health 组合入口
├── server/              # REST API、认证、路由、metrics、health
├── openapi.yaml         # API 文档
├── Dockerfile           # 容器构建示例
├── docker-compose.yml   # 一键启动示例
└── EXERCISES.md         # 阶段答辩练习
```

## 🚀 运行

测试：

```bash
go test ./stage-2-business/capstone-2-blog-api/...
```

运行学习报告：

```bash
go run ./stage-2-business/capstone-2-blog-api
```

容器构建与启动示例：

```bash
docker compose -f stage-2-business/capstone-2-blog-api/docker-compose.yml up --build
```

服务默认监听 `:8080`，可通过 `ADDR=:8081` 覆盖。

## 🌐 API 速览

| 方法 | 路径 | 认证 | 说明 |
|---|---|---:|---|
| POST | `/api/register` | 否 | 注册用户 |
| POST | `/api/login` | 否 | 登录并返回 bearer token |
| GET | `/api/articles?tag=go&page=1&page_size=10` | 否 | 文章列表、标签过滤、分页 |
| POST | `/api/articles` | 是 | 创建文章 |
| GET | `/api/articles/{id}` | 否 | 文章详情，命中 TTL cache |
| PUT | `/api/articles/{id}` | 是 | 更新文章并失效缓存 |
| DELETE | `/api/articles/{id}` | 是 | 删除文章并失效缓存 |
| POST | `/api/articles/{id}/comments` | 是 | 新增根评论或带 `parent_id` 的回复 |
| GET | `/livez` | 否 | liveness 健康检查 |
| GET | `/readyz` | 否 | readiness 健康检查 |
| GET | `/metrics` | 否 | Prometheus 风格指标 |

示例请求：

```bash
curl -X POST http://localhost:8080/api/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","password":"secret"}'

TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","password":"secret"}' | jq -r .token)

curl -X POST http://localhost:8080/api/articles \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"title":"Stage 2","body":"capstone","tags":["go","api"]}'
```

## 🔐 安全说明

本项目的密码哈希和 HMAC bearer token 是为了教学可读性而实现的最小版本：生产环境应使用 bcrypt/argon2id 存储密码，并使用包含过期时间、签发者、受众等声明的 JWT 或会话方案。

## 🏗️ 分层说明

```text
server  ── HTTP 协议、认证、状态码、JSON
  │
  ├── auth ── 密码与 token
  ├── store ── 业务数据读写接口形态
  ├── cache ── 文章详情缓存与失效
  └── observability / chapter-11 packages ── health、metrics、trace
```

这个分层刻意保持简单：业务还在单体内，但 HTTP、认证、数据、缓存和可观测性已经拆开，方便第 12 章继续演进到 clean architecture。

## ✅ 完成标准

- [x] 所有列出的章节知识点至少综合使用过一次。
- [x] 代码可运行、有测试、有 README。
- [x] 接口文档 `openapi.yaml` 齐全。
- [x] 端到端冒烟测试覆盖注册、登录、创建文章、列表、评论和 metrics。
- [x] 阶段答辩能讲清楚当前分层，以及未来替换 MySQL / Redis / OTel 的位置。
