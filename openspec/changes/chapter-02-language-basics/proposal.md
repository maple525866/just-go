## Why

第 02 章 `02-language-basics/` 目前只有占位 README 与 `.gitkeep`，学习者在完成第 01 章后无法继续系统掌握 Go 的核心语法。作为阶段一的第二站，本章必须提供**可运行、可测试、带练习**的示例代码，覆盖变量/常量/基本类型、控制流、函数与包组织等全部基础语法，让学习者从"能跑 Hello World"过渡到"能独立写出符合 Go 风格的短程序"。

## What Changes

- **新增**第 02 章可运行示例代码（位于 `stage-1-syntax/02-language-basics/`）：
  - 一个 `main` 包入口，串联各子包演示，输出一份"语法基础报告"。
  - 三个按主题拆分的子包，分别演示：
    - `vars/`：变量声明（`var` / `:=`）、零值、常量与 `iota`、基本类型与类型转换。
    - `control/`：`if` / `for`（三种形式）/ `switch` / `defer`。
    - `funcs/`：多返回值、命名返回、可变参数、闭包。
  - 各子包及入口的表驱动单元测试，延续第 01 章测试范式。
- **新增**第 02 章配套学习材料：
  - `EXERCISES.md`：3~5 道由浅入深的练习题（含验收标准）。
  - 在章节 README 中填充"📦 本章产出"（示例文件清单、运行命令）与"✅ 自测清单"。
- **更新** `docs/glossary.md`：追加本章引入的术语（零值 / iota / rune / defer / 闭包等）。

## Capabilities

### New Capabilities

- `language-basics-tutorial`：定义第 02 章"语法基础"作为可运行学习单元的契约——入口程序、三个主题子包、单元测试、练习题与运行指引必须齐备，且 `go run` / `go build ./...` / `go test ./...` 均能成功。

### Modified Capabilities

（无——`learning-curriculum` 已在第 01 章放宽"已落地章节可含源码"约束，本章无需再次修改。）

## Impact

- **代码**：在 `stage-1-syntax/02-language-basics/` 下新增首批 `.go` 文件（入口 + 3 个子包 + 测试）。仓库继续满足 `go build ./...` 与 `go test ./...` 通过。
- **依赖**：不引入任何第三方依赖，仅使用标准库（`fmt` / `strings` / `testing` 等）。
- **文档**：更新章节 README、新增 EXERCISES.md、追加术语表条目。
- **未来提案**：本章沿用第 01 章确立的"章节内代码 + 子包拆分 + 测试 + 练习 + README 产出说明"范式，为后续 `chapter-03-*` 提供参照。
