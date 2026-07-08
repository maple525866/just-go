# 10. 缓存与消息

> 阶段：② 业务工程 ｜ 难度：⭐⭐⭐⭐☆ ｜ 预计耗时：2 天

## 🎯 学习目标

掌握 Redis 缓存的常见模式、缓存风险治理、分布式锁思想，以及轻量消息队列的生产 / 消费 / ack 语义。

## 🧩 关键知识点

- Redis 风格 key-value、TTL、`SET NX`、Compare-And-Delete
- 缓存模式：Cache-Aside、Read-Through、Write-Through
- 缓存三大问题：穿透 / 雪崩 / 击穿 的工程对策
- 分布式锁：token、TTL、持有者释放
- 消息队列：发布、消费、ack、未 ack 重投
- 至少一次 / 最多一次 / 精确一次语义的差异

## 📦 本章产出

本章使用内存组件模拟 Redis 与轻量消息 broker 的核心语义，保证无需启动外部 Redis / NATS / Kafka 即可运行测试。生产项目中可将这些接口替换为 go-redis、NATS 或 Kafka 客户端。

```text
stage-2-business/10-cache-and-mq/
├── main.go          # 输出缓存与消息学习报告
├── cachex/          # TTL store、缓存模式、负缓存、TTL jitter、singleflight、锁
└── mqdemo/          # 内存 broker、Publish / Fetch / Ack / RequeueExpired
```

运行测试：

```bash
go test ./stage-2-business/10-cache-and-mq/...
```

运行示例：

```bash
go run ./stage-2-business/10-cache-and-mq
```

真实组件替换方向：

| 本章教学组件 | 生产替换 | 关注点 |
|---|---|---|
| `cachex.Store` | `github.com/redis/go-redis/v9` | TTL、连接池、序列化、错误处理 |
| `cachex.LockManager` | Redis `SET key token NX PX` + Lua 删除 | token 校验、续约、时钟与网络抖动 |
| `mqdemo.Broker` | NATS / Kafka | ack、重试、offset、consumer group |

## 🔗 前置依赖

- 第 09 章

## 📚 推荐扩展阅读

- [go-redis](https://github.com/redis/go-redis)
- [NATS Go Client](https://github.com/nats-io/nats.go)
- [Apache Kafka Go Client (segmentio/kafka-go)](https://github.com/segmentio/kafka-go)

## ✅ 自测清单

- [ ] 能用 Cache-Aside 模式封装一层数据缓存。
- [ ] 能区分 Read-Through 与 Write-Through 的职责边界。
- [ ] 能讲清缓存穿透 / 雪崩 / 击穿三种问题及其对策。
- [ ] 能用 token + TTL 解释 Redis 分布式锁的安全释放条件。
- [ ] 能写一组消息生产 / 消费 / ack / 重投示例。
- [ ] 能解释至少一次、最多一次、精确一次语义在实现上的差异。
