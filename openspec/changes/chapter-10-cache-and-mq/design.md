## Context

`stage-2-business/10-cache-and-mq/` 当前只有占位 README 与 `.gitkeep`。第 10 章目标是从数据库持久化继续扩展到高频业务工程组件：缓存与消息系统。ROADMAP 提到 Redis、NATS/Kafka、缓存模式、缓存三大问题、分布式锁和消息语义。

为保持学习仓库可在无外部服务的本地/CI 环境稳定运行，本章使用内存组件模拟 Redis 与轻量 broker 的核心语义：TTL、原子 SetNX、Cache-Aside、Read-Through、Write-Through、负缓存、singleflight 风格互斥加载、ack/retry。README 会说明如何替换为 go-redis、NATS 或 Kafka 客户端。

## Goals / Non-Goals

**Goals:**

- 提供可运行的第 10 章入口程序，输出缓存与消息学习报告。
- 用测试覆盖 Cache-Aside、Read-Through、Write-Through 三种模式。
- 演示缓存穿透、雪崩、击穿的工程对策：负缓存、TTL 抖动、互斥加载。
- 演示 Redis 风格分布式锁：token、TTL、只允许持有者释放。
- 演示消息 broker：发布、消费、ack、未 ack 重投，帮助理解至少一次语义。
- README 与练习题替换占位内容，形成完整学习单元。

**Non-Goals:**

- 不要求本地启动 Redis、NATS 或 Kafka。
- 不实现生产级分布式锁续约、Redlock 或精确一次语义。
- 不实现 Kafka 分区、consumer group rebalance、offset 持久化等复杂机制。
- 不把缓存接入前两章的 Web/DB 示例；本章保持独立可运行。

## Decisions

### 1. 使用内存 Redis 风格 store

`cachex` 子包提供带 TTL 的 key-value store、SetNX、CompareAndDelete 等 Redis 风格原语，再在其上实现缓存模式和锁。

**Rationale:** 初学者可直接看到缓存策略行为，并通过测试断言 TTL、命中、未命中和 token 校验，无需外部服务。

**Alternative considered:** 直接引入 go-redis 并要求本地 Redis。该方案更贴近生产，但会让章节测试依赖外部进程。

### 2. 缓存问题以可断言策略呈现

- 穿透：不存在数据写入短 TTL 负缓存。
- 雪崩：为 TTL 添加确定性 jitter 计算函数。
- 击穿：为同一 key 的并发 miss 使用 singleflight 风格互斥加载。

**Rationale:** 这些问题可通过函数行为和计数器测试说明本质，比只写概念更可执行。

**Alternative considered:** 通过真实压测演示。该方案更复杂且不稳定，留给后续性能章节。

### 3. 消息队列使用内存 broker

`mqdemo` 子包提供 Publish、Fetch、Ack 和 RequeueExpired，未 ack 消息会重新变为可投递，用于解释至少一次语义。

**Rationale:** 消息语义核心是状态转移和 ack，而不是特定 broker API；内存 broker 可稳定测试。

**Alternative considered:** 引入 NATS。NATS 轻量但仍需外部 server 或嵌入 server，超出本章必要范围。

## Risks / Trade-offs

- [Risk] 内存实现可能被误认为生产实现 → Mitigation：README 明确其为教学模拟，并列出替换真实客户端的步骤。
- [Risk] 分布式锁未覆盖续约 → Mitigation：练习要求学习者补充续约，并在正文标注本章只实现 token + TTL + 持有者释放。
- [Risk] 消息 broker 简化了 consumer group 语义 → Mitigation：README 区分本章的至少一次模型和 Kafka/NATS 的真实能力。
