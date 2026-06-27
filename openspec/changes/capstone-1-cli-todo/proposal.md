## Why

阶段一 01-07 章已覆盖 Go 语法、复合类型、接口错误、并发、标准库与工程化基础，但仍需要一个综合项目把这些知识串成可运行的小产品。CLI Todo 作为阶段一 capstone，可以验证学习者是否能独立组织命令行程序、设计数据模型、处理错误、持久化文件并建立测试与 benchmark 质量闭环。

## What Changes

- 在 `stage-1-syntax/capstone-1-cli-todo/` 下新增一个可运行的命令行 Todo 程序。
- 支持 `add`、`list`、`done`、`delete`、`clear` 子命令，并提供清晰的 help / usage 输出。
- 使用 JSON 文件持久化 todo 数据，包含 task id、title、done 状态、created/updated 时间。
- 使用自定义错误表达命令解析、任务不存在、存储读写失败等场景。
- 使用 goroutine + channel 实现受控的异步保存示例，确保可测试且不会泄漏 goroutine。
- 补充表驱动单元测试、benchmark、README 功能说明与 EXERCISES。
- 不引入第三方 CLI 框架，优先使用标准库实现，保持阶段一学习目标聚焦。

## Capabilities

### New Capabilities
- `cli-todo-capstone`: 覆盖阶段一 CLI Todo 综合项目的命令行交互、文件持久化、错误处理、异步保存、测试、benchmark、文档和练习。

### Modified Capabilities
- `learning-curriculum`: 将阶段一 capstone 从未落地占位项目更新为已落地项目，允许其目录包含源码、测试、benchmark 与练习材料。

## Impact

- 主要影响目录：`stage-1-syntax/capstone-1-cli-todo/`。
- 可能新增 OpenSpec 规格：`openspec/changes/capstone-1-cli-todo/specs/cli-todo-capstone/spec.md`。
- 可能修改现有规格：`openspec/changes/capstone-1-cli-todo/specs/learning-curriculum/spec.md`。
- 验证命令包括 `go test ./stage-1-syntax/capstone-1-cli-todo/...`、`go test -bench=. ./stage-1-syntax/capstone-1-cli-todo/...`、`go run ./stage-1-syntax/capstone-1-cli-todo --help`、`go test ./...` 和 `go build ./...`。
