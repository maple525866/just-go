## Why

第 01 章 `01-hello-go/` 目前只有占位 README 与 `.gitkeep`，学习者无法"跑通第一个 Go 程序"。作为整条学习路径（L0 → L4）的起点章，它必须提供一份**可运行、可测试、带练习**的最小可用内容，让零基础学习者完成"装好工具链 → 跑通程序 → 理解 `go run / build / mod / test`"的第一次闭环。

## What Changes

- **新增**第 01 章可运行示例代码（位于 `stage-1-syntax/01-hello-go/`）：
  - 一个 `main` 包的入口程序，打印问候语并演示 `package` / `import` / `func main`。
  - 一个被 `main` 调用的内部子包，演示**包的拆分、import 路径与可见性（首字母大小写）**。
  - 对应的表驱动单元测试，演示 `go test` 的基本用法。
- **新增**第 01 章配套学习材料：
  - `EXERCISES.md`：3~5 道由浅入深的练习题（含验收标准）。
  - 在章节 README 中填充"📦 本章产出"（示例文件清单、运行命令）与"✅ 自测清单"。
- **更新** `stage-1-syntax/01-hello-go/README.md`：把"待 OpenSpec change 填充"占位替换为实际产出说明与运行指引。
- **更新** `docs/glossary.md`：追加本章引入的术语（package / module / GOPATH 等）。
- **修改**约束：放宽 `learning-curriculum` 中"章节目录 MUST 仅含 README 与 .gitkeep、MUST NOT 含 .go 文件"的限制——**已落地章节**允许包含示例源码与测试（未落地章节仍保持占位）。

## Capabilities

### New Capabilities

- `hello-go-tutorial`：定义第 01 章"Hello, Go"作为可运行学习单元的契约——入口程序、子包拆分演示、单元测试、练习题与运行指引必须齐备，且 `go run` / `go build ./...` / `go test ./...` 均能成功。

### Modified Capabilities

- `learning-curriculum`：放宽"章节目录仅含 README 与 .gitkeep、禁止 .go 文件"的要求，改为"**未落地**章节保持占位、**已落地**章节可包含示例源码与测试"，以容纳本章及后续章节的实际代码产出。

## Impact

- **代码**：在 `stage-1-syntax/01-hello-go/` 下新增首批 `.go` 文件（入口 + 子包 + 测试）。仓库从"零 Go 源码"进入"有可编译包"状态，`go build ./...` 与 `go test ./...` 必须通过。
- **依赖**：不引入任何第三方依赖，仅使用标准库（`fmt` / `testing`）。
- **文档**：更新章节 README、新增 EXERCISES.md、追加术语表条目。
- **未来提案**：本章确立的"章节内代码 + 测试 + 练习 + README 产出说明"结构，将成为后续 `chapter-NN-*` 的参照范式。
