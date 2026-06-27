## Context

`stage-1-syntax/07-engineering/` 当前只有占位 README 与 `.gitkeep`。第 07 章是阶段一最后一个常规章节，重点从语言能力转向工程能力：测试、基准测试、质量门禁、调试和性能分析。

仓库已经配置 GitHub Actions 执行 `go vet`、`go test -race -count=1 ./...`、`go build ./...` 和 `golangci-lint`。本章不重复搭 CI，而是提供可运行示例和命令说明，帮助学习者理解这些门禁在本地如何使用。

## Goals / Non-Goals

**Goals:**

- 提供可运行的第 07 章入口程序，输出工程化基础学习报告。
- 提供适合表驱动测试、子测试和 benchmark 的小型函数。
- 用代码和文档说明 module、go.work、语义化版本、go vet、golangci-lint、dlv、slog、pprof。
- 让 `go test` 和 `go test -bench` 都能在本章目录内成功运行。
- README 与练习题替换占位内容，形成完整学习单元。

**Non-Goals:**

- 不引入新的 CI workflow 或修改已有 CI 语义。
- 不要求本章测试实际安装和运行 `dlv`、`golangci-lint` 或打开 pprof UI。
- 不生成真实 profile 文件作为仓库产物，避免污染工作目录。
- 不把阶段二业务工程内容提前引入本章。

## Decisions

### 1. 使用主题子包承载工程概念

- `moduleinfo/`：说明 module、go.work 与语义化版本概念。
- `calc/`：提供适合表驱动测试、子测试和 benchmark 的小函数。
- `quality/`：汇总 go vet、golangci-lint 和 CI 本地对齐命令。
- `debugx/`：演示 `log/slog` 输出与 dlv/IDE 调试命令摘要。
- `profile/`：说明 pprof CPU / 内存 / 阻塞 profile 的使用入口和适用场景。

**Rationale:** 工程化概念中很多是命令和流程，子包应提供可断言摘要和小型可 benchmark 代码，而不是伪造外部工具运行结果。

**Alternative considered:** 只写 README，不写代码。该方案更像文档章节，但不符合本仓库“每章独立可运行 demo”的风格。

### 2. benchmark 使用确定性纯函数

`calc` 包提供 `NormalizeWords`、`Fibonacci` 或类似纯函数供 benchmark 使用，避免 I/O 和并发噪音影响结果。

**Rationale:** 初学者学习 benchmark 时需要先理解“同一输入、多次运行、可比较”的基本原则。

**Alternative considered:** benchmark HTTP 或文件操作。该方案更真实，但阶段一容易引入过多变量。

### 3. 外部工具以命令清单和说明呈现

`golangci-lint`、`dlv`、`pprof` 不作为单元测试前置依赖；本章通过 README、EXERCISES 和导出摘要说明如何运行。

**Rationale:** 保持仓库本地验证轻量稳定，同时不回避真实工程工具。

**Alternative considered:** 在测试中直接执行这些工具。该方案容易因环境未安装而失败，不适合作为学习仓库基础门禁。

## Risks / Trade-offs

- [Risk] 工程化章节代码示例偏抽象 → Mitigation：入口报告和 README 明确命令与真实使用场景，EXERCISES 要求实际运行 benchmark / vet。
- [Risk] 不实际运行 dlv/pprof 可能显得不够实战 → Mitigation：提供可复制命令和练习验收，让学习者在本地手动完成。
- [Risk] benchmark 结果在不同机器不同 → Mitigation：验收只要求命令跑通和理解指标，不固定 ns/op 数值。
- [Risk] golangci-lint 本地未安装 → Mitigation：README 区分 CI 已配置与本地可选安装，并保留 `go vet` / `go test` / `go build` 基础验证。
