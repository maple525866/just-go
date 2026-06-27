## Why

第 07 章目前仍是占位章节，尚未提供可运行示例、单元测试和练习材料。完成本章可以把阶段一前 6 章的“会写 Go”推进到“会做 Go 工程”，覆盖测试、benchmark、质量门禁、调试与 profile 的基础闭环。

## What Changes

- 在 `stage-1-syntax/07-engineering/` 下新增可运行入口程序，串联 module/go.work 概念、testing、benchmark、lint、debug、pprof 示例与说明。
- 新增多个主题子包，分别演示表驱动测试对象、benchmark 对象、工程元信息、日志/调试摘要、pprof 风险和优化方向。
- 补充 `_test.go` 与 benchmark 函数，使学习者可运行 `go test` 与 `go test -bench`。
- 将本章 `README.md` 从占位内容更新为实际产出说明、自测清单和运行命令，并新增 `EXERCISES.md`。
- 不引入第三方依赖；`golangci-lint`、`dlv`、`pprof` 以文档化命令和可断言摘要呈现，不强制本章测试依赖额外工具安装。

## Capabilities

### New Capabilities
- `engineering-tutorial`: 覆盖第 07 章 Go 工程化基础学习单元的可运行代码、测试、benchmark、文档和练习。

### Modified Capabilities
- `learning-curriculum`: 将第 07 章从未落地占位章节更新为已落地章节，允许其目录包含源码、测试、benchmark 与练习材料。

## Impact

- 主要影响目录：`stage-1-syntax/07-engineering/`。
- 可能新增 OpenSpec 规格：`openspec/changes/chapter-07-engineering/specs/engineering-tutorial/spec.md`。
- 可能修改现有规格：`openspec/changes/chapter-07-engineering/specs/learning-curriculum/spec.md`。
- 验证命令包括 `go test ./stage-1-syntax/07-engineering/...`、`go test -bench=. ./stage-1-syntax/07-engineering/...`、`go run ./stage-1-syntax/07-engineering`、`go test ./...` 和 `go build ./...`。
