# language-basics-tutorial Specification

## Purpose
TBD - created by archiving change chapter-02-language-basics. Update Purpose after archive.
## Requirements
### Requirement: 第 02 章 SHALL 提供可运行的入口程序

第 02 章目录 `stage-1-syntax/02-language-basics/` SHALL 包含一个 `package main` 的入口文件 `main.go`，其 `main` 函数 MUST 通过调用各主题子包的导出函数来组装并输出一份"语法基础报告"，且 MUST 能通过 `go run` 成功执行。

#### Scenario: go run 跑通入口程序

- **WHEN** 学习者在 `stage-1-syntax/02-language-basics/` 目录执行 `go run .`
- **THEN** 程序 MUST 以退出码 0 结束，并在标准输出打印包含变量、控制流、函数演示内容的报告文本

#### Scenario: 入口程序演示多子包协作

- **WHEN** 阅读者查看 `main.go`
- **THEN** 该文件 MUST 通过 `import` 引入本章至少两个主题子包（`vars` / `control` / `funcs`），并调用其导出函数，而非把全部逻辑写在 `main.go` 内

### Requirement: 第 02 章 SHALL 通过 vars 子包演示变量、常量与基本类型

第 02 章 SHALL 包含子目录 `vars/`（`package vars`），其中 MUST 演示：`var` 与 `:=` 声明、零值、常量与 `iota`、基本类型（数值 / 布尔 / 字符串）及显式类型转换。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/02-language-basics/vars`。

#### Scenario: vars 子包演示 iota 常量组

- **WHEN** 阅读者查看 `vars/vars.go`
- **THEN** 该文件 MUST 使用 `const (...)` 配合 `iota` 定义一组相关常量，且 MUST 含有至少一个导出函数供外部调用

#### Scenario: vars 子包演示类型转换

- **WHEN** 学习者在测试中或入口程序输出中查看类型转换相关逻辑
- **THEN** 该子包 MUST 包含至少一处显式类型转换（如 `int` → `float64`），且转换结果 MUST 可被单元测试断言

#### Scenario: vars 子包含未导出标识符

- **WHEN** 阅读者查看 `vars` 子包
- **THEN** 该包 MUST 含有至少一个首字母小写的未导出标识符，延续第 01 章可见性演示

### Requirement: 第 02 章 SHALL 通过 control 子包演示控制流

第 02 章 SHALL 包含子目录 `control/`（`package control`），其中 MUST 演示：`if`（含初始化语句形式）、`for` 的三种形式（类 while / 经典三段式 / range——本章 range 仅用于字符串或整数序列，不引入 slice 专题）、`switch`（含无表达式形式）、`defer`。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/02-language-basics/control`。

#### Scenario: control 子包演示 for 三种形式

- **WHEN** 阅读者查看 `control/control.go`
- **THEN** 该文件 MUST 含有至少两处不同形式的 `for` 循环（如 `for condition`、`for init; cond; post`、`for range`）

#### Scenario: control 子包演示 defer

- **WHEN** 阅读者查看 `control` 子包中演示 defer 的导出函数
- **THEN** 该函数 MUST 在返回前通过 `defer` 追加内容（如汇总行），且该行为 MUST 可被单元测试断言（无需捕获 stdout）

#### Scenario: control 子包提供分级函数

- **WHEN** 学习者对某个数值调用 `control` 包的导出分级函数（如 `LetterGrade`）
- **THEN** 该函数 MUST 使用 `switch` 返回对应的等级字符串，且边界值 MUST 有测试覆盖

### Requirement: 第 02 章 SHALL 通过 funcs 子包演示函数特性

第 02 章 SHALL 包含子目录 `funcs/`（`package funcs`），其中 MUST 演示：多返回值、命名返回值、可变参数、闭包。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/02-language-basics/funcs`。

#### Scenario: funcs 子包演示多返回值

- **WHEN** 阅读者查看 `funcs/funcs.go`
- **THEN** 该文件 MUST 含有一个返回至少两个值的导出函数（如 `MinMax`），且 MUST 有表驱动测试覆盖

#### Scenario: funcs 子包演示可变参数

- **WHEN** 阅读者查看 `funcs/funcs.go`
- **THEN** 该文件 MUST 含有一个接受可变参数（`...T`）的导出函数（如 `Average`），且 MUST 有测试覆盖空参数与多参数场景

#### Scenario: funcs 子包演示闭包

- **WHEN** 阅读者查看 `funcs/funcs.go`
- **THEN** 该文件 MUST 含有一个返回函数的导出函数（如 `MakeGrader`），闭包 MUST 捕获外部变量并在多次调用间保持状态或配置

### Requirement: 第 02 章 SHALL 提供单元测试

第 02 章 SHALL 在 `vars/`、`control/`、`funcs/` 各至少包含一个 `_test.go` 测试文件，且 MUST 采用表驱动（table-driven）+ `t.Run` 子测试的写法。

#### Scenario: go test 全部通过

- **WHEN** 学习者在仓库根目录执行 `go test ./stage-1-syntax/02-language-basics/...`
- **THEN** 所有测试 MUST 通过，命令 MUST 以退出码 0 返回

#### Scenario: 各子包测试采用表驱动风格

- **WHEN** 阅读者查看本章任一 `_test.go` 文件
- **THEN** 该文件 MUST 使用形如 `[]struct{ ... }` 的用例表并配合 `t.Run` 运行子测试

### Requirement: 第 02 章 SHALL 提供练习题与产出说明

第 02 章 SHALL 包含一个 `EXERCISES.md`，列出 3~5 道由浅入深的练习题（每题含验收标准）；同时章节 `README.md` 的"📦 本章产出"段落 MUST 被替换为实际内容（示例文件清单 + 运行命令），"✅ 自测清单"MUST 列出可勾选的检查项，且 MUST NOT 再包含"待 OpenSpec change 填充"占位语。

#### Scenario: EXERCISES.md 含带验收标准的练习

- **WHEN** 阅读者打开 `stage-1-syntax/02-language-basics/EXERCISES.md`
- **THEN** 该文件 MUST 包含 3 到 5 道练习题，且每道题 MUST 给出明确的验收标准（如"运行 X 后输出 Y"或"测试通过条件"）

#### Scenario: README 产出说明不再是占位

- **WHEN** 阅读者打开 `stage-1-syntax/02-language-basics/README.md` 的"📦 本章产出"段落
- **THEN** 该段落 MUST NOT 再包含"待 OpenSpec change 填充"占位语，且 MUST 列出本章 `.go` 文件清单与运行命令（`go run .` / `go test ./...`）

#### Scenario: README 自测清单与知识点对齐

- **WHEN** 阅读者查看该 README 的"✅ 自测清单"
- **THEN** 该清单 MUST 包含与 ROADMAP 关键知识点对应的可勾选项（如 `var` vs `:=`、`for` 三种形式、defer、闭包等）

