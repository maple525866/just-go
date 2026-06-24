# composite-types-tutorial Specification

## Purpose
TBD - created by archiving change chapter-03-composite-types. Update Purpose after archive.
## Requirements
### Requirement: 第 03 章 SHALL 提供可运行的入口程序

第 03 章目录 `stage-1-syntax/03-composite-types/` SHALL 包含一个 `package main` 的入口文件 `main.go`，其 `main` 函数 MUST 通过调用各主题子包的导出函数来组装并输出一份"班级花名册报告"，且 MUST 能通过 `go run` 成功执行。

#### Scenario: go run 跑通入口程序

- **WHEN** 学习者在 `stage-1-syntax/03-composite-types/` 目录执行 `go run .`
- **THEN** 程序 MUST 以退出码 0 结束，并在标准输出打印包含切片、map、struct、指针演示内容的报告文本

#### Scenario: 入口程序演示多子包协作

- **WHEN** 阅读者查看 `main.go`
- **THEN** 该文件 MUST 通过 `import` 引入本章至少两个主题子包（`seq` / `dict` / `model` / `ptr`），并调用其导出函数，而非把全部逻辑写在 `main.go` 内

### Requirement: 第 03 章 SHALL 通过 seq 子包演示数组与切片

第 03 章 SHALL 包含子目录 `seq/`（`package seq`），其中 MUST 演示：数组与切片的区别、`len`/`cap`、`append` 扩容、以及切片共享底层数组的踩坑。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/03-composite-types/seq`。

#### Scenario: seq 子包演示 append 与 len/cap

- **WHEN** 阅读者查看 `seq/seq.go`
- **THEN** 该文件 MUST 含有至少一个导出函数，演示对切片 `append` 并观察 `len` 与 `cap` 的变化，且结果 MUST 可被单元测试断言

#### Scenario: seq 子包演示切片共享底层数组踩坑

- **WHEN** 学习者调用 seq 子包中演示共享底层数组的导出函数
- **THEN** 该函数 MUST 以可断言的方式（返回值而非仅打印）体现"修改子切片会影响原切片"，且该行为 MUST 有单元测试覆盖

#### Scenario: seq 子包含未导出标识符

- **WHEN** 阅读者查看 `seq` 子包
- **THEN** 该包 MUST 含有至少一个首字母小写的未导出标识符，延续既有可见性演示

### Requirement: 第 03 章 SHALL 通过 dict 子包演示 map

第 03 章 SHALL 包含子目录 `dict/`（`package dict`），其中 MUST 演示：map 的声明与读写、`comma-ok` 查询缺失键、零值陷阱，并以注释形式说明 map 并发不安全。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/03-composite-types/dict`。

#### Scenario: dict 子包演示 comma-ok 查询

- **WHEN** 学习者用 dict 子包的导出函数查询一个不存在的键
- **THEN** 该函数 MUST 通过 `value, ok := m[key]` 形式区分"键不存在"与"值为零值"，且两种情况 MUST 有测试覆盖

#### Scenario: dict 子包标注并发不安全

- **WHEN** 阅读者查看 `dict/dict.go`
- **THEN** 该文件 MUST 含有注释说明 map 并发写不安全（无需引入并发代码）

### Requirement: 第 03 章 SHALL 通过 model 子包演示 struct

第 03 章 SHALL 包含子目录 `model/`（`package model`），其中 MUST 演示：struct 定义、字段标签（tag）、通过嵌入实现组合。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/03-composite-types/model`。

#### Scenario: model 子包演示字段标签

- **WHEN** 阅读者查看 `model/model.go`
- **THEN** 该文件 MUST 定义至少一个含字段标签（如 `json:"..."`）的结构体

#### Scenario: model 子包演示嵌入组合

- **WHEN** 阅读者查看 `model` 子包的结构体定义
- **THEN** 该包 MUST 含有一个通过嵌入其他结构体实现字段提升（组合优于继承）的示例，且嵌入字段的可访问性 MUST 有测试或导出函数体现

### Requirement: 第 03 章 SHALL 通过 ptr 子包演示指针与值/引用语义

第 03 章 SHALL 包含子目录 `ptr/`（`package ptr`），其中 MUST 演示：取址与解引用、值接收者 vs 指针接收者对原值修改是否生效的差异、值类型 vs 引用类型的传递语义。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/03-composite-types/ptr`。

#### Scenario: ptr 子包演示值接收者不改原值

- **WHEN** 学习者调用一个值接收者方法修改字段后再读取原对象
- **THEN** 原对象 MUST 保持不变，且该行为 MUST 有单元测试断言

#### Scenario: ptr 子包演示指针接收者改原值

- **WHEN** 学习者调用一个指针接收者方法修改字段后再读取原对象
- **THEN** 原对象 MUST 被改动，且该行为 MUST 有单元测试断言

### Requirement: 第 03 章 SHALL 提供单元测试

第 03 章 SHALL 在 `seq/`、`dict/`、`model/`、`ptr/` 各至少包含一个 `_test.go` 测试文件，且 MUST 采用表驱动（table-driven）+ `t.Run` 子测试的写法。

#### Scenario: go test 全部通过

- **WHEN** 学习者在仓库根目录执行 `go test ./stage-1-syntax/03-composite-types/...`
- **THEN** 所有测试 MUST 通过，命令 MUST 以退出码 0 返回

#### Scenario: 各子包测试采用表驱动风格

- **WHEN** 阅读者查看本章任一 `_test.go` 文件
- **THEN** 该文件 MUST 使用形如 `[]struct{ ... }` 的用例表并配合 `t.Run` 运行子测试

### Requirement: 第 03 章 SHALL 提供练习题与产出说明

第 03 章 SHALL 包含一个 `EXERCISES.md`，列出 3~5 道由浅入深的练习题（每题含验收标准）；同时章节 `README.md` 的"📦 本章产出"段落 MUST 被替换为实际内容（示例文件清单 + 运行命令），"✅ 自测清单"MUST 列出可勾选的检查项，且 MUST NOT 再包含"待 OpenSpec change 填充"或"尚未实现"等占位语。

#### Scenario: EXERCISES.md 含带验收标准的练习

- **WHEN** 阅读者打开 `stage-1-syntax/03-composite-types/EXERCISES.md`
- **THEN** 该文件 MUST 包含 3 到 5 道练习题，且每道题 MUST 给出明确的验收标准（如"运行 X 后输出 Y"或"测试通过条件"）

#### Scenario: README 产出说明不再是占位

- **WHEN** 阅读者打开 `stage-1-syntax/03-composite-types/README.md` 的"📦 本章产出"段落
- **THEN** 该段落 MUST NOT 再包含"待 OpenSpec change 填充"或"尚未实现"占位语，且 MUST 列出本章 `.go` 文件清单与运行命令（`go run .` / `go test ./...`）

#### Scenario: README 自测清单与知识点对齐

- **WHEN** 阅读者查看该 README 的"✅ 自测清单"
- **THEN** 该清单 MUST 包含与 ROADMAP 关键知识点对应的可勾选项（如切片扩容、map 并发不安全、值接收者 vs 指针接收者、切片共享底层数组、struct 嵌入等）

