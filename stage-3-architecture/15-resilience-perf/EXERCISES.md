# 第 15 章练习：韧性与性能

所有练习都应保持默认测试无需外部服务。需要 Redis、Envoy、vegeta、wrk 或 k6 的练习必须使用独立命令或 build tag，不能破坏 `go test ./...`。

## 练习 1：替换为 `x/time/rate`（基础）

### 目标

用 `golang.org/x/time/rate` 替换教学版 token bucket，并保持 Gateway 行为不变。

### 约束

- Gateway 仍必须在调用下游前拒绝超额请求；
- 测试不得使用脆弱 `time.Sleep`；
- README 说明手写 limiter 与库实现的 API 差异。

### 验收标准

- 限流成功、耗尽和恢复测试通过；
- `429` 响应仍包含可解释的重试等待时间；
- `go test ./stage-3-architecture/15-resilience-perf/... -count=1` 通过。

## 练习 2：分布式限流设计（进阶）

### 目标

设计两个 Gateway 实例共享同一配额的分布式限流方案。

### 约束

- 明确限流维度：全局、租户、用户或 API key；
- 明确算法：GCRA、sliding window、token bucket 或 fixed window；
- 明确 Redis/存储不可用时 fail-open 还是 fail-closed；
- 不把 Redis 集成测试放入默认测试路径。

### 验收标准

- 文档说明一致性、延迟、可用性和误差权衡；
- 单元测试覆盖时钟偏差、存储超时和突发流量；
- 集成验证命令能启动依赖、执行请求并清理数据。

## 练习 3：熔断策略调优（进阶）

### 目标

把 breaker 配置改为可热更新，并比较不同失败阈值和半开探测数量。

### 约束

- 已存在请求不得被配置更新中断；
- 配置非法时保留旧策略；
- 半开探测必须限制并发。

### 验收标准

- 测试覆盖 closed → open → half-open → closed；
- 测试覆盖非法配置不生效；
- README 记录至少两组配置在压测下的现象差异。

## 练习 4：缓存降级（进阶）

### 目标

把静态 fallback 改为“最后一次成功响应缓存”。

### 约束

- 缓存必须有 TTL；
- 降级响应仍必须标记 `degraded`；
- 不能缓存 4xx 错误响应。

### 验收标准

- 上游成功后缓存写入；
- 上游失败且缓存未过期时返回降级缓存；
- 缓存过期后返回稳定错误；
- 并发读写无 race。

## 练习 5：pprof 定位热点（基础）

### 目标

给 demo 增加一个可触发 CPU 热点的 endpoint，并用 pprof 定位。

### 约束

- endpoint 默认只用于本地；
- 测试只验证 handler 可执行，不断言固定耗时；
- 文档记录采样命令和 top 输出解读。

### 验收标准

- `go tool pprof` 能看到新增热点函数；
- README 说明 CPU profile 与 heap profile 的区别；
- 默认测试仍然稳定。

## 练习 6：压测报告（挑战）

### 目标

使用 vegeta、wrk、hey 或 k6 对 Gateway 进行压测并产出报告。

### 约束

- 报告必须包含测试环境、命令、持续时间、并发或速率；
- 至少比较正常、限流、上游慢响应和熔断打开四种场景；
- 不把压测结果写成跨机器固定门禁。

### 验收标准

- 报告包含 p50/p95/p99、吞吐、错误率和降级比例；
- 能解释限流和熔断对尾延迟的影响；
- 能给出下一步容量或配置建议。

## 练习 7：迁移到网关或服务网格（挑战）

### 目标

把部分治理能力迁移到 Envoy、Kong、APISIX、Istio 或 Linkerd 的配置中。

### 约束

- 代码内策略和平台策略不能重复导致双重重试风暴；
- 明确哪些策略留在应用内，哪些交给基础设施；
- 外部依赖验证使用独立命令。

### 验收标准

- 配置能实现限流、重试或熔断中的至少两项；
- 文档比较应用内治理与基础设施治理的可观测性和发布风险；
- 回滚步骤明确。

## 完成检查

完成任一练习后至少执行：

```bash
gofmt -w stage-3-architecture/15-resilience-perf
go test ./stage-3-architecture/15-resilience-perf/... -count=1
go test -race -count=1 ./stage-3-architecture/15-resilience-perf/...
go vet ./stage-3-architecture/15-resilience-perf/...
```
