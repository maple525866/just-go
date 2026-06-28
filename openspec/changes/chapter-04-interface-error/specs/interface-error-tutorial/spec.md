## ADDED Requirements

### Requirement: 第 04 章 SHALL 提供可运行的入口程序

第 04 章目录 `stage-1-syntax/04-interface-error/` SHALL 包含一个 `package main` 的入口文件 `main.go`，其 `main` 函数 MUST 通过调用各主题子包的导出函数来组装并输出一份"接口、错误与泛型学习报告"，且 MUST 能通过 `go run` 成功执行。

#### Scenario: go run 跑通入口程序
- **WHEN** 学习者在 `stage-1-syntax/04-interface-error/` 目录执行 `go run .`
- **THEN** 程序 MUST 以退出码 0 结束，并在标准输出打印包含 interface、error、generic 演示内容的报告文本

#### Scenario: 入口程序演示多子包协作
- **WHEN** 阅读者查看 `main.go`
- **THEN** 该文件 MUST 通过 `import` 引入本章至少两个主题子包（`iface` / `apperr` / `generic`），并调用其导出函数，而非把全部逻辑写在 `main.go` 内

### Requirement: 第 04 章 SHALL 通过 iface 子包演示接口

第 04 章 SHALL 包含子目录 `iface/`（`package iface`），其中 MUST 演示：interface 隐式实现、小接口设计、`any`、类型断言与 type switch。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/04-interface-error/iface`。

#### Scenario: iface 子包演示隐式实现
- **WHEN** 阅读者查看 `iface/iface.go`
- **THEN** 该文件 MUST 定义至少一个小接口和至少一个没有显式声明但通过方法集满足该接口的具体类型

#### Scenario: iface 子包演示接受接口返回结构体
- **WHEN** 学习者调用 iface 子包中负责格式化报告的导出函数
- **THEN** 该函数 MUST 接受接口类型参数，并返回具体结构体或字符串结果，且该行为 MUST 有单元测试覆盖

#### Scenario: iface 子包演示 any 与类型分支
- **WHEN** 学习者调用 iface 子包中分析 `any` 值的导出函数
- **THEN** 该函数 MUST 通过类型断言或 type switch 区分至少三类输入，并返回可断言的分类结果

### Requirement: 第 04 章 SHALL 通过 apperr 子包演示错误处理

第 04 章 SHALL 包含子目录 `apperr/`（`package apperr`），其中 MUST 演示：`error` 接口、sentinel error、自定义错误类型、`fmt.Errorf` + `%w` 包装、`errors.Is` 与 `errors.As`。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/04-interface-error/apperr`。

#### Scenario: apperr 子包演示 errors.Is
- **WHEN** 学习者调用 apperr 子包中返回包装错误的导出函数
- **THEN** 调用方 MUST 能使用 `errors.Is` 判断该错误链中包含指定 sentinel error，且该行为 MUST 有单元测试覆盖

#### Scenario: apperr 子包演示 errors.As
- **WHEN** 学习者调用 apperr 子包中返回自定义错误类型的导出函数
- **THEN** 调用方 MUST 能使用 `errors.As` 提取自定义错误并读取其字段，且该行为 MUST 有单元测试覆盖

#### Scenario: apperr 子包提供错误摘要
- **WHEN** 学习者调用 apperr 子包的导出摘要函数
- **THEN** 该函数 MUST 返回可展示的错误处理要点，且 MUST 包含 `%w`、`errors.Is`、`errors.As` 三个关键字中的至少两个

### Requirement: 第 04 章 SHALL 通过 generic 子包演示泛型

第 04 章 SHALL 包含子目录 `generic/`（`package generic`），其中 MUST 演示：类型参数、约束、泛型函数的基础用法。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/04-interface-error/generic`。

#### Scenario: generic 子包提供泛型 Map 与 Filter
- **WHEN** 学习者调用 generic 子包的 `Map` 或 `Filter` 导出函数
- **THEN** 这些函数 MUST 使用类型参数实现对不同元素类型切片的转换或过滤，且行为 MUST 有单元测试覆盖

#### Scenario: generic 子包演示约束
- **WHEN** 阅读者查看 `generic/generic.go`
- **THEN** 该文件 MUST 定义至少一个类型集约束，并用该约束实现一个可对数值类型工作的导出函数

### Requirement: 第 04 章 SHALL 提供单元测试

第 04 章 SHALL 在 `iface/`、`apperr/`、`generic/` 各至少包含一个 `_test.go` 测试文件，且 MUST 采用表驱动（table-driven）+ `t.Run` 子测试的写法。

#### Scenario: go test 全部通过
- **WHEN** 学习者在仓库根目录执行 `go test ./stage-1-syntax/04-interface-error/...`
- **THEN** 所有测试 MUST 通过，命令 MUST 以退出码 0 返回

#### Scenario: 各子包测试采用表驱动风格
- **WHEN** 阅读者查看本章任一 `_test.go` 文件
- **THEN** 该文件 MUST 使用形如 `[]struct{ ... }` 的用例表并配合 `t.Run` 运行子测试

### Requirement: 第 04 章 SHALL 提供练习题与产出说明

第 04 章 SHALL 包含一个 `EXERCISES.md`，列出 3~5 道由浅入深的练习题（每题含验收标准）；同时章节 `README.md` 的"📦 本章产出"段落 MUST 被替换为实际内容（示例文件清单 + 运行命令），"✅ 自测清单"MUST 列出可勾选的检查项，且 MUST NOT 再包含"待 OpenSpec change 填充"或"尚未实现"等占位语。

#### Scenario: EXERCISES.md 含带验收标准的练习
- **WHEN** 阅读者打开 `stage-1-syntax/04-interface-error/EXERCISES.md`
- **THEN** 该文件 MUST 包含 3 到 5 道练习题，且每道题 MUST 给出明确的验收标准（如"运行 X 后输出 Y"或"测试通过条件"）

#### Scenario: README 产出说明不再是占位
- **WHEN** 阅读者打开 `stage-1-syntax/04-interface-error/README.md` 的"📦 本章产出"段落
- **THEN** 该段落 MUST NOT 再包含"待 OpenSpec change 填充"或"尚未实现"占位语，且 MUST 列出本章 `.go` 文件清单与运行命令（`go run .` / `go test ./...`）

#### Scenario: README 自测清单与知识点对齐
- **WHEN** 阅读者查看该 README 的"✅ 自测清单"
- **THEN** 该清单 MUST 包含与 ROADMAP 关键知识点对应的可勾选项（如隐式实现、any、类型断言、type switch、errors.Is、errors.As、泛型约束等）
