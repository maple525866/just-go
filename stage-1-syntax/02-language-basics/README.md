# 02. 语法基础

> 阶段：① 语法精通 ｜ 难度：⭐⭐☆☆☆ ｜ 预计耗时：1 天

## 🎯 学习目标

掌握变量、常量、基本类型、控制流、函数、包的全部基础语法。

## 🧩 关键知识点

- 变量声明（`var` / `:=`）、零值、常量与 `iota`
- 基本类型（数值 / 布尔 / 字符串 / `rune` / `byte`）与类型转换
- 控制流（`if` / `for` / `switch` / `defer`）
- 函数：多返回值、命名返回、可变参数、闭包
- 包的组织、可见性（首字母大小写）

## 📦 本章产出

**建议阅读顺序**：`vars/` → `control/` → `funcs/` → `main.go`

**示例代码（`.go` 文件清单）：**

| 文件 | 职责 |
| ---- | ---- |
| `main.go` | `package main` 入口，串联三子包输出「语法基础报告」；含纯函数 `reportTitle` / `buildReport` |
| `main_test.go` | `package main`，表驱动测试标题与报告结构 |
| `vars/vars.go` | `var` / `:=`、零值、`iota` 常量组、类型转换、`FormatScore` / `ZeroValueDemo` |
| `vars/vars_test.go` | 表驱动 + `t.Run` 测试 `FormatScore` 与 `ZeroValueDemo` |
| `control/control.go` | `if` 初始化语句、`for` 两种形式、`switch` 分级、`defer` 页脚（`RunReport`） |
| `control/control_test.go` | 表驱动测试 `LetterGrade` 边界、`Summarize`、`RunReport` |
| `funcs/funcs.go` | `MinMax`（多返回值 + 可变参数）、`Average`（命名返回）、`MakeGrader`（闭包） |
| `funcs/funcs_test.go` | 表驱动测试 `MinMax` / `Average` / `MakeGrader` |

子包 import 路径示例：`just-go/stage-1-syntax/02-language-basics/vars`（由 `module just-go` + 目录路径拼成）。

**学习材料：** [`EXERCISES.md`](./EXERCISES.md) —— 5 道由浅入深、含验收标准的练习题。

**运行命令：**

```bash
# 在本章目录运行入口程序
cd stage-1-syntax/02-language-basics
go run .

# 在仓库根目录运行本章测试
go test ./stage-1-syntax/02-language-basics/...
```

## 🔗 前置依赖

- 第 01 章

## 📚 推荐扩展阅读

- 《Go 程序设计语言》第 1~3 章
- [Tour of Go § Basics](https://go.dev/tour/basics/1)

## ✅ 自测清单

- [ ] 能解释 `var x int` 与 `x := 0` 的差异（包级只能用 `var`；函数内可用 `:=`）
- [ ] 能默写出 `for` 的三种形式（`for cond`、`for init; cond; post`、`for range`）
- [ ] 能讲清楚为什么 Go 没有 `while`（`for condition { }` 即 while）
- [ ] 能说明 `RunReport` 为何用命名返回值配合 `defer` 才能追加页脚
- [ ] 能在 `funcs/MakeGrader` 中看出闭包如何捕获 `threshold`
- [ ] 能在本章目录运行 `go run .` 并看到完整语法基础报告
- [ ] 能在仓库根运行 `go test ./stage-1-syntax/02-language-basics/...` 并看到全部 PASS
