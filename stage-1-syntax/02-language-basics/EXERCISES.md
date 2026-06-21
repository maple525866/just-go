# 02. 语法基础 — 练习题

> 以下练习只用本章已学概念：变量/常量、基本类型、控制流、函数与包。
> 在 `stage-1-syntax/02-language-basics/` 目录下完成并验证。

## 练习 1：理解 `var` 与 `:=`（⭐）

在 `vars/vars.go` 的 `DemoSubjects` 中，把 `mathScore := 92` 改成包级 `var mathScore int = 92`（需移到函数外），函数内改为直接使用该变量。

**验收标准**：`go build ./...` 退出码 0；运行 `go run .` 仍能打印 Math 成绩行。

## 练习 2：默写 `for` 的三种形式（⭐⭐）

在 `control/control.go` 新增导出函数 `RepeatChar(ch rune, n int) string`，分别用：
1. `for i := 0; i < n; i++`（经典三段式）
2. `for n > 0 { n-- }`（类 while）
3. `for i, r := range "abc"` 的思路处理 rune（本题只需前两种形式之一 + 字符串拼接）

**验收标准**：新增表驱动测试 `RepeatChar('x', 3) == "xxx"`；`go test ./stage-1-syntax/02-language-basics/...` 全部 PASS。

## 练习 3：用 `defer` 追加日志行（⭐⭐）

仿照 `RunReport`，在 `control` 包新增 `RunWithTimestamp(body string) string`，用 `defer` 在返回字符串末尾追加一行 `logged at ready`（无需真实时间，固定字符串即可）。

**验收标准**：测试断言返回值同时包含 `body` 与 `logged at ready`；理解为何需用**命名返回值**才能让 `defer` 修改最终返回字符串。

## 练习 4：写一个闭包计数器（⭐⭐）

在 `funcs/funcs.go` 新增 `MakeCounter() func() int`，每次调用返回递增的整数（从 1 开始），演示闭包捕获外部变量。

**验收标准**：测试中连续三次调用 `c()` 分别得到 1、2、3。

## 练习 5：解释 Go 没有 `while`（⭐）

不写代码：阅读 `control/Summarize` 与 `LetterGrade`，用一句话说明「Go 用哪种 `for` 形式替代 `while`」，并指出本仓库中哪段代码体现了这一点。

**验收标准**：能口头或写在注释里：`for condition { }` 即 while；例如 `Summarize` 里 `for i := 0; i < len(scores); i++` 前的逻辑可改写为 `for len(scores) > 0` 风格的类 while（或指出其他现有循环）。
