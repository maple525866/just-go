## 1. 子包：vars

- [x] 1.1 创建 `stage-1-syntax/02-language-basics/vars/vars.go`：`package vars`，用 `const (...)` + `iota` 定义科目常量组；用 `var` / `:=` 声明分数变量；含未导出辅助标识符演示可见性
- [x] 1.2 在 `vars.go` 中实现导出函数 `FormatScore(name string, score int) string`（含 `int` → `float64` 显式类型转换算百分比）与 `ZeroValueDemo() string`（演示各基本类型零值）
- [x] 1.3 在 `vars.go` 顶部用注释标注完整 import 路径 `just-go/stage-1-syntax/02-language-basics/vars`
- [x] 1.4 创建 `stage-1-syntax/02-language-basics/vars/vars_test.go`：表驱动 + `t.Run` 覆盖 `FormatScore`（含边界分）与 `ZeroValueDemo`

## 2. 子包：control

- [x] 2.1 创建 `stage-1-syntax/02-language-basics/control/control.go`：`package control`，实现导出函数 `LetterGrade(score int) string`（用 `switch` 分级 A/B/C/D/F）
- [x] 2.2 在 `control.go` 中实现 `Summarize(scores []int) string`：用 `for` 至少两种形式（如经典三段式 + `for range`）汇总；`if` 含初始化语句形式处理空输入
- [x] 2.3 在 `control.go` 中实现 `RunReport(body string) string`：内部用 `defer` 追加汇总行后返回完整报告字符串（可测，不捕获 stdout）
- [x] 2.4 创建 `stage-1-syntax/02-language-basics/control/control_test.go`：表驱动 + `t.Run` 覆盖 `LetterGrade` 边界值、`Summarize` 空/多元素、`RunReport` 的 defer 追加行为

## 3. 子包：funcs

- [x] 3.1 创建 `stage-1-syntax/02-language-basics/funcs/funcs.go`：`package funcs`，实现 `MinMax(nums ...int) (min, max int, ok bool)`（多返回值 + 可变参数）
- [x] 3.2 在 `funcs.go` 中实现 `Average(nums ...int) float64`（命名返回值可选）与 `MakeGrader(threshold int) func(int) string`（闭包返回分级函数）
- [x] 3.3 创建 `stage-1-syntax/02-language-basics/funcs/funcs_test.go`：表驱动 + `t.Run` 覆盖 `MinMax`（空/单/多）、`Average`、`MakeGrader` 多次调用行为

## 4. 入口程序

- [x] 4.1 创建 `stage-1-syntax/02-language-basics/main.go`：`package main`，import `vars` / `control` / `funcs` 三个子包，组装并 `fmt.Println` 输出"语法基础报告"
- [x] 4.2 创建 `stage-1-syntax/02-language-basics/main_test.go`：对入口处可测的纯逻辑（如报告标题拼装）做表驱动测试（不捕获 stdout）

## 5. 学习材料

- [x] 5.1 创建 `stage-1-syntax/02-language-basics/EXERCISES.md`：3~5 道由浅入深的练习题，覆盖 `var` vs `:=`、`for` 三种形式、defer、闭包等，每题含明确验收标准
- [x] 5.2 更新 `stage-1-syntax/02-language-basics/README.md` 的"📦 本章产出"段落：移除"待 OpenSpec change 填充"占位，改列 `.go` 文件清单（main.go / vars/ / control/ / funcs/ / *_test.go）+ 运行命令（`go run .`、`go test ./...`）
- [x] 5.3 更新该 README 的"✅ 自测清单"：把检查项改为与本章实际产出对应的可勾选项（对齐 ROADMAP 关键知识点）

## 6. 跨章文档

- [x] 6.1 更新 `docs/glossary.md`：追加本章术语（零值 / iota / rune / byte / defer / 闭包 / 可变参数 / 命名返回值等），并在"出现章节"列标注 02

## 7. 验证

- [x] 7.1 在 `stage-1-syntax/02-language-basics/` 执行 `go run .`，确认退出码 0 且打印语法基础报告
- [x] 7.2 在仓库根目录执行 `go build ./...`，确认退出码 0
- [x] 7.3 在仓库根目录执行 `go test ./stage-1-syntax/02-language-basics/...`，确认全部测试通过
- [x] 7.4 执行 `go vet ./...`，确认无告警
- [x] 7.5 运行 `openspec validate chapter-02-language-basics --strict`，确认本 change 全部产出物合规
