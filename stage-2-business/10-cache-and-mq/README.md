# 10. 缓存与消息

> 阶段：② 业务工程 ｜ 难度：⭐⭐⭐⭐☆ ｜ 预计耗时：2 天

## 🎯 学习目标

Redis 缓存的常见模式 + NATS / Kafka 消息系统的入门用法。

## 🧩 关键知识点

- Redis 客户端（`go-redis`）基础与连接池
- 缓存模式：Cache-Aside、Read-Through、Write-Through
- 缓存三大问题：穿透 / 雪崩 / 击穿 的工程对策
- 分布式锁（Redis 实现）
- 消息队列：NATS（轻量）或 Kafka（重量）的 Go 客户端
- 消息的至少一次 / 最多一次 / 精确一次语义

## 📦 本章产出（待 OpenSpec change 填充）

> ⚠️ 当前本章内容尚未实现。
>
> 请通过 `/opsx-propose chapter-10-cache-and-mq` 来落地本章内容。

## 🔗 前置依赖

- 第 09 章

## 📚 推荐扩展阅读

- [go-redis](https://github.com/redis/go-redis)
- [NATS Go Client](https://github.com/nats-io/nats.go)
- [Apache Kafka Go Client (segmentio/kafka-go)](https://github.com/segmentio/kafka-go)

## ✅ 自测清单（落地后填充）

- [ ] 能用 Cache-Aside 模式封装一层数据缓存
- [ ] 能讲清并演示缓存穿透 / 雪崩 / 击穿三种问题及其对策
- [ ] 能用 Redis 实现一把安全的分布式锁（含续约）
- [ ] 能用 NATS 或 Kafka 写一组生产 / 消费示例
- [ ] 能解释三种消息语义在实现上的差异
