## ADDED Requirements

### Requirement: 第 15 章 SHALL 提供自包含的韧性 Gateway 教学切片

第 15 章 `stage-3-architecture/15-resilience-perf/` SHALL 提供一个可运行、可测试且默认无需外部服务的商品详情韧性 Gateway 教学切片。该切片 MUST 覆盖 fake upstream、token bucket 限流、bulkhead 隔离、指数退避重试、`gobreaker/v2` 熔断、timeout、显式 fallback 和 pprof 示例。

#### Scenario: 默认测试无需外部依赖

- **WHEN** 开发者在仓库根目录执行 `go test ./stage-3-architecture/15-resilience-perf/... -count=1`
- **THEN** 测试 MUST 在不启动 Docker、Redis、Envoy、service mesh、vegeta、wrk、hey 或外部 upstream 服务的情况下通过

#### Scenario: Gateway 按固定策略顺序编排

- **WHEN** Gateway 处理商品详情请求
- **THEN** 策略顺序 MUST 为 limit → bulkhead → timeout → retry → breaker → upstream → fallback
- **AND** 限流、隔离、重试、熔断、超时和降级路径 MUST 有确定性测试覆盖

#### Scenario: 降级响应显式标记

- **WHEN** 上游失败且 Gateway 返回 fallback 数据
- **THEN** 响应 MUST 显式包含 `degraded` 语义
- **AND** 调用方 MUST 能区分降级数据与完整真实数据

### Requirement: 第 15 章 SHALL 使用 gobreaker/v2 封装熔断

第 15 章 SHALL 集成 `github.com/sony/gobreaker/v2` 作为 circuit breaker 依赖，并通过窄接口包装，避免 Gateway 直接依赖第三方库的完整 API。

#### Scenario: breaker 打开后拒绝调用

- **WHEN** 上游持续失败达到配置阈值
- **THEN** breaker MUST 从 closed 转为 open
- **AND** 后续调用 MUST 在进入真实 upstream 前被拒绝

#### Scenario: 依赖版本受控

- **WHEN** 开发者查看 `go.mod`
- **THEN** `github.com/sony/gobreaker/v2` MUST 使用版本 `v2.4.0`
- **AND** module go version MUST 保持 `go 1.24`

### Requirement: 第 15 章 SHALL 提供性能观察与压测指导

第 15 章 SHALL 注册 `net/http/pprof` 示例，提供 heap workload，并在中文学习材料中说明 pprof 与压测命令的使用边界。默认自动化测试 MUST 验证入口可用，但 MUST NOT 依赖机器固定性能阈值。

#### Scenario: pprof 示例可被注册

- **WHEN** demo 组合根创建 HTTP server
- **THEN** `/debug/pprof/` 相关 handler MUST 可注册
- **AND** 默认测试 MUST 不要求公网暴露 pprof

#### Scenario: 压测工具仅作为手动命令

- **WHEN** 学习者阅读 Chapter 15 README 或练习
- **THEN** vegeta、wrk、hey、k6 等压测工具 MUST 被说明为手动验证工具
- **AND** 默认 `go test ./...` MUST NOT 要求这些工具存在
