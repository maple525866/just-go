# 第 05 章练习：并发编程

> 建议先运行 `go test -race ./stage-1-syntax/05-concurrency/...`，确认示例代码全部通过且没有数据竞争，再开始练习。

## 练习 1：实现一个可取消的平方计算

在 `goroutine` 包中新增一个函数，接收 `context.Context` 和一组整数，为每个整数并发计算平方；如果 context 被取消，应尽快停止并返回取消状态。

**验收标准：**

- 正常输入 `[]int{2,3,4}` 时能返回 `4,9,16`。
- 传入已取消的 context 时不应启动长期阻塞的 goroutine，并返回可断言的取消结果。
- 使用表驱动 + `t.Run` 编写测试。
- 运行 `go test -race ./stage-1-syntax/05-concurrency/goroutine` 通过。

## 练习 2：实现 fan-in 合并 channel

在 `channel` 包中实现一个函数，将两个只读 channel 的值合并到一个输出切片中。要求两个输入 channel 都关闭后函数才返回。

**验收标准：**

- 两个输入分别发送 `1,2` 和 `3,4` 时，最终结果包含四个值。
- 测试不依赖固定调度顺序；可以排序后断言。
- 不出现 goroutine 泄漏。
- 运行 `go test -race ./stage-1-syntax/05-concurrency/channel` 通过。

## 练习 3：为 ScoreBoard 增加快照方法

在 `syncx.ScoreBoard` 中新增 `Snapshot() map[string]int`，返回当前分数表的拷贝。

**验收标准：**

- `Snapshot` 必须使用读锁保护 map 读取。
- 修改返回的 map 不会影响 `ScoreBoard` 内部状态。
- 使用表驱动 + `t.Run` 编写测试。
- 运行 `go test -race ./stage-1-syntax/05-concurrency/syncx` 通过。

## 练习 4：实现带 timeout 的 worker

在 `ctx` 包中新增一个 worker 函数：如果任务在 timeout 内到达则处理任务，否则返回 timeout 状态。

**验收标准：**

- 任务立即到达时返回 `processed:<job>`。
- 没有任务时在 timeout 后返回可断言的超时状态。
- 测试用较小但稳定的 timeout，不依赖精确耗时。
- 运行 `go test -race ./stage-1-syntax/05-concurrency/ctx` 通过。

## 练习 5：补充一个并发坑案例说明

在 `pitfall` 包中新增一个并发坑，例如“关闭已关闭的 channel 会 panic”或“循环变量捕获导致结果错乱”，并给出规避方式。

**验收标准：**

- `Summaries()` 至少返回 4 个风险项。
- 新风险项包含风险名称、成因说明和规避方式。
- 不提交会实际 panic、race 或 deadlock 的代码。
- 运行 `go test -race ./stage-1-syntax/05-concurrency/pitfall` 通过。
