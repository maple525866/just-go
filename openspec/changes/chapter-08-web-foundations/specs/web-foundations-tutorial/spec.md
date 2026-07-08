## ADDED Requirements

### Requirement: 第 08 章 SHALL 提供可运行的 HTTP 服务入口程序

第 08 章目录 `stage-2-business/08-web-foundations/` SHALL 包含一个 `package main` 的入口文件 `main.go`，其 `main` 函数 MUST 通过调用 `server` 子包构造 HTTP handler，并 MUST 支持以环境变量 `ADDR` 指定监听地址。

#### Scenario: go run 跑通入口程序
- **WHEN** 学习者在仓库根目录执行 `go run ./stage-2-business/08-web-foundations`
- **THEN** 程序 MUST 能成功编译并启动 HTTP 服务，默认监听地址 MUST 为 `:8080`

#### Scenario: 入口程序不承载全部业务逻辑
- **WHEN** 阅读者查看 `stage-2-business/08-web-foundations/main.go`
- **THEN** 该文件 MUST 通过 import 引入本章子包并调用其导出函数构造服务，而非把路由、handler、中间件和业务存储全部写在 `main.go` 内

### Requirement: 第 08 章 SHALL 演示标准库 ServeMux 与 chi router

第 08 章 SHALL 包含 `server` 子包，提供标准库 `http.ServeMux` 与 chi router 两种构造函数。chi router MUST 至少包含 `GET /healthz`、`GET /api/articles`、`POST /api/articles`、`GET /api/articles/{id}` 四类端点。

#### Scenario: chi router 支持 REST 端点
- **WHEN** 测试代码向 chi router 发送 `GET /healthz`、`GET /api/articles`、`POST /api/articles`、`GET /api/articles/{id}` 请求
- **THEN** router MUST 分别返回符合 REST 语义的状态码与 JSON 响应

#### Scenario: URL 参数用于文章详情
- **WHEN** 测试代码向 `GET /api/articles/{id}` 发送存在的文章 ID
- **THEN** handler MUST 使用 URL 参数读取 ID 并返回对应文章 JSON

#### Scenario: 标准库 ServeMux 示例可运行
- **WHEN** 测试代码使用 `server.NewStdMux` 构造 handler 并请求 `GET /healthz`
- **THEN** handler MUST 返回 200 状态码和 JSON 健康响应

### Requirement: 第 08 章 SHALL 提供 JSON 请求响应与 REST 错误约定

第 08 章 SHALL 包含 `model` 与 `response` 子包，定义文章创建请求、文章响应、列表响应、错误响应，并统一使用 `Content-Type: application/json` 输出 JSON。

#### Scenario: 创建文章返回 201 与 Location
- **WHEN** 客户端向 `POST /api/articles` 发送合法 JSON 请求体
- **THEN** 服务 MUST 返回 201 状态码、`Location` 响应头和新建文章 JSON

#### Scenario: 非法 JSON 返回 400
- **WHEN** 客户端向 `POST /api/articles` 发送无法解析的 JSON
- **THEN** 服务 MUST 返回 400 状态码，并返回包含错误码和消息的 JSON 错误响应

#### Scenario: 未找到文章返回 404
- **WHEN** 客户端请求不存在的 `GET /api/articles/{id}`
- **THEN** 服务 MUST 返回 404 状态码，并返回结构化 JSON 错误响应

### Requirement: 第 08 章 SHALL 使用 validator 做请求体参数校验

第 08 章 SHALL 包含 `validation` 子包，封装 `go-playground/validator/v10`，并对创建文章请求校验 title、body、tags 字段。title MUST 非空且长度不超过 80，body MUST 非空且长度不少于 10，tags 中每个标签 MUST 非空且长度不超过 20。

#### Scenario: 校验失败返回 422
- **WHEN** 客户端向 `POST /api/articles` 发送字段缺失或长度非法的 JSON 请求体
- **THEN** 服务 MUST 返回 422 状态码，并返回可读的字段级校验错误列表

#### Scenario: 校验成功创建文章
- **WHEN** 客户端发送 title、body、tags 均满足约束的 JSON 请求体
- **THEN** 服务 MUST 创建文章并返回包含 ID、title、body、tags 的 JSON 响应

### Requirement: 第 08 章 SHALL 提供可组合中间件链

第 08 章 SHALL 包含 `middleware` 子包，至少提供 Recover、请求日志、CORS、简易限流、request ID/context 传播中间件。中间件 MUST 遵循 `func(http.Handler) http.Handler` 形态。

#### Scenario: Recover 中间件捕获 panic
- **WHEN** 被包装的 handler 发生 panic
- **THEN** Recover 中间件 MUST 返回 500 状态码和 JSON 错误响应，而不是让 panic 逃逸到测试进程

#### Scenario: CORS 中间件写入跨域响应头
- **WHEN** 客户端发送带 Origin 的请求
- **THEN** CORS 中间件 MUST 写入 `Access-Control-Allow-Origin` 响应头

#### Scenario: 简易限流返回 429
- **WHEN** 同一限流器容量被耗尽后继续接收请求
- **THEN** 限流中间件 MUST 返回 429 状态码和 JSON 错误响应

#### Scenario: request ID 写入 context 与响应头
- **WHEN** 请求未携带 request ID
- **THEN** 中间件 MUST 生成 request ID，写入请求 context，并在响应头中返回同一个 request ID

### Requirement: 第 08 章 SHALL 提供 HTTP 单元测试

第 08 章 SHALL 为 `server`、`middleware`、`validation`、`response` 或相关子包提供 `_test.go` 测试，且 MUST 使用 `httptest` 覆盖真实 handler 行为。

#### Scenario: 章节测试全部通过
- **WHEN** 学习者在仓库根目录执行 `go test ./stage-2-business/08-web-foundations/...`
- **THEN** 所有测试 MUST 通过，命令 MUST 以退出码 0 返回

#### Scenario: 测试断言 HTTP 语义
- **WHEN** 阅读者查看本章测试文件
- **THEN** 测试 MUST 断言状态码、JSON 响应、关键响应头和至少一个错误响应场景

### Requirement: 第 08 章 SHALL 提供练习题与产出说明

第 08 章 SHALL 包含一个 `EXERCISES.md`，列出 3~5 道由浅入深的练习题（每题含验收标准）；同时章节 `README.md` 的“📦 本章产出”段落 MUST 被替换为实际内容（示例文件清单 + 运行命令），自测清单 MUST 与 ROADMAP 关键知识点对齐，且 MUST NOT 再包含“待 OpenSpec change 填充”或“尚未实现”等占位语。

#### Scenario: EXERCISES.md 含带验收标准的练习
- **WHEN** 阅读者打开 `stage-2-business/08-web-foundations/EXERCISES.md`
- **THEN** 该文件 MUST 包含 3 到 5 道练习题，且每道题 MUST 给出明确的验收标准

#### Scenario: README 产出说明不再是占位
- **WHEN** 阅读者打开 `stage-2-business/08-web-foundations/README.md` 的“📦 本章产出”段落
- **THEN** 该段落 MUST NOT 再包含“待 OpenSpec change 填充”或“尚未实现”占位语，且 MUST 列出本章 `.go` 文件清单与运行命令

#### Scenario: README 自测清单与知识点对齐
- **WHEN** 阅读者查看该 README 的“✅ 自测清单”
- **THEN** 该清单 MUST 包含与 ROADMAP 关键知识点对应的可勾选项，包括 `net/http`、router、中间件、validator、JSON、context 和 REST 状态码
