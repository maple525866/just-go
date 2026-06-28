## Why

第 06 章目前仍是占位章节，尚未提供可运行示例、单元测试和练习材料。完成本章可以帮助学习者把前面语法、接口、并发基础落到常用标准库 API 上，具备处理文本、文件、序列化、时间、HTTP 与只读反射的基本能力。

## What Changes

- 在 `stage-1-syntax/06-stdlib-essentials/` 下新增可运行入口程序，串联 fmt、io/bufio、os/os/exec、net/http、encoding/json/xml、time、reflect 示例。
- 新增多个主题子包，分别演示格式化、流式读写与缓冲、临时文件与外部命令、HTTP handler/client 基础、JSON/XML 序列化、时间/定时器/ticker、只读反射。
- 为各主题子包补充表驱动单元测试，使用标准库测试工具（如 `httptest`、临时目录）保证行为可断言且不依赖外部网络。
- 将本章 `README.md` 从占位内容更新为实际产出说明、自测清单和运行命令，并新增 `EXERCISES.md`。
- 不引入第三方依赖，不修改公共模块名，不影响已完成章节的运行方式。

## Capabilities

### New Capabilities
- `stdlib-essentials-tutorial`: 覆盖第 06 章 Go 标准库精要学习单元的可运行代码、测试、文档和练习。

### Modified Capabilities
- `learning-curriculum`: 将第 06 章从未落地占位章节更新为已落地章节，允许其目录包含源码、测试与练习材料。

## Impact

- 主要影响目录：`stage-1-syntax/06-stdlib-essentials/`。
- 可能新增 OpenSpec 规格：`openspec/changes/chapter-06-stdlib-essentials/specs/stdlib-essentials-tutorial/spec.md`。
- 可能修改现有规格：`openspec/changes/chapter-06-stdlib-essentials/specs/learning-curriculum/spec.md`。
- 验证命令包括 `go test ./stage-1-syntax/06-stdlib-essentials/...`、`go run ./stage-1-syntax/06-stdlib-essentials`、`go test ./...` 和 `go build ./...`。
