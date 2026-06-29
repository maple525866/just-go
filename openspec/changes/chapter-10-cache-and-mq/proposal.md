## Why

第 10 章目前仍是占位章节，尚未提供缓存与消息系统示例、测试和练习材料。完成本章可以承接第 09 章持久化层，帮助学习者掌握 Redis 缓存常见模式、缓存风险治理、分布式锁思想，以及轻量消息队列的生产/消费语义。

## What Changes

- 在 `stage-2-business/10-cache-and-mq/` 下新增可运行入口程序，输出缓存与消息学习报告。
- 新增主题子包，分别演示 Cache-Aside、Read-Through、Write-Through、缓存穿透/雪崩/击穿对策、Redis 风格分布式锁、消息发布/订阅与 ack/retry 语义。
- 使用内存实现模拟 Redis 与轻量消息 broker，保证本地与 CI 无需启动外部 Redis/NATS/Kafka 服务。
- 补充 `_test.go` 覆盖缓存命中/未命中、TTL、负缓存、互斥加载、分布式锁 token 校验、消息至少一次投递与 ack。
- 将本章 `README.md` 从占位内容更新为实际产出说明、自测清单和运行命令，并新增 `EXERCISES.md`。
- 不引入生产客户端依赖；README 提供 go-redis、NATS/Kafka 的真实替换路径。

## Capabilities

### New Capabilities
- `cache-and-mq-tutorial`: 覆盖第 10 章缓存与消息学习单元的可运行代码、测试、文档和练习。

### Modified Capabilities
- `learning-curriculum`: 将第 10 章从未落地占位章节更新为已落地章节，允许其目录包含源码、测试与练习材料。

## Impact

- 主要影响目录：`stage-2-business/10-cache-and-mq/`。
- 新增 OpenSpec 规格：`openspec/changes/chapter-10-cache-and-mq/specs/cache-and-mq-tutorial/spec.md`。
- 修改现有规格：`openspec/changes/chapter-10-cache-and-mq/specs/learning-curriculum/spec.md`。
- 验证命令包括 `go test ./stage-2-business/10-cache-and-mq/...`、`go run ./stage-2-business/10-cache-and-mq`、`go test ./...` 和 `go build ./...`。
