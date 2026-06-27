## ADDED Requirements

### Requirement: 第 06 章 SHALL 提供可运行的入口程序

第 06 章目录 `stage-1-syntax/06-stdlib-essentials/` SHALL 包含一个 `package main` 的入口文件 `main.go`，其 `main` 函数 MUST 通过调用各主题子包的导出函数来组装并输出一份"标准库精要学习报告"，且 MUST 能通过 `go run` 成功执行。

#### Scenario: go run 跑通入口程序
- **WHEN** 学习者在 `stage-1-syntax/06-stdlib-essentials/` 目录执行 `go run .`
- **THEN** 程序 MUST 以退出码 0 结束，并在标准输出打印包含 fmt、io/bufio、os/exec、net/http、encoding、time、reflect 演示内容的报告文本

#### Scenario: 入口程序演示多子包协作
- **WHEN** 阅读者查看 `main.go`
- **THEN** 该文件 MUST 通过 `import` 引入本章至少四个主题子包（`format` / `stream` / `system` / `web` / `codec` / `clock` / `inspect`），并调用其导出函数，而非把全部逻辑写在 `main.go` 内

### Requirement: 第 06 章 SHALL 通过 format 子包演示 fmt

第 06 章 SHALL 包含子目录 `format/`（`package format`），其中 MUST 演示：`fmt.Sprintf` 与 `fmt.Fprintf` 的格式化输出。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/06-stdlib-essentials/format`。

#### Scenario: format 子包演示格式化字符串
- **WHEN** 学习者调用 format 子包的导出函数格式化学习进度
- **THEN** 该函数 MUST 使用 `fmt.Sprintf` 或 `fmt.Fprintf` 生成可断言字符串

### Requirement: 第 06 章 SHALL 通过 stream 子包演示 io 与 bufio

第 06 章 SHALL 包含子目录 `stream/`（`package stream`），其中 MUST 演示：`io.Reader` / `io.Writer` 抽象、`io.Copy` 流式拷贝、`bufio.Scanner` 按行读取。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/06-stdlib-essentials/stream`。

#### Scenario: stream 子包演示 io.Copy
- **WHEN** 学习者调用 stream 子包中的拷贝导出函数
- **THEN** 该函数 MUST 通过 `io.Copy` 从 reader 复制到 writer，并返回可断言结果

#### Scenario: stream 子包演示 bufio.Scanner
- **WHEN** 学习者调用 stream 子包中的按行读取导出函数
- **THEN** 该函数 MUST 使用 `bufio.Scanner` 读取多行文本，并返回行切片供测试断言

### Requirement: 第 06 章 SHALL 通过 system 子包演示 os 与 os/exec

第 06 章 SHALL 包含子目录 `system/`（`package system`），其中 MUST 演示：临时文件读写、环境变量读取、外部命令执行。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/06-stdlib-essentials/system`。

#### Scenario: system 子包演示文件读写
- **WHEN** 学习者调用 system 子包中的文件读写导出函数
- **THEN** 该函数 MUST 使用 `os` 在临时位置写入并读取内容，且不污染固定工作目录

#### Scenario: system 子包演示 os/exec
- **WHEN** 学习者调用 system 子包中的外部命令导出函数
- **THEN** 该函数 MUST 使用 `os/exec` 执行一个安全、可预期的命令，并返回输出或错误供调用方处理

### Requirement: 第 06 章 SHALL 通过 web 子包演示 net/http

第 06 章 SHALL 包含子目录 `web/`（`package web`），其中 MUST 演示：`net/http` handler、`http.Client`、`httptest` 本地测试。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/06-stdlib-essentials/web`。

#### Scenario: web 子包演示 HTTP handler
- **WHEN** 学习者调用 web 子包提供的 handler
- **THEN** 该 handler MUST 返回可断言的状态码和响应体，且测试 MUST 使用 `httptest`

#### Scenario: web 子包演示 HTTP client
- **WHEN** 学习者调用 web 子包中的 client 导出函数
- **THEN** 该函数 MUST 通过传入的 URL 请求本地测试服务，并返回响应内容供测试断言

### Requirement: 第 06 章 SHALL 通过 codec 子包演示 JSON 与 XML

第 06 章 SHALL 包含子目录 `codec/`（`package codec`），其中 MUST 演示：`encoding/json` 与 `encoding/xml` 的序列化和反序列化。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/06-stdlib-essentials/codec`。

#### Scenario: codec 子包演示 JSON round trip
- **WHEN** 学习者调用 codec 子包中的 JSON 导出函数
- **THEN** 该函数 MUST 将结构体编码为 JSON 再解码回来，并返回可断言结果

#### Scenario: codec 子包演示 XML round trip
- **WHEN** 学习者调用 codec 子包中的 XML 导出函数
- **THEN** 该函数 MUST 将结构体编码为 XML 再解码回来，并返回可断言结果

### Requirement: 第 06 章 SHALL 通过 clock 子包演示 time

第 06 章 SHALL 包含子目录 `clock/`（`package clock`），其中 MUST 演示：时间格式化、`time.Duration`、timer 或 ticker。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/06-stdlib-essentials/clock`。

#### Scenario: clock 子包演示时间格式化与 duration
- **WHEN** 学习者调用 clock 子包中的时间摘要导出函数
- **THEN** 该函数 MUST 使用 `time.Format` 和 `time.Duration` 返回可断言字符串或结构体

#### Scenario: clock 子包演示 ticker
- **WHEN** 学习者调用 clock 子包中的 ticker 导出函数
- **THEN** 该函数 MUST 使用 `time.Ticker` 产生有限次数 tick，并返回 tick 次数供测试断言

### Requirement: 第 06 章 SHALL 通过 inspect 子包演示只读 reflect

第 06 章 SHALL 包含子目录 `inspect/`（`package inspect`），其中 MUST 演示：使用 `reflect` 读取类型名、字段名和 struct tag，且 MUST NOT 修改值。该子包的完整 import 路径 MUST 为 `just-go/stage-1-syntax/06-stdlib-essentials/inspect`。

#### Scenario: inspect 子包演示读取字段与 tag
- **WHEN** 学习者调用 inspect 子包中的导出函数分析结构体
- **THEN** 该函数 MUST 使用 `reflect` 返回字段名与 tag 信息，且行为 MUST 有单元测试覆盖

#### Scenario: inspect 子包不演示反射修改
- **WHEN** 阅读者查看 `inspect/inspect.go`
- **THEN** 该文件 MUST 仅使用 reflect 读取元数据，不得包含通过反射修改字段值的教学示例

### Requirement: 第 06 章 SHALL 提供单元测试

第 06 章 SHALL 在各主题子包至少包含一个 `_test.go` 测试文件，且 MUST 采用表驱动（table-driven）+ `t.Run` 子测试的写法。

#### Scenario: go test 全部通过
- **WHEN** 学习者在仓库根目录执行 `go test ./stage-1-syntax/06-stdlib-essentials/...`
- **THEN** 所有测试 MUST 通过，命令 MUST 以退出码 0 返回

#### Scenario: HTTP 测试不依赖外网
- **WHEN** 阅读者查看 web 子包测试
- **THEN** 测试 MUST 使用 `httptest` 或等价本地测试服务，不得依赖公开互联网服务

### Requirement: 第 06 章 SHALL 提供练习题与产出说明

第 06 章 SHALL 包含一个 `EXERCISES.md`，列出 3~5 道由浅入深的练习题（每题含验收标准）；同时章节 `README.md` 的"📦 本章产出"段落 MUST 被替换为实际内容（示例文件清单 + 运行命令），"✅ 自测清单"MUST 列出可勾选的检查项，且 MUST NOT 再包含"待 OpenSpec change 填充"或"尚未实现"等占位语。

#### Scenario: EXERCISES.md 含带验收标准的练习
- **WHEN** 阅读者打开 `stage-1-syntax/06-stdlib-essentials/EXERCISES.md`
- **THEN** 该文件 MUST 包含 3 到 5 道练习题，且每道题 MUST 给出明确的验收标准（如"运行 X 后输出 Y"或"测试通过条件"）

#### Scenario: README 产出说明不再是占位
- **WHEN** 阅读者打开 `stage-1-syntax/06-stdlib-essentials/README.md` 的"📦 本章产出"段落
- **THEN** 该段落 MUST NOT 再包含"待 OpenSpec change 填充"或"尚未实现"占位语，且 MUST 列出本章 `.go` 文件清单与运行命令（`go run .` / `go test ./...`）

#### Scenario: README 自测清单与知识点对齐
- **WHEN** 阅读者查看该 README 的"✅ 自测清单"
- **THEN** 该清单 MUST 包含与 ROADMAP 关键知识点对应的可勾选项（如 fmt、io.Reader、bufio.Scanner、os/exec、net/http、JSON/XML、time.Ticker、reflect 等）
