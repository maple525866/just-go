## Why

第 03 章 `03-composite-types/` 目前只有占位 README，学习者在掌握第 02 章基础语法后无法继续深入 Go 的核心数据建模能力。作为阶段一的第三站，本章必须提供**可运行、可测试、带练习**的示例代码，覆盖数组 / 切片 / map / struct / 指针，并讲清「值语义 vs 引用语义」，让学习者从"会写短程序"过渡到"能用复合类型组织真实数据"。

## What Changes

- **新增**第 03 章可运行示例代码（位于 `stage-1-syntax/03-composite-types/`）：
  - 一个 `main` 包入口，串联各子包演示，输出一份"班级花名册报告"。
  - 四个按主题拆分的子包，分别演示：
    - `seq/`：数组与切片、`len`/`cap`、`append` 扩容、**切片共享底层数组的踩坑**。
    - `dict/`：map 的 `comma-ok` 读取、零值陷阱、并发不安全（仅注释说明）。
    - `model/`：struct 的字段标签（tag）、嵌入（组合优于继承）。
    - `ptr/`：取址 / 解引用、值接收者 vs 指针接收者、值类型 vs 引用类型的传递成本。
  - 各子包及入口的表驱动单元测试，延续第 01/02 章测试范式。
- **新增**第 03 章配套学习材料：
  - `EXERCISES.md`：3~5 道由浅入深的练习题（含验收标准）。
  - 在章节 README 中填充"📦 本章产出"（示例文件清单、运行命令）与"✅ 自测清单"。
- **更新** `docs/glossary.md`：追加本章引入的术语（数组 / 切片 / 容量 cap / map / 结构体 / 字段标签 / 嵌入 / 指针 / 值语义 vs 引用语义 / 指针接收者等）。

## Capabilities

### New Capabilities

- `composite-types-tutorial`：定义第 03 章"复合类型"作为可运行学习单元的契约——入口程序、四个主题子包（seq / dict / model / ptr）、表驱动单元测试、练习题与运行指引必须齐备，且 `go run` / `go build ./...` / `go test ./...` 均能成功。

### Modified Capabilities

（无——`learning-curriculum` 已在第 01 章放宽"已落地章节可含源码"约束，本章无需再次修改。）

## Impact

- **代码**：在 `stage-1-syntax/03-composite-types/` 下新增 `.go` 文件（入口 + 4 个子包 + 测试）。仓库继续满足 `go build ./...` 与 `go test ./...` 通过。
- **依赖**：不引入任何第三方依赖，仅使用标准库（`fmt` / `strings` / `testing` 等）。
- **文档**：更新章节 README、新增 EXERCISES.md、追加术语表条目。
- **未来提案**：本章沿用第 01/02 章确立的"章节内代码 + 子包拆分 + 测试 + 练习 + README 产出说明"范式，为后续 `chapter-04-*` 提供参照。
