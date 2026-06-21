## Context

`just-go` 是一份"教科书式"的 Go 学习仓库，约定**每一章 = 一个独立 OpenSpec change**。奠基提案 `bootstrap-learning-repo` 已建好全仓骨架，但所有章节目录目前都只含占位 `README.md` + `.gitkeep`，无任何 `.go` 源码。

本章是第 01 章 "Hello, Go"，面向 **L0 完全零基础** 学习者，是整条路径的第一次"动手"。它的内容必须足够小（避免在第 1 章引入未学概念），又必须完整覆盖本章关键知识点：`package` / `import` / `func main`、`go run/build/test`、以及"包拆分与可见性"这一最容易被新手忽略的点。

约束：
- 仓库为单模块（`module just-go`，`go 1.24`），无第三方依赖。
- 本章只能使用截至本章已介绍的概念 + 标准库 `fmt` / `testing`。
- 落地后 `go build ./...` 与 `go test ./...` 必须通过。

## Goals / Non-Goals

**Goals:**
- 提供一个可 `go run` 跑通的入口程序，输出问候语。
- 通过一个被 `main` 调用的内部子包，直观演示**包拆分、import 路径、首字母大小写决定可见性**。
- 提供表驱动单元测试，让学习者第一次体验 `go test`。
- 提供 3~5 道练习题（含验收标准）与填充后的 README 产出说明 / 自测清单。
- 放宽 `learning-curriculum` 对"章节目录禁止 .go"的约束，使本章及后续章节的代码合法。

**Non-Goals:**
- 不涉及并发、泛型、接口、错误包装等后续章节内容。
- 不引入任何第三方库、CI 配置、Makefile。
- 不覆盖工具链安装的逐步截图（README 用文字 + 官方链接指引即可）。

## Decisions

### 决策 1：目录结构——入口在章节根，子包独立目录

```text
stage-1-syntax/01-hello-go/
├── README.md            # 填充产出说明 + 自测清单
├── EXERCISES.md         # 练习题
├── main.go              # package main，入口，调用 greeting 子包
├── main_test.go         # 对入口可测逻辑的表驱动测试
├── greeting/
│   ├── greeting.go      # package greeting，导出 Greet()，含未导出 helper
│   └── greeting_test.go # 表驱动测试 Greet()
└── .gitkeep             # 保留（无害）
```

理由：把"可复用逻辑"放进 `greeting` 子包，`main.go` 仅负责组装与输出，这样能在最小代码量下同时演示「包拆分 + import 路径 `just-go/stage-1-syntax/01-hello-go/greeting` + 大小写可见性」。

**备选**：所有代码塞进单个 `main.go`。否决——无法演示 import 路径与可见性这两个本章关键知识点。

### 决策 2：把可测逻辑与副作用（打印）分离

`main()` 只调用 `fmt.Println(greeting.Greet(name))`；真正的字符串构造逻辑放在 `greeting.Greet()`（纯函数，返回 string）。这样测试无需捕获 stdout 即可断言，符合 Go 测试惯例，也给新手示范"纯函数易测"的工程直觉。

**备选**：测试中重定向 `os.Stdout` 断言打印内容。否决——对第 1 章过于复杂。

### 决策 3：用表驱动测试作为第一个测试范式

即便逻辑简单，也采用 `[]struct{ name, in, want }` 表驱动 + `t.Run` 子测试，提前为第 07 章"工程化测试"埋下范式锚点，保持全仓测试风格一致。

### 决策 4：spec 拆分——新增 `hello-go-tutorial` + 修改 `learning-curriculum`

新增能力 `hello-go-tutorial` 描述本章产出契约（可运行 / 可测 / 有练习）。同时对 `learning-curriculum` 中"章节目录仅含 README + .gitkeep、禁 .go"的既有 Requirement 做 MODIFIED，区分"未落地（保持占位）"与"已落地（可含源码与测试）"两种状态，避免本章代码与既有 spec 冲突。

## Risks / Trade-offs

- [新手对 import 路径中的模块前缀困惑] → 在 README 与代码注释中显式标注完整 import 路径 `just-go/stage-1-syntax/01-hello-go/greeting`，并解释其由 `go.mod` 的 `module just-go` + 目录路径拼成。
- [放宽 learning-curriculum 约束后，"已落地"无客观判定标准] → 在 MODIFIED Requirement 中用"章节内存在 `.go` 文件"作为已落地的可观测判据；同时保留"仓库整体可 `go build ./...`"作为底线场景。
- [章节根目录同时含 main.go 与子目录，新手易混淆] → 通过 README 的"📦 本章产出"段落给出清晰的文件清单与各自职责说明。

## Migration Plan

纯增量，无破坏性变更：新增 `.go` 文件与文档，修改 README/术语表。回滚仅需 `git revert` 本 change 的提交即可恢复占位状态。
