## Why

第 04 章目前仍是占位章节，尚未提供可运行示例、单元测试和练习材料。完成本章可以承接第 03 章的复合类型基础，帮助学习者掌握 Go 中最核心的抽象机制、错误处理习惯与泛型入门。

## What Changes

- 在 `stage-1-syntax/04-interface-error/` 下新增可运行入口程序，串联接口、错误和泛型示例。
- 新增多个主题子包，分别演示 interface 隐式实现、`any` / 类型断言 / type switch、小接口设计、自定义错误包装与 `errors.Is` / `errors.As`、泛型类型参数与约束。
- 为各主题子包补充表驱动单元测试，确保示例行为可断言。
- 将本章 `README.md` 从占位内容更新为实际产出说明、自测清单和运行命令，并新增 `EXERCISES.md`。
- 不引入第三方依赖，不修改公共模块名，不影响已完成章节的运行方式。

## Capabilities

### New Capabilities
- `interface-error-tutorial`: 覆盖第 04 章接口、错误处理与泛型学习单元的可运行代码、测试、文档和练习。

### Modified Capabilities
- `learning-curriculum`: 将第 04 章从未落地占位章节更新为已落地章节，允许其目录包含源码、测试与练习材料。

## Impact

- 主要影响目录：`stage-1-syntax/04-interface-error/`。
- 可能新增 OpenSpec 规格：`openspec/changes/chapter-04-interface-error/specs/interface-error-tutorial/spec.md`。
- 可能修改现有规格：`openspec/changes/chapter-04-interface-error/specs/learning-curriculum/spec.md`。
- 验证命令包括 `go test ./stage-1-syntax/04-interface-error/...`、`go run ./stage-1-syntax/04-interface-error`、`go test ./...` 和 `go build ./...`。
