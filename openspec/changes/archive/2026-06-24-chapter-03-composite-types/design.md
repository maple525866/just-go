## Context

`just-go` 约定**每一章 = 一个独立 OpenSpec change**。第 01 章 `chapter-01-hello-go`、第 02 章 `chapter-02-language-basics` 已落地，学习者能跑通程序、掌握变量 / 控制流 / 函数 / 包。第 03 章"复合类型"是阶段一的第三站，面向 **已熟悉基础语法、但还不会用复合类型组织数据** 的学习者，需系统覆盖 ROADMAP 列出的全部关键知识点（数组 / 切片 / map / struct / 指针 + 值语义 vs 引用语义）。

约束：
- 单模块（`module just-go`，`go 1.24`），无第三方依赖。
- 本章只能使用截至本章已介绍的概念 + 标准库（`fmt` / `strings` / `sort` / `testing` 等）。
- 落地后 `go build ./...` 与 `go test ./...` 必须通过。
- 延续第 01/02 章范式：子包拆分、纯函数可测、表驱动测试、EXERCISES.md + README 产出说明。

## Goals / Non-Goals

**Goals:**
- 提供一个可 `go run` 跑通的入口程序，串联输出一份"班级花名册报告"，直观展示本章各知识点。
- 通过四个主题子包（`seq` / `dict` / `model` / `ptr`）拆分示例，演示包组织与可见性。
- 各子包提供表驱动单元测试，覆盖核心纯函数逻辑。
- 明确演示两个高频踩坑：**切片共享底层数组**、**值接收者 vs 指针接收者对修改是否生效的影响**。
- 提供 3~5 道练习题（含验收标准）与填充后的 README 产出说明 / 自测清单。
- 追加 `docs/glossary.md` 本章术语。

**Non-Goals:**
- 不涉及接口、错误处理（第 04 章）、并发（第 05 章）、泛型。map 的并发不安全仅以注释点到为止，不引入 `sync` / goroutine。
- 不引入第三方库、CI 配置、Makefile。
- 不做交互式 CLI（`flag` / `cobra` 第 07 章才系统覆盖）。

## Decisions

### 决策 1：目录结构——入口 + 四个主题子包

```text
stage-1-syntax/03-composite-types/
├── README.md
├── EXERCISES.md
├── main.go              # package main，串联各子包，打印花名册报告
├── main_test.go         # 对入口可测逻辑的表驱动测试
├── seq/
│   ├── seq.go           # 数组/切片、len/cap、append 扩容、共享底层数组踩坑
│   └── seq_test.go
├── dict/
│   ├── dict.go          # map 声明、comma-ok、零值、并发不安全（注释）
│   └── dict_test.go
├── model/
│   ├── model.go         # struct、字段 tag、嵌入（组合）
│   └── model_test.go
└── ptr/
    ├── ptr.go           # 取址/解引用、值 vs 指针接收者、传递成本
    └── ptr_test.go
```

理由：按 ROADMAP 四个知识簇自然分组，每个子包职责单一，便于逐文件阅读；延续第 02 章"子包拆分 + import 路径"的范式。子包完整 import 路径形如 `just-go/stage-1-syntax/03-composite-types/seq`。

**备选**：把数组/切片/map 合并为单个 `collections/` 包。否决——map 的零值陷阱与切片的共享底层数组是两类独立踩坑，分包更利于聚焦。

### 决策 2：统一叙事——"班级花名册"串联各子包

入口程序模拟一份班级花名册：
- `seq` 包：用切片承载一组分数，演示 `append` 扩容与 `len`/`cap`，并用一个"子切片修改污染原切片"的函数演示共享底层数组。
- `dict` 包：用 `map[string]int` 承载"姓名→分数"，演示 `comma-ok` 查询缺失键、零值返回，以及如何安全统计。
- `model` 包：定义 `Student` 结构体（含 tag），通过嵌入 `Contact` 演示组合优于继承。
- `ptr` 包：提供值接收者与指针接收者两个方法，演示只有指针接收者能改动原值；并用一个函数对比值传递与指针传递。

`main()` 调用各子包导出函数，组装并 `fmt.Println` 输出报告。可测逻辑留在子包纯函数中，`main` 仅负责编排。

**备选**：每个子包独立 `main` 演示。否决——学习者需多次 `go run` 不同目录，体验碎片化。

### 决策 3：把可测逻辑与副作用（打印）分离

与第 01/02 章一致：子包函数返回 `string` / `[]int` / `(T, bool)` 等纯结果，`main` 负责打印。测试直接断言返回值，无需捕获 stdout。

### 决策 4：切片共享底层数组的踩坑用"可断言"的方式演示

在 `seq` 包提供一个函数（如 `SubSliceMutationDemo`），它返回"修改子切片后原切片也被改动"的前后快照字符串或两个切片，让测试能直接断言"共享"这一事实，而非靠肉眼观察 stdout。

### 决策 5：指针接收者用"修改是否生效"对比演示

在 `ptr` 包定义同一结构体的值接收者方法与指针接收者方法（如 `Student.WithBonusValue` 返回新值 vs `(*Student).AddBonus` 原地修改），测试断言"值接收者不改原对象、指针接收者改原对象"，把值语义 vs 引用语义讲清楚。

### 决策 6：spec 仅新增 `composite-types-tutorial`

本章不修改 `learning-curriculum`（第 01 章已放宽"已落地章节可含源码"约束）。新增能力 `composite-types-tutorial` 描述本章产出契约。

### 决策 7：表驱动测试延续既有范式

各子包测试采用 `[]struct{ name, in, want }` + `t.Run`，保持全仓测试风格一致。

## Risks / Trade-offs

- [知识点过多导致示例子包膨胀] → 每个子包控制在 2~4 个导出函数 + 少量未导出辅助标识符；README 给出"阅读顺序"指引。
- [切片共享底层数组的踩坑难以稳定复现] → 用固定 `cap` 的切片构造确定性场景（不触发扩容才共享），测试断言确定结果。
- [map 并发不安全无法在不引入并发的前提下演示] → 仅以代码注释 + README/自测清单中文字说明点到为止，不写 goroutine。
- [指针/接收者对新手抽象] → 在 `ptr.go` 顶部注释 + README 产出说明给出简短解释，EXERCISES.md 设专项练习。
- [与第 04 章接口/错误边界] → 本章查询缺失键用 `comma-ok`（返回 `bool`）而非 `error`，避免提前引入错误处理。

## Migration Plan

纯增量，无破坏性变更：新增 `.go` 文件与文档，更新 README / 术语表。回滚仅需 `git revert` 本 change 的提交即可恢复占位状态。

## Open Questions

（无——本章范围与第 01/02 章范式对齐，技术决策已足够明确，可直接进入实现。）
