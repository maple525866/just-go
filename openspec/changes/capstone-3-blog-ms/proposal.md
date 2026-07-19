## Why

阶段三需要一个综合项目，把整洁架构、DDD、微服务通信、韧性治理与性能观测串成可运行的端到端案例。当前 Capstone 3 仍是占位目录，学习者无法验证如何从单体博客演进为边界清晰、可独立测试和部署的微服务系统。

Tracking issue: [#37](https://github.com/maple525866/just-go/issues/37)

## What Changes

- 实现 `user-svc`、`post-svc`、`comment-svc` 三个独立限界上下文，分别负责身份、文章和评论。
- 使用 gRPC 定义并实现服务间契约，并由 HTTP API Gateway 聚合为面向客户端的 API。
- 在 Gateway 中加入鉴权、限流、超时、重试、熔断和显式降级。
- 使用 W3C `traceparent` 传播端到端 trace，并提供结构化内存 span 导出，默认测试无需外部可观测平台。
- 提供 Dockerfile 与 docker-compose，支持一条命令启动完整示例。
- 增加端到端冒烟测试、服务级测试、中文 README、练习和架构取舍说明。
- 更新 `ROADMAP.md` 与 Capstone 3 完成清单。

## Capabilities

### New Capabilities

- `blog-microservices-capstone`: 定义三个博客微服务、gRPC 合同、API Gateway、端到端追踪、韧性策略、容器化交付和测试要求。

### Modified Capabilities

- `learning-curriculum`: 将 Capstone 3 从占位状态更新为已实现的微服务综合项目，并明确其学习产出。

## Impact

- 主要新增代码位于 `stage-3-architecture/capstone-3-blog-ms/`。
- 新增 protobuf 合同及生成的 Go 文件、服务实现、Gateway、测试和容器化配置。
- 复用仓库现有 gRPC、protobuf 和 `gobreaker/v2` 依赖，默认实现不要求数据库、消息队列或外部 tracing 服务。
- 更新 `ROADMAP.md` 和对应 OpenSpec change；不修改前三个阶段章节的既有行为。
