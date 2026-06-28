## Context

`stage-1-syntax/05-concurrency/` 当前只有占位 README 与 `.gitkeep`。前 4 章采用一致的章节落地方式：章节根目录提供 `main.go` 串联主题子包，主题子包通过导出函数返回可断言结果，测试采用表驱动 + `t.Run`，README 与 EXERCISES 补充学习说明。

第 05 章涉及并发，示例若直接依赖真实调度顺序、长时间 sleep 或故意死锁，会导致测试不稳定。因此本章设计重点是“可运行、可断言、可安全演示”，用受控 channel、context deadline 和同步原语展示概念，用文档和返回值说明危险模式。

## Goals / Non-Goals

**Goals:**

- 提供可运行的第 05 章入口程序，输出 goroutine、channel、sync、context 和并发坑主题报告。
- 通过独立子包组织并发概念，保持示例短小且可测试。
- 所有测试必须稳定，避免依赖不可预测调度或无限阻塞。
- 使用 `go test -race` 验证本章示例没有数据竞争。
- README 与练习题替换占位内容，形成完整学习单元。

**Non-Goals:**

- 不深入 GMP 调度器内部实现，只做浅层生命周期说明。
- 不实现生产级 worker pool、任务队列或复杂 pipeline。
- 不故意写会挂住测试的 deadlock / goroutine leak 示例。
- 不引入第三方并发库。

## Decisions

### 1. 使用五个主题子包承载并发概念

- `goroutine/`：演示 goroutine 启动、WaitGroup 等待和结果收集。
- `channel/`：演示无缓冲 / 有缓冲 channel、close + range、select timeout。
- `syncx/`：演示 Mutex、RWMutex、Once 的安全用法。
- `ctx/`：演示 context 取消、超时与协作式退出。
- `pitfall/`：以安全、可断言的方式说明 data race、goroutine 泄漏、deadlock 的成因和规避方式。

**Rationale:** 与前几章多子包结构一致，也避免单个包混杂过多并发概念。

**Alternative considered:** 用一个大 package 展示所有并发 API。该方案文件少，但测试和阅读边界不清晰。

### 2. 使用确定性同步替代脆弱 sleep

示例优先用 channel 握手、WaitGroup、context cancellation 确认行为；只有 timeout 场景使用很短且受控的 deadline。

**Rationale:** 并发测试最容易因为时序不稳定变 flaky，本章必须让学习者看到可靠写法。

**Alternative considered:** 在示例中大量使用 `time.Sleep`。该方案更直观但不适合作为测试基线。

### 3. 并发坑通过“安全模拟 + 文档说明”展示

`pitfall` 包不写真正会产生 data race、泄漏或死锁的代码；改为返回风险说明和安全替代方案，并用测试断言说明完整性。

**Rationale:** 教学需要解释坑，但仓库质量门禁不能引入 race 或挂死测试。

**Alternative considered:** 添加被跳过的失败示例测试。这样会增加维护成本，也容易被误运行导致体验差。

## Risks / Trade-offs

- [Risk] 并发示例过于简化，学习者误以为生产并发很简单 → Mitigation：README 和练习明确标注示例边界，并把 worker pool 等复杂模式留到练习或后续阶段。
- [Risk] timeout 测试在慢机器上偶发失败 → Mitigation：测试只断言语义，不依赖精确耗时；deadline 留足余量。
- [Risk] `pitfall` 包不真正触发 race/deadlock，冲击力不足 → Mitigation：用清晰注释解释为什么不提交危险代码，并提供 `go test -race` 作为正确验证方式。
- [Risk] 子包数量多导致章节显得分散 → Mitigation：入口报告和 README 给出推荐阅读顺序。
