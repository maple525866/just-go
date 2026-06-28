## Context

`stage-1-syntax/06-stdlib-essentials/` 当前只有占位 README 与 `.gitkeep`。前几章已经形成一致落地方式：章节根目录提供 `main.go` 串联主题子包，主题子包通过导出函数返回可断言结果，测试采用表驱动 + `t.Run`，README 与 EXERCISES 补充学习说明。

第 06 章覆盖多个标准库包，容易变成 API 罗列。设计上按使用场景拆分子包，并保证所有示例可本地、可确定、无外部网络依赖地运行。

## Goals / Non-Goals

**Goals:**

- 提供可运行的第 06 章入口程序，输出标准库精要学习报告。
- 通过独立子包组织 fmt、io/bufio、os/exec、net/http、encoding、time、reflect 概念。
- 所有关键行为都通过返回值或标准库测试工具断言，而不是依赖真实外部环境。
- 使用 `httptest` 演示 HTTP 基础，避免访问外网。
- README 与练习题替换占位内容，形成完整学习单元。

**Non-Goals:**

- 不深入阶段二的 Web 框架、路由、中间件或生产 HTTP 服务治理。
- 不演示危险的 shell 拼接或依赖平台差异较大的系统命令。
- 不使用 reflect 做字段修改或动态调用，只演示读取类型/字段/tag。
- 不引入第三方依赖。

## Decisions

### 1. 使用六个主题子包覆盖标准库场景

- `format/`：演示 `fmt.Sprintf`、`fmt.Fprintf` 格式化输出。
- `stream/`：演示 `io.Reader` / `io.Writer`、`io.Copy`、`bufio.Scanner`。
- `system/`：演示 `os` 临时文件读写与 `os/exec` 安全执行简单命令。
- `web/`：演示 `net/http` handler、client 和 `httptest`。
- `codec/`：演示 `encoding/json` 与 `encoding/xml` 序列化/反序列化。
- `clock/`：演示 `time` 格式化、duration、timer/ticker。
- `inspect/`：演示 reflect 只读查看类型、字段和 tag。

**Rationale:** 按标准库常见使用场景拆分，便于学习者逐个阅读和运行测试。

**Alternative considered:** 按每个标准库包建一个子包。该方案更贴近包名，但会导致章节碎片化。

### 2. 所有 I/O 示例使用内存或临时目录

流式读写优先用 `strings.Reader` / `bytes.Buffer`；文件示例使用 `t.TempDir` 或 `os.CreateTemp`；HTTP 示例使用 `httptest.Server`。

**Rationale:** 保证测试可重复，不污染工作目录，不依赖外部网络或机器状态。

**Alternative considered:** 直接读写固定文件或请求公开 URL。该方案更接近真实场景，但会引入脆弱依赖。

### 3. os/exec 示例选择跨平台 Go 命令

外部命令示例优先执行 `go env GOVERSION`，避免依赖 Unix 专属命令如 `echo`、`cat`。

**Rationale:** 仓库本身依赖 Go 工具链，且跨平台稳定。

**Alternative considered:** 使用 shell 命令展示管道。该方案容易受平台和 shell 差异影响，并带来安全误导。

## Risks / Trade-offs

- [Risk] 标准库覆盖面太广导致每个示例偏浅 → Mitigation：本章定位“精要”，练习提供扩展任务，深入 Web 留到阶段二。
- [Risk] `os/exec` 示例在极少数环境中找不到 `go` → Mitigation：测试可在失败时报告明确错误；仓库验证本身已经依赖 Go。
- [Risk] time/ticker 测试可能因时间敏感而 flaky → Mitigation：示例使用短且有余量的定时，测试只断言语义不依赖精确时间。
- [Risk] reflect 容易被误用 → Mitigation：代码与 README 明确“只读，不滥用”，不展示修改私有字段等技巧。
