## ADDED Requirements

### Requirement: 第 05 章 SHALL 提供可运行的入口程序

第 05 章目录 `stage-1-syntax/05-concurrency/` SHALL 包含一个 `package main` 的入口文件 `main.go`，其 `main` 函数 MUST 通过调用各主题子包的导出函数来组装并输出一份"并发编程学习报告"，且 MUST 能通过 `go run` 成功执行。

#### Scenario: go run 跑通入口程序
- **WHEN** 学习者在 `stage-1-syntax/05-concurrency/` 目录执行 `go run .`
- **THEN** 程序 MUST 以退出码 0 结束，并在标准输出打印包含 goroutine、channel、sync、context、并发坑演示内容的报告文本

#### Scenario: 入口程序演示多子包协作
- **WHEN** 阅读者查看 `main.go`
- **THEN** 该文件 MUST 通过 `import` 引入本章至少三个主题子包（`goroutine` / `channel` / `syncx` / `ctx` / `pitfall`），并调用其导出函数，而非把全部逻辑写在 `main.go` 内

### Requirement: 第 05 章 SHALL 通过 goroutine 子包演示 goroutine 生命周期

第 05 章 SHALL 包含子目录 `goroutine/`（`package goroutine`），其中 MUST 演示：goroutine 启动、并发任务完成等待、结果收集。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/05-concurrency/goroutine`。

#### Scenario: goroutine 子包演示任务并发执行
- **WHEN** 学习者调用 goroutine 子包中的导出函数启动多个任务
- **THEN** 该函数 MUST 使用 goroutine 执行任务，并通过同步机制等待完成后返回可断言结果

#### Scenario: goroutine 子包不泄漏后台任务
- **WHEN** goroutine 子包导出函数返回
- **THEN** 它启动的 goroutine MUST 已经完成或被明确取消，且测试 MUST 不依赖无限等待

### Requirement: 第 05 章 SHALL 通过 channel 子包演示 channel 通信

第 05 章 SHALL 包含子目录 `channel/`（`package channel`），其中 MUST 演示：无缓冲 channel、有缓冲 channel、close + range、`select` timeout。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/05-concurrency/channel`。

#### Scenario: channel 子包演示无缓冲通信
- **WHEN** 学习者调用 channel 子包中的无缓冲通信导出函数
- **THEN** 该函数 MUST 使用无缓冲 channel 完成一次发送接收配对，并返回可断言结果

#### Scenario: channel 子包演示 close 与 range
- **WHEN** 学习者调用 channel 子包中的 close/range 导出函数
- **THEN** 该函数 MUST 关闭 channel 并通过 range 读取全部值，且返回值 MUST 可被单元测试断言

#### Scenario: channel 子包演示 select timeout
- **WHEN** 学习者调用 channel 子包中的 timeout 导出函数
- **THEN** 该函数 MUST 使用 `select` 和 `time.After` 或 timer 实现超时分支，且测试 MUST 能断言超时结果

### Requirement: 第 05 章 SHALL 通过 syncx 子包演示 sync 原语

第 05 章 SHALL 包含子目录 `syncx/`（`package syncx`），其中 MUST 演示：`sync.Mutex`、`sync.RWMutex`、`sync.WaitGroup`、`sync.Once` 的基础用法。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/05-concurrency/syncx`。

#### Scenario: syncx 子包演示 Mutex 保护共享状态
- **WHEN** 学习者调用 syncx 子包中并发累加的导出函数
- **THEN** 该函数 MUST 使用 `sync.Mutex` 或等价同步方式保护共享计数，并返回正确总数

#### Scenario: syncx 子包演示 Once 只执行一次
- **WHEN** 学习者调用 syncx 子包中多次并发触发初始化的导出函数
- **THEN** 该函数 MUST 使用 `sync.Once` 保证初始化逻辑只执行一次，且结果 MUST 有单元测试覆盖

### Requirement: 第 05 章 SHALL 通过 ctx 子包演示 context

第 05 章 SHALL 包含子目录 `ctx/`（`package ctx`），其中 MUST 演示：`context.Context` 的取消、超时与协作式退出。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/05-concurrency/ctx`。

#### Scenario: ctx 子包演示取消信号
- **WHEN** 学习者调用 ctx 子包中使用取消信号的导出函数
- **THEN** worker MUST 观察 `ctx.Done()` 后退出，并返回可断言的取消结果

#### Scenario: ctx 子包演示超时
- **WHEN** 学习者调用 ctx 子包中使用 timeout 的导出函数
- **THEN** 该函数 MUST 在超时后返回 `context deadline exceeded` 或等价可断言状态

### Requirement: 第 05 章 SHALL 通过 pitfall 子包说明常见并发坑

第 05 章 SHALL 包含子目录 `pitfall/`（`package pitfall`），其中 MUST 说明：data race、goroutine 泄漏、channel deadlock 的成因与规避方式。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/05-concurrency/pitfall`。

#### Scenario: pitfall 子包以安全方式说明风险
- **WHEN** 阅读者查看 `pitfall/pitfall.go`
- **THEN** 该文件 MUST 不包含会实际挂死测试或引入数据竞争的代码，但 MUST 通过注释或导出结果说明至少三类并发坑

#### Scenario: pitfall 子包提供可断言摘要
- **WHEN** 学习者调用 pitfall 子包的导出摘要函数
- **THEN** 该函数 MUST 返回包含 data race、goroutine leak、deadlock 或对应中文说明的可断言结果

### Requirement: 第 05 章 SHALL 提供单元测试

第 05 章 SHALL 在 `goroutine/`、`channel/`、`syncx/`、`ctx/`、`pitfall/` 各至少包含一个 `_test.go` 测试文件，且 MUST 采用表驱动（table-driven）+ `t.Run` 子测试的写法。

#### Scenario: go test -race 全部通过
- **WHEN** 学习者在仓库根目录执行 `go test -race ./stage-1-syntax/05-concurrency/...`
- **THEN** 所有测试 MUST 通过，命令 MUST 以退出码 0 返回

#### Scenario: 各子包测试采用表驱动风格
- **WHEN** 阅读者查看本章任一 `_test.go` 文件
- **THEN** 该文件 MUST 使用形如 `[]struct{ ... }` 的用例表并配合 `t.Run` 运行子测试

### Requirement: 第 05 章 SHALL 提供练习题与产出说明

第 05 章 SHALL 包含一个 `EXERCISES.md`，列出 3~5 道由浅入深的练习题（每题含验收标准）；同时章节 `README.md` 的"📦 本章产出"段落 MUST 被替换为实际内容（示例文件清单 + 运行命令），"✅ 自测清单"MUST 列出可勾选的检查项，且 MUST NOT 再包含"待 OpenSpec change 填充"或"尚未实现"等占位语。

#### Scenario: EXERCISES.md 含带验收标准的练习
- **WHEN** 阅读者打开 `stage-1-syntax/05-concurrency/EXERCISES.md`
- **THEN** 该文件 MUST 包含 3 到 5 道练习题，且每道题 MUST 给出明确的验收标准（如"运行 X 后输出 Y"或"测试通过条件"）

#### Scenario: README 产出说明不再是占位
- **WHEN** 阅读者打开 `stage-1-syntax/05-concurrency/README.md` 的"📦 本章产出"段落
- **THEN** 该段落 MUST NOT 再包含"待 OpenSpec change 填充"或"尚未实现"占位语，且 MUST 列出本章 `.go` 文件清单与运行命令（`go run .` / `go test -race ./...`）

#### Scenario: README 自测清单与知识点对齐
- **WHEN** 阅读者查看该 README 的"✅ 自测清单"
- **THEN** 该清单 MUST 包含与 ROADMAP 关键知识点对应的可勾选项（如 goroutine、channel、select、sync、context、data race、goroutine 泄漏、deadlock 等）
