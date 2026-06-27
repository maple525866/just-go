## ADDED Requirements

### Requirement: 第 07 章 SHALL 提供可运行的入口程序

第 07 章目录 `stage-1-syntax/07-engineering/` SHALL 包含一个 `package main` 的入口文件 `main.go`，其 `main` 函数 MUST 通过调用各主题子包的导出函数来组装并输出一份"工程化基础学习报告"，且 MUST 能通过 `go run` 成功执行。

#### Scenario: go run 跑通入口程序
- **WHEN** 学习者在 `stage-1-syntax/07-engineering/` 目录执行 `go run .`
- **THEN** 程序 MUST 以退出码 0 结束，并在标准输出打印包含 module/testing/benchmark/lint/debug/pprof 演示内容的报告文本

#### Scenario: 入口程序演示多子包协作
- **WHEN** 阅读者查看 `main.go`
- **THEN** 该文件 MUST 通过 `import` 引入本章至少三个主题子包（`moduleinfo` / `calc` / `quality` / `debugx` / `profile`），并调用其导出函数，而非把全部逻辑写在 `main.go` 内

### Requirement: 第 07 章 SHALL 通过 moduleinfo 子包说明模块与版本语义

第 07 章 SHALL 包含子目录 `moduleinfo/`（`package moduleinfo`），其中 MUST 说明：Go module、`go.work`、语义化版本的基础概念。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/07-engineering/moduleinfo`。

#### Scenario: moduleinfo 子包提供模块摘要
- **WHEN** 学习者调用 moduleinfo 子包的导出函数
- **THEN** 该函数 MUST 返回包含 module、go.work、semantic version 关键概念的可断言摘要

### Requirement: 第 07 章 SHALL 通过 calc 子包演示 testing 与 benchmark

第 07 章 SHALL 包含子目录 `calc/`（`package calc`），其中 MUST 提供适合表驱动测试、子测试和 benchmark 的纯函数。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/07-engineering/calc`。

#### Scenario: calc 子包提供可测试纯函数
- **WHEN** 学习者调用 calc 子包的导出函数
- **THEN** 该函数 MUST 对固定输入返回确定结果，且 MUST 有表驱动 + `t.Run` 子测试覆盖

#### Scenario: calc 子包提供 benchmark
- **WHEN** 学习者执行 `go test -bench=. ./stage-1-syntax/07-engineering/...`
- **THEN** 命令 MUST 能运行至少一个 benchmark 函数并以退出码 0 返回

### Requirement: 第 07 章 SHALL 通过 quality 子包说明质量门禁

第 07 章 SHALL 包含子目录 `quality/`（`package quality`），其中 MUST 说明：`go vet`、`go test -race`、`go build`、`golangci-lint` 与 CI 对齐。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/07-engineering/quality`。

#### Scenario: quality 子包提供本地验证命令
- **WHEN** 学习者调用 quality 子包的导出函数
- **THEN** 该函数 MUST 返回包含 `go vet ./...`、`go test -race -count=1 ./...`、`go build ./...`、`golangci-lint run` 的命令清单

### Requirement: 第 07 章 SHALL 通过 debugx 子包说明调试与日志

第 07 章 SHALL 包含子目录 `debugx/`（`package debugx`），其中 MUST 演示：`log/slog` 日志输出，并说明 `dlv` / IDE 断点调试入口。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/07-engineering/debugx`。

#### Scenario: debugx 子包演示 slog 输出
- **WHEN** 学习者调用 debugx 子包的导出函数
- **THEN** 该函数 MUST 使用 `log/slog` 写入可断言的日志文本，且行为 MUST 有单元测试覆盖

#### Scenario: debugx 子包提供调试命令摘要
- **WHEN** 学习者调用 debugx 子包的调试摘要函数
- **THEN** 该函数 MUST 返回包含 `dlv test` 或 IDE breakpoint 的调试入口说明

### Requirement: 第 07 章 SHALL 通过 profile 子包说明 pprof

第 07 章 SHALL 包含子目录 `profile/`（`package profile`），其中 MUST 说明：pprof CPU / memory / blocking profile 的用途与命令入口。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/07-engineering/profile`。

#### Scenario: profile 子包提供 profile 类型摘要
- **WHEN** 学习者调用 profile 子包的导出函数
- **THEN** 该函数 MUST 返回至少三类 profile（CPU、memory、blocking）的用途说明

#### Scenario: profile 子包提供 pprof 命令提示
- **WHEN** 学习者调用 profile 子包的命令摘要函数
- **THEN** 该函数 MUST 返回包含 `go test -bench`、`-cpuprofile` 或 `go tool pprof` 的命令提示

### Requirement: 第 07 章 SHALL 提供单元测试与 benchmark

第 07 章 SHALL 在各主题子包至少包含一个 `_test.go` 测试文件或由相邻测试覆盖，且 MUST 采用表驱动（table-driven）+ `t.Run` 子测试的写法；同时 MUST 至少包含一个 `BenchmarkXxx` 函数。

#### Scenario: go test 全部通过
- **WHEN** 学习者在仓库根目录执行 `go test ./stage-1-syntax/07-engineering/...`
- **THEN** 所有测试 MUST 通过，命令 MUST 以退出码 0 返回

#### Scenario: benchmark 可运行
- **WHEN** 学习者在仓库根目录执行 `go test -bench=. ./stage-1-syntax/07-engineering/...`
- **THEN** 至少一个 benchmark MUST 被执行，命令 MUST 以退出码 0 返回

### Requirement: 第 07 章 SHALL 提供练习题与产出说明

第 07 章 SHALL 包含一个 `EXERCISES.md`，列出 3~5 道由浅入深的练习题（每题含验收标准）；同时章节 `README.md` 的"📦 本章产出"段落 MUST 被替换为实际内容（示例文件清单 + 运行命令），"✅ 自测清单"MUST 列出可勾选的检查项，且 MUST NOT 再包含"待 OpenSpec change 填充"或"尚未实现"等占位语。

#### Scenario: EXERCISES.md 含带验收标准的练习
- **WHEN** 阅读者打开 `stage-1-syntax/07-engineering/EXERCISES.md`
- **THEN** 该文件 MUST 包含 3 到 5 道练习题，且每道题 MUST 给出明确的验收标准（如"运行 X 后输出 Y"或"测试通过条件"）

#### Scenario: README 产出说明不再是占位
- **WHEN** 阅读者打开 `stage-1-syntax/07-engineering/README.md` 的"📦 本章产出"段落
- **THEN** 该段落 MUST NOT 再包含"待 OpenSpec change 填充"或"尚未实现"占位语，且 MUST 列出本章 `.go` 文件清单与运行命令（`go run .` / `go test ./...` / `go test -bench=. ./...`）

#### Scenario: README 自测清单与知识点对齐
- **WHEN** 阅读者查看该 README 的"✅ 自测清单"
- **THEN** 该清单 MUST 包含与 ROADMAP 关键知识点对应的可勾选项（如 module、go.work、表驱动测试、benchmark、go vet、golangci-lint、dlv、slog、pprof 等）
