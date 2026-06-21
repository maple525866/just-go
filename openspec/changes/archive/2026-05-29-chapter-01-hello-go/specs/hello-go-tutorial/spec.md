## ADDED Requirements

### Requirement: 第 01 章 SHALL 提供可运行的入口程序

第 01 章目录 `stage-1-syntax/01-hello-go/` SHALL 包含一个 `package main` 的入口文件 `main.go`，其 `main` 函数 MUST 通过调用子包导出函数来产生问候语输出，且 MUST 能通过 `go run` 成功执行。

#### Scenario: go run 跑通入口程序

- **WHEN** 学习者在 `stage-1-syntax/01-hello-go/` 目录执行 `go run .`
- **THEN** 程序 MUST 以退出码 0 结束，并在标准输出打印一行包含问候语的文本

#### Scenario: 入口程序演示包拆分

- **WHEN** 阅读者查看 `main.go`
- **THEN** 该文件 MUST 通过 `import` 引入本章子包，并调用其导出函数，而非把全部逻辑写在 `main.go` 内

### Requirement: 第 01 章 SHALL 通过子包演示 import 路径与可见性

第 01 章 SHALL 包含一个独立子目录 `greeting/`（`package greeting`），其中 MUST 至少有一个**导出**函数（首字母大写）供 `main` 调用，并 SHOULD 包含一个**未导出**的辅助标识符（首字母小写）以演示可见性规则。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/01-hello-go/greeting`。

#### Scenario: 子包导出函数可被外部调用

- **WHEN** `main.go` 以 `just-go/stage-1-syntax/01-hello-go/greeting` 路径导入该子包
- **THEN** 它 MUST 能调用子包的导出函数（首字母大写），且该调用 MUST 编译通过

#### Scenario: 未导出标识符不可被跨包访问

- **WHEN** 阅读者查看 `greeting` 子包
- **THEN** 该包 MUST 含有至少一个首字母小写的未导出标识符，用以演示"仅包内可见"的规则

### Requirement: 第 01 章 SHALL 提供单元测试

第 01 章 SHALL 至少包含一个 `_test.go` 测试文件，且 MUST 采用表驱动（table-driven）+ `t.Run` 子测试的写法，以便学习者第一次体验 `go test`。

#### Scenario: go test 全部通过

- **WHEN** 学习者在仓库根目录执行 `go test ./stage-1-syntax/01-hello-go/...`
- **THEN** 所有测试 MUST 通过，命令 MUST 以退出码 0 返回

#### Scenario: 测试采用表驱动风格

- **WHEN** 阅读者查看本章测试文件
- **THEN** 该文件 MUST 使用形如 `[]struct{ ... }` 的用例表并配合 `t.Run` 运行子测试

### Requirement: 第 01 章 SHALL 提供练习题与产出说明

第 01 章 SHALL 包含一个 `EXERCISES.md`，列出 3~5 道由浅入深的练习题（每题含验收标准）；同时章节 `README.md` 的"📦 本章产出"段落 MUST 被替换为实际内容（示例文件清单 + 运行命令），"✅ 自测清单"MUST 列出可勾选的检查项。

#### Scenario: EXERCISES.md 含带验收标准的练习

- **WHEN** 阅读者打开 `stage-1-syntax/01-hello-go/EXERCISES.md`
- **THEN** 该文件 MUST 包含 3 到 5 道练习题，且每道题 MUST 给出明确的验收标准（如"运行 X 后输出 Y"）

#### Scenario: README 产出说明不再是占位

- **WHEN** 阅读者打开 `stage-1-syntax/01-hello-go/README.md` 的"📦 本章产出"段落
- **THEN** 该段落 MUST NOT 再包含"待 OpenSpec change 填充"占位语，且 MUST 列出本章 `.go` 文件清单与运行命令（`go run .` / `go test ./...`）
