## Context

`just-go` 约定**每一章 = 一个独立 OpenSpec change**。第 01 章 `chapter-01-hello-go` 已落地，学习者能跑通第一个程序并理解 `package` / `import` / `go test`。第 02 章 "语法基础" 是阶段一的第二站，面向 **L1 已写过 Hello World 但语法不熟** 的学习者，需系统覆盖 ROADMAP 列出的全部关键知识点。

约束：
- 单模块（`module just-go`，`go 1.24`），无第三方依赖。
- 本章只能使用截至本章已介绍的概念 + 标准库（`fmt` / `strings` / `testing` 等）。
- 落地后 `go build ./...` 与 `go test ./...` 必须通过。
- 延续第 01 章范式：子包拆分、纯函数可测、表驱动测试、EXERCISES.md + README 产出说明。

## Goals / Non-Goals

**Goals:**
- 提供一个可 `go run` 跑通的入口程序，串联输出一份"语法基础报告"，直观展示本章各知识点。
- 通过三个主题子包（`vars` / `control` / `funcs`）拆分示例，演示包组织与可见性。
- 各子包提供表驱动单元测试，覆盖核心纯函数逻辑。
- 提供 3~5 道练习题（含验收标准）与填充后的 README 产出说明 / 自测清单。
- 追加 `docs/glossary.md` 本章术语。

**Non-Goals:**
- 不涉及复合类型（slice / map / struct——第 03 章）、接口、错误处理、并发、泛型。
- 不引入第三方库、CI 配置、Makefile。
- 不做交互式 CLI 框架（如 `flag` / `cobra——第 07 章才系统覆盖）。

## Decisions

### 决策 1：目录结构——入口 + 三个主题子包

```text
stage-1-syntax/02-language-basics/
├── README.md
├── EXERCISES.md
├── main.go              # package main，串联各子包，打印报告
├── main_test.go         # 对入口可测逻辑的表驱动测试
├── vars/
│   ├── vars.go          # var / := / 零值 / const / iota / 类型转换
│   └── vars_test.go
├── control/
│   ├── control.go       # if / for(3种) / switch / defer
│   └── control_test.go
├── funcs/
│   ├── funcs.go         # 多返回值 / 命名返回 / 可变参数 / 闭包
│   └── funcs_test.go
└── .gitkeep
```

理由：按 ROADMAP 知识点自然分组，每个子包职责单一，便于学习者逐文件阅读；同时延续第 01 章"子包拆分 + import 路径"的范式。

**备选**：所有示例塞进单个 `main.go`。否决——无法演示包组织，且文件过长不利于学习。

### 决策 2：统一叙事——"成绩报告"串联各子包

入口程序模拟一份学生成绩报告：
- `vars` 包：定义科目常量（`iota`）、分数变量（`var` / `:=`）、零值演示、类型转换（如 `int` → `float64` 算百分比）。
- `control` 包：用 `if` / `for` / `switch` 对分数分级（A/B/C/D/F），用 `defer` 在报告结束时打印汇总行。
- `funcs` 包：提供 `MinMax`（多返回值）、`Average`（可变参数）、`MakeGrader`（闭包返回分级函数）等纯函数。

`main()` 调用各子包导出函数，组装并 `fmt.Println` 输出报告。可测逻辑留在子包纯函数中，`main` 仅负责编排。

**备选**：每个子包独立 `main` 演示。否决——学习者需多次 `go run` 不同目录，体验碎片化。

### 决策 3：把可测逻辑与副作用（打印）分离

与第 01 章一致：子包函数返回 `string` 或 `(T, error)` 等纯结果，`main` 负责打印。测试直接断言返回值，无需捕获 stdout。

`defer` 演示例外：在 `control` 包中提供一个 `RunReport(fn func() string) string` 函数，内部 `defer` 追加汇总行后返回完整报告字符串——既演示 `defer` 又保持可测。

### 决策 4：表驱动测试延续第 01 章范式

各子包测试采用 `[]struct{ name, in, want }` + `t.Run`，保持全仓测试风格一致，为第 07 章"工程化测试"埋伏笔。

### 决策 5：spec 仅新增 `language-basics-tutorial`

本章不修改 `learning-curriculum`（第 01 章已放宽"已落地章节可含源码"约束）。新增能力 `language-basics-tutorial` 描述本章产出契约。

## Risks / Trade-offs

- [知识点过多导致示例子包膨胀] → 每个子包控制在 2~3 个导出函数 + 1~2 个未导出辅助标识符；README 给出"阅读顺序"指引。
- [defer 演示难以在纯函数中体现] → 用 `RunReport` 包装函数封装 defer 逻辑，测试断言返回字符串包含 defer 追加的内容。
- [iota 与类型转换对新手抽象] → 在 `vars.go` 顶部注释 + README 产出说明中给出简短解释，EXERCISES.md 设专项练习。
- [与第 03 章 composite types 边界] → 本章仅用基本类型与字符串拼接，不使用 slice / map / struct 作为核心数据结构（最多在可变参数 `...int` 中触及 slice 底层，但不展开讲解）。

## Migration Plan

纯增量，无破坏性变更：新增 `.go` 文件与文档，更新 README / 术语表。回滚仅需 `git revert` 本 change 的提交即可恢复占位状态。

## Open Questions

（无——本章范围与第 01 章范式对齐，技术决策已足够明确，可直接进入实现。）
