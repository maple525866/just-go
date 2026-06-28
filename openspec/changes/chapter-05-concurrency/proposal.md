## Why

第 05 章目前仍是占位章节，尚未提供可运行示例、单元测试和练习材料。完成本章可以承接接口与错误处理基础，帮助学习者掌握 Go 最具代表性的并发模型与常见并发风险。

## What Changes

- 在 `stage-1-syntax/05-concurrency/` 下新增可运行入口程序，串联 goroutine、channel、sync、context 和并发坑示例。
- 新增多个主题子包，分别演示 goroutine 生命周期、channel 通信与关闭、`select` 超时、`sync` 常用原语、`context` 取消 / 超时，以及 data race / goroutine 泄漏 / deadlock 的规避说明。
- 为各主题子包补充可确定的表驱动单元测试，避免依赖不稳定的时间竞态。
- 将本章 `README.md` 从占位内容更新为实际产出说明、自测清单和运行命令，并新增 `EXERCISES.md`。
- 不引入第三方依赖，不修改公共模块名，不影响已完成章节的运行方式。

## Capabilities

### New Capabilities
- `concurrency-tutorial`: 覆盖第 05 章 Go 并发编程学习单元的可运行代码、测试、文档和练习。

### Modified Capabilities
- `learning-curriculum`: 将第 05 章从未落地占位章节更新为已落地章节，允许其目录包含源码、测试与练习材料。

## Impact

- 主要影响目录：`stage-1-syntax/05-concurrency/`。
- 可能新增 OpenSpec 规格：`openspec/changes/chapter-05-concurrency/specs/concurrency-tutorial/spec.md`。
- 可能修改现有规格：`openspec/changes/chapter-05-concurrency/specs/learning-curriculum/spec.md`。
- 验证命令包括 `go test -race ./stage-1-syntax/05-concurrency/...`、`go run ./stage-1-syntax/05-concurrency`、`go test ./...` 和 `go build ./...`。
