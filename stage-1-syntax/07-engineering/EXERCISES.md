# 第 07 章练习：工程化基础

> 建议先运行 `go test ./stage-1-syntax/07-engineering/...` 和 `go test -bench=. ./stage-1-syntax/07-engineering/...`，确认示例代码与 benchmark 都能跑通。

## 练习 1：为 calc 增加更多边界测试

为 `calc.Fibonacci` 补充更多边界用例，例如负数、较大的输入、连续值关系。

**验收标准：**

- 使用表驱动 + `t.Run`。
- 至少新增 3 个测试用例。
- 所有测试通过：`go test ./stage-1-syntax/07-engineering/calc`。

## 练习 2：新增一个 benchmark 并比较结果

为 `calc.Fibonacci` 新增 `BenchmarkFibonacci`，并运行 benchmark。

**验收标准：**

- benchmark 函数命名符合 `BenchmarkXxx`。
- 循环体使用 `b.N`。
- 命令 `go test -bench=Fibonacci ./stage-1-syntax/07-engineering/calc` 能输出 benchmark 结果。
- 能在笔记中解释 `ns/op` 的含义。

## 练习 3：跑通本地质量门禁

在仓库根目录依次运行 README 中的质量门禁命令。

**验收标准：**

- `go vet ./...` 通过。
- `go test -race -count=1 ./...` 通过。
- `go build ./...` 通过。
- 如果本地已安装 `golangci-lint`，运行 `golangci-lint run`；如果未安装，记录安装链接和跳过原因。

## 练习 4：用 slog 增加调试上下文

扩展 `debugx.LogExample` 或新增一个日志函数，让日志包含 request id、module name 或耗时字段。

**验收标准：**

- 测试断言日志文本包含新增字段。
- 不使用裸 `fmt.Println` 作为主要调试输出。
- `go test ./stage-1-syntax/07-engineering/debugx` 通过。

## 练习 5：生成一次 pprof CPU profile

使用本章 benchmark 生成 CPU profile，并用 `go tool pprof` 打开。

**验收标准：**

- 成功执行类似命令：`go test -bench=. -cpuprofile=cpu.out ./stage-1-syntax/07-engineering/calc`。
- 成功执行 `go tool pprof cpu.out` 并能看到 `top` 输出。
- 练习结束后删除 `cpu.out`，不要把 profile 文件提交到仓库。
