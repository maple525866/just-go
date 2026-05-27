# 15. 韧性与性能

> 阶段：③ 架构进阶 ｜ 难度：⭐⭐⭐⭐⭐ ｜ 预计耗时：3 天

## 🎯 学习目标

从"能跑"到"扛得住"——限流 / 熔断 / 降级 / 重试，外加性能调优。

## 🧩 关键知识点

- 限流：令牌桶 / 漏桶 / 滑动窗口
- 熔断器（circuit breaker）原理与实现（`gobreaker`）
- 降级与服务隔离
- 超时与重试策略（指数退避 / 抖动）
- 性能调优：`pprof` 实战、火焰图、GC 调优
- 压测：`vegeta` / `wrk`

## 📦 本章产出（待 OpenSpec change 填充）

> ⚠️ 当前本章内容尚未实现。
>
> 请通过 `/opsx-propose chapter-15-resilience-perf` 来落地本章内容。

## 🔗 前置依赖

- 第 14 章

## 📚 推荐扩展阅读

- Netflix Hystrix 论文
- [gobreaker](https://github.com/sony/gobreaker)
- 《数据密集型应用系统设计》Martin Kleppmann
- [Go pprof tutorial](https://go.dev/blog/pprof)

## ✅ 自测清单（落地后填充）

- [ ] 能用令牌桶实现一段可配置的限流中间件
- [ ] 能用 `gobreaker` 给一段远程调用加熔断
- [ ] 能写出一段带指数退避抖动的重试函数
- [ ] 能用 `pprof` 找出一段代码的内存分配热点
- [ ] 能用 `vegeta` 或 `wrk` 压测一个接口并解读结果
