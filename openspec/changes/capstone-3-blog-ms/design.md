## Context

Capstone 3 是阶段三的收束项目，需要在不依赖外部数据库、注册中心或 tracing 后端的默认测试中，演示三个限界上下文、gRPC、HTTP Gateway、链路追踪和韧性治理。仓库已有 gRPC/protobuf 与 `gobreaker/v2` 依赖，也已有 Capstone 2 的博客领域模型可供概念参考，但本项目必须保持目录独立并体现微服务边界。

## Goals / Non-Goals

**Goals:**

- 三个服务拥有独立内存存储、应用服务和 gRPC transport，可独立启动与测试。
- Gateway 提供注册、登录、文章和评论的 HTTP 聚合 API。
- 请求携带同一 trace 上下文穿过 Gateway 和全部 gRPC 调用，并可在测试中验证。
- 文章详情聚合作者与评论；评论服务不可用时返回显式降级结果。
- Gateway 对评论聚合调用应用超时、重试、熔断和限流。
- 默认端到端测试完全进程内运行，不依赖 Docker 或外部服务。
- Docker Compose 可启动三个服务和 Gateway。

**Non-Goals:**

- 实现生产级持久化、分布式事务、消息总线、服务注册中心或全局限流。
- 与 Capstone 2 共用源码或迁移其运行数据。
- 实现完整社交功能、富文本、搜索、附件和管理后台。
- 把教学用内存 tracing exporter 伪装成 OpenTelemetry SDK 的生产替代品。

## Decisions

### 每个限界上下文使用独立进程与内存仓储

`user-svc`、`post-svc`、`comment-svc` 各自拥有模型和仓储，跨上下文只交换 ID 与 protobuf DTO。相比共享数据库或共享 Go model 包，这能清晰展示数据所有权；代价是示例重启后数据丢失。

### 使用单个 protobuf 文件承载教学合同

合同集中在 `api/blog/v1/blog.proto`，生成代码提交到仓库。相比拆成多个 proto module，这更适合一周 capstone，并保留服务级 package 与 RPC 边界。生成文件只能通过修改 `.proto` 后重新生成。

### HTTP Gateway 聚合，服务间使用 gRPC

客户端使用易于实验的 JSON/HTTP；Gateway 通过 gRPC 调用后端。文章写入只调用 post-svc，文章详情再并行查询 user-svc 与 comment-svc。避免分布式写事务，使失败语义可解释。

### 自包含 W3C trace context 与内存 exporter

使用 `traceparent` 格式生成/解析 trace ID 和 span ID，并通过 gRPC metadata 传播；每层记录结构化 span。默认测试断言同一 trace 覆盖所有服务。生产系统应替换为 OpenTelemetry SDK/Collector。

### 韧性只包围评论聚合依赖

评论列表是文章详情的可降级附加数据。Gateway 对该调用执行短超时、有限重试和 circuit breaker；失败时返回空评论并设置 `comments_degraded: true`。用户与文章主体失败则返回错误，不静默伪造核心数据。

### 使用进程内 bufconn 做端到端测试

测试启动三个真实 gRPC server 和一个真实 HTTP Gateway，但连接使用 `bufconn`，既覆盖序列化与拦截器，又避免随机端口和 Docker 依赖。独立二进制仍监听 TCP 供 Compose 使用。

## Risks / Trade-offs

- [内存仓储无法展示持久化与数据迁移] → README 明确该边界，并把数据库演进列为练习。
- [手写 trace context 容易被误认为完整 OTel] → 命名与文档明确其教学用途，列出替换为 OTel SDK 的步骤。
- [Gateway 重试可能放大故障] → 仅重试可判定的临时 gRPC 状态，限制尝试次数，并让所有尝试共享总超时。
- [跨服务引用可能悬空] → 写入时由 Gateway 验证作者和文章存在；说明生产环境仍需事件驱动校验或补偿。
- [单一 proto 文件随功能增长会膨胀] → 本项目保持小型合同；练习要求按 bounded context 拆分。

## Migration Plan

1. 先生成并验证 protobuf 合同。
2. 实现三个服务及其独立测试。
3. 实现 trace、韧性组件和 Gateway。
4. 加入进程内端到端测试与可执行入口。
5. 补齐 Docker、Compose、README、练习和路线图。

本 change 只替换占位目录，无线上数据迁移。回滚时可删除新增实现并恢复占位 README。

## Open Questions

- 生产演进时优先选择数据库 outbox 事件同步，还是保留同步校验并增加补偿？
- trace、指标与日志在部署平台上应由应用直接导出，还是统一经过 OpenTelemetry Collector？
