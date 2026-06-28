## Why

第 08 章目前仍是占位章节，尚未提供可运行 Web 服务示例、HTTP handler 测试和练习材料。完成本章可以把阶段一的工程化基础推进到阶段二业务工程入口，让学习者掌握 `net/http`、路由、中间件、请求校验、JSON 编解码与 REST 状态码约定。

## What Changes

- 在 `stage-2-business/08-web-foundations/` 下新增可运行 HTTP 服务入口程序，提供健康检查、文章列表、文章创建和文章详情等 REST 示例。
- 新增主题子包，分别演示标准库 `net/http`/`ServeMux`、chi 路由与 URL 参数、中间件链、JSON 请求/响应、validator 参数校验、请求 context 传播与 REST 错误响应。
- 补充 `_test.go`，使用 `httptest` 覆盖路由、handler、中间件、校验失败、JSON 响应与 context 传播场景。
- 将本章 `README.md` 从占位内容更新为实际产出说明、自测清单和运行命令，并新增 `EXERCISES.md`。
- 引入轻量业务 Web 依赖：`github.com/go-chi/chi/v5` 与 `github.com/go-playground/validator/v10`。

## Capabilities

### New Capabilities
- `web-foundations-tutorial`: 覆盖第 08 章 Web 基础学习单元的可运行代码、HTTP 测试、文档和练习。

### Modified Capabilities
- `learning-curriculum`: 将第 08 章从未落地占位章节更新为已落地章节，允许其目录包含源码、测试与练习材料。

## Impact

- 主要影响目录：`stage-2-business/08-web-foundations/`。
- 新增 OpenSpec 规格：`openspec/changes/chapter-08-web-foundations/specs/web-foundations-tutorial/spec.md`。
- 修改现有规格：`openspec/changes/chapter-08-web-foundations/specs/learning-curriculum/spec.md`。
- 更新依赖：`go.mod` / `go.sum` 增加 chi 与 validator 相关依赖。
- 验证命令包括 `go test ./stage-2-business/08-web-foundations/...`、`go run ./stage-2-business/08-web-foundations`、`go test ./...` 和 `go build ./...`。
