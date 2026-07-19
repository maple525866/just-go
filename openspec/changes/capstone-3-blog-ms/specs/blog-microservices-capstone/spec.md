## ADDED Requirements

### Requirement: 系统 SHALL 提供三个独立博客微服务

系统 SHALL 提供 `user-svc`、`post-svc` 和 `comment-svc`，每个服务 MUST 拥有独立的数据所有权、应用逻辑、gRPC 接口和可独立执行的测试。

#### Scenario: 服务独立启动

- **WHEN** 分别启动三个服务二进制
- **THEN** 每个服务 MUST 在未启动其他两个服务时完成初始化并监听自己的 gRPC 地址

#### Scenario: 服务独立测试

- **WHEN** 执行任一服务包的测试
- **THEN** 测试 MUST 不依赖 Docker、外部数据库或其他服务进程

### Requirement: user-svc SHALL 负责身份与鉴权

`user-svc` SHALL 支持用户注册、登录、令牌校验和按 ID 查询用户，并 MUST 拒绝重复用户名、错误密码和无效令牌。

#### Scenario: 注册并登录

- **WHEN** 客户端用合法用户名和密码注册后登录
- **THEN** 服务 MUST 返回可被后续请求验证的 bearer token

#### Scenario: 无效凭据

- **WHEN** 客户端使用错误密码或无效 token
- **THEN** 服务 MUST 返回可区分的未认证错误

### Requirement: post-svc SHALL 负责文章生命周期

`post-svc` SHALL 支持创建、查询、列表、更新和删除文章，并 MUST 使用作者 ID 表示跨服务所有权。

#### Scenario: 作者管理文章

- **WHEN** 已认证作者创建、更新或删除自己的文章
- **THEN** 服务 MUST 持久化对应变化并返回最新文章状态

#### Scenario: 非作者修改文章

- **WHEN** 其他用户尝试更新或删除文章
- **THEN** 服务 MUST 返回禁止访问错误且 MUST NOT 修改文章

### Requirement: comment-svc SHALL 负责嵌套评论与软删除

`comment-svc` SHALL 支持创建顶层或回复评论、按文章列出评论树和软删除评论。

#### Scenario: 创建嵌套评论

- **WHEN** 客户端为文章创建评论并以已有评论作为 parent
- **THEN** 列表结果 MUST 将该评论放入父评论的 replies

#### Scenario: 软删除评论

- **WHEN** 评论被删除
- **THEN** 评论节点 MUST 保留层级和回复，但正文 MUST 被清空并标记为 deleted

### Requirement: API Gateway SHALL 聚合后端服务

Gateway SHALL 提供 JSON/HTTP API，完成注册、登录、文章写入、文章列表、文章详情和评论写入，并 MUST 将 bearer token 交由 `user-svc` 验证。

#### Scenario: 获取聚合文章详情

- **WHEN** 客户端查询存在的文章
- **THEN** Gateway MUST 返回文章、作者信息和评论树

#### Scenario: 受保护写请求

- **WHEN** 客户端未携带有效 bearer token 发起写请求
- **THEN** Gateway MUST 返回 HTTP 401 且 MUST NOT 调用对应写 RPC

### Requirement: 系统 SHALL 传播端到端 trace context

Gateway 和三个 gRPC 服务 SHALL 使用 W3C `traceparent` 兼容格式传播 trace context，并 SHALL 为每个处理阶段记录结构化 span。

#### Scenario: 跨服务链路

- **WHEN** 一个文章详情请求经过 Gateway、post-svc、user-svc 和 comment-svc
- **THEN** 所有记录的 span MUST 共享同一个 trace ID，并包含正确的父子关系

#### Scenario: 接收外部 trace

- **WHEN** 客户端发送合法 `traceparent`
- **THEN** Gateway MUST 延续其 trace ID，并在响应中返回 trace 标识

### Requirement: Gateway SHALL 对可降级依赖应用韧性策略

Gateway SHALL 对评论聚合调用应用本地限流、总超时、有限重试和 circuit breaker，并 MUST 只对临时错误重试。

#### Scenario: 评论服务临时失败后恢复

- **WHEN** 评论 RPC 第一次返回临时错误且下一次成功
- **THEN** Gateway MUST 在超时预算内重试并返回完整详情

#### Scenario: 评论服务持续失败

- **WHEN** 评论 RPC 持续失败或 breaker 已打开
- **THEN** Gateway MUST 返回文章主体和作者、空评论列表以及显式 `comments_degraded` 标记

#### Scenario: 请求超过本地配额

- **WHEN** Gateway 的详情请求令牌桶耗尽
- **THEN** Gateway MUST 在调用后端前返回 HTTP 429

### Requirement: 项目 SHALL 提供自包含端到端验证

项目 SHALL 提供启动真实 gRPC server 和 HTTP Gateway 的进程内端到端测试，并 MUST 覆盖注册、登录、创建文章、添加嵌套评论、查询聚合详情和 trace 传播。

#### Scenario: 默认冒烟测试

- **WHEN** 在仓库根目录执行 Capstone 3 测试
- **THEN** 完整冒烟流程 MUST 在无需外部服务的情况下通过

### Requirement: 项目 SHALL 提供容器化运行环境

项目 SHALL 为三个服务和 Gateway 提供 Docker 构建入口，并 SHALL 提供 `docker-compose.yml`。

#### Scenario: Compose 启动完整系统

- **WHEN** 学习者执行 `docker compose up --build`
- **THEN** 三个 gRPC 服务和 Gateway MUST 被编排到同一网络，Gateway API MUST 可从宿主机访问
