## Context

`stage-1-syntax/capstone-1-cli-todo/` 当前只有占位 README 与 `.gitkeep`。阶段一常规章节已经覆盖 Go 的基础语法、复合类型、接口错误、并发、标准库与工程化，本 capstone 需要用一个完整但小型的 CLI Todo 项目把这些能力串起来。

项目必须保持阶段一约束：只用标准库，不引入第三方 CLI 框架；所有行为可本地测试；数据持久化使用临时目录或可配置路径，避免测试污染用户环境。

## Goals / Non-Goals

**Goals:**

- 提供可运行的 CLI Todo 程序，支持 add / list / done / delete / clear / help。
- 使用 struct + slice 建模任务集合，使用 JSON 文件持久化本地数据。
- 使用自定义错误区分命令解析、任务不存在、存储读写失败等场景。
- 提供受控异步保存示例，展示 goroutine + channel 的协作式退出。
- 提供表驱动测试、benchmark、README 和 EXERCISES。
- 保持 `go test ./...` 与 `go build ./...` 通过。

**Non-Goals:**

- 不实现多用户、网络同步、数据库或复杂 TUI。
- 不引入 cobra/urfave/cli 等第三方框架。
- 不依赖用户真实 home 目录作为测试数据目录。
- 不实现长期后台 daemon；异步保存仅作为本进程内受控示例。

## Decisions

### 1. 使用标准库手写命令分发

根目录 `main.go` 调用 `app.Run(args, stdout, stderr)`，`app` 包负责解析子命令并调用 domain/store 逻辑。

**Rationale:** 阶段一重点是基础语法与标准库，手写分发比引入 CLI 框架更适合教学。

**Alternative considered:** 使用 cobra。该方案更接近生产 CLI，但会把学习重点转移到第三方框架。

### 2. 分层为 app / todo / store / asyncsave

- `todo/`：任务模型、列表操作、自定义错误。
- `store/`：JSON 文件 load/save。
- `asyncsave/`：goroutine + channel 包装保存操作。
- `app/`：命令解析、输出格式、错误映射。
- 根目录 `main.go`：真实 CLI 入口。

**Rationale:** 保持结构清晰，同时不过度架构化，方便测试各层。

**Alternative considered:** 全部放在 main 包。该方案文件少，但测试和复用差，也不利于展示阶段一综合能力。

### 3. 默认数据路径可被环境变量覆盖

CLI 默认使用 `JUST_GO_TODO_FILE` 环境变量指定数据文件；未设置时使用当前目录下 `.just-go-todos.json`。

**Rationale:** 测试可以用临时文件隔离，用户也能明确指定文件。

**Alternative considered:** 使用用户 home 目录。该方案更像真实工具，但测试和清理更复杂。

### 4. 异步保存必须可关闭、可等待

`asyncsave` 包提供 worker，接收保存请求，支持 `Close()` 等待 goroutine 退出；测试确保无无限阻塞。

**Rationale:** 展示 goroutine + channel 同时避免 goroutine leak。

**Alternative considered:** 每次保存都 fire-and-forget。该方案容易丢数据和泄漏 goroutine。

## Risks / Trade-offs

- [Risk] CLI 功能过多导致 capstone 过大 → Mitigation：限定为本地 JSON 单文件 Todo，不做 TUI/网络/数据库。
- [Risk] 异步保存让测试变复杂 → Mitigation：核心 store 仍提供同步 Save；异步 worker 单独测试，app 可用同步保存保持命令确定性。
- [Risk] 默认数据文件污染仓库 → Mitigation：README 标注 `JUST_GO_TODO_FILE`，测试全部使用临时目录；`.just-go-todos.json` 可加入忽略说明但不必修改全局配置。
- [Risk] benchmark 结果随机器变化 → Mitigation：只要求 benchmark 可运行，不固定性能数值。
