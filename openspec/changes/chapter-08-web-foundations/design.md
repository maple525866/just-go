## Context

`stage-2-business/08-web-foundations/` 当前只有占位 README 与 `.gitkeep`。第 08 章是阶段二入口，目标是从阶段一的 CLI/标准库示例过渡到可测试的 HTTP 服务，覆盖 `net/http`、路由、中间件、参数校验、JSON 编解码、REST 状态码与 request context 生命周期。

本仓库每章都是独立可运行 demo，因此本章应提供可 `go run` 的小服务，也要提供可在无外部服务依赖下稳定运行的 `httptest` 单元测试。阶段二后续章节会引入数据库、缓存、可观测性，本章只保留内存数据和 Web 层概念。

## Goals / Non-Goals

**Goals:**

- 提供可运行的第 08 章入口程序，启动一个小型博客文章 HTTP API。
- 同时演示标准库 `http.ServeMux` 与 chi router，让学习者理解二者职责差异和 URL 参数提取。
- 提供可组合中间件链，覆盖 Recover、请求日志、CORS、简易限流和 request ID/context 传播。
- 使用 `go-playground/validator` 对 JSON 请求体做参数校验，并返回结构化错误响应。
- 使用 `httptest` 为 handler、router、middleware、JSON 编解码、校验失败和状态码约定提供测试。
- README 与练习题替换占位内容，形成完整学习单元。

**Non-Goals:**

- 不引入数据库、Redis、消息队列或持久化；文章数据使用内存 store。
- 不实现生产级认证、权限、分布式限流或复杂 CORS 策略。
- 不使用 Gin/Echo；本章以标准库 + chi 形成对比，避免同时引入过多框架。
- 不提供 OpenAPI/Swagger；接口文档以后续 capstone 综合项目处理。

## Decisions

### 1. 使用小型博客文章 API 作为贯穿案例

本章 API 包含 `GET /healthz`、`GET /api/articles`、`POST /api/articles`、`GET /api/articles/{id}` 四类端点。文章模型只包含 ID、Title、Body 与 Tags，足以演示列表、创建、详情、JSON 编解码和 REST 状态码。

**Rationale:** 博客文章 API 与阶段二 capstone 主题一致，能为后续数据库、缓存和可观测性章节铺垫，同时又足够小，不会把业务复杂度压到 Web 基础章节。

**Alternative considered:** 使用 Todo API。该方案简单，但阶段一 capstone 已覆盖 Todo，阶段二继续使用博客语境更连贯。

### 2. 子包按 Web 责任拆分

- `model/`：文章请求、响应、错误响应等 JSON 类型。
- `store/`：内存文章仓库，提供 handler 可测试依赖。
- `response/`：JSON 编码、错误响应、状态码约定。
- `middleware/`：Recover、日志、CORS、简易限流、request ID/context helper。
- `server/`：标准库 ServeMux 示例与 chi router 示例。
- `validation/`：validator 封装与错误字段格式化。

**Rationale:** 初学者需要看到 handler 不是孤立函数，而是围绕 model、store、response、middleware 协作；拆包能避免一个巨大 `main.go`，也延续前序章节“主题子包 + 可断言函数”的风格。

**Alternative considered:** 所有内容放在 `main.go`。该方案更短，但不利于演示包职责、测试粒度和后续扩展。

### 3. HTTP 测试全部使用 `httptest`

handler 和 middleware 测试使用 `httptest.NewRequest`、`httptest.NewRecorder` 与 router 的 `ServeHTTP`，不启动真实网络端口。

**Rationale:** 测试稳定、无端口冲突、无外部依赖，并能精准断言状态码、响应头和 JSON body。

**Alternative considered:** 在测试中 `ListenAndServe` 启动真实 server。该方案更接近运行环境，但会引入端口占用和并发清理问题。

### 4. 依赖选择：chi + validator

chi 用于演示主流 router 的方法路由、URL 参数和 middleware 链；validator 用于结构体 tag 参数校验。标准库 ServeMux 示例保留在 `server` 子包中，帮助学习者理解框架不是必需品。

**Rationale:** 两个依赖轻量、社区常见，与 ROADMAP 推荐项一致。chi 更贴近 `net/http` 接口，适合从标准库过渡。

**Alternative considered:** Gin。Gin 功能完整但抽象更多，本章重点是理解 `http.Handler` 生命周期，chi 更贴合。

## Risks / Trade-offs

- [Risk] 引入第三方依赖可能增加初学者心智负担 → Mitigation：README 明确先看标准库 ServeMux，再看 chi 的增量价值。
- [Risk] 内存 store 可能被误解为生产持久化方案 → Mitigation：文档和练习明确说明持久化留到第 09 章。
- [Risk] 简易限流不是生产级 → Mitigation：命名和 README 标注为教学版，只演示 middleware 形态与 429 状态码。
- [Risk] handler 示例过多导致本章范围膨胀 → Mitigation：仅保留文章列表、创建、详情三个核心业务端点。
