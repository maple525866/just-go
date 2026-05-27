# 13. DDD 战术模式

> 阶段：③ 架构进阶 ｜ 难度：⭐⭐⭐⭐⭐ ｜ 预计耗时：3 天

## 🎯 学习目标

实体 / 值对象 / 聚合根 / 领域事件 / 仓储模式的 Go 落地。

## 🧩 关键知识点

- 实体（Entity）与值对象（Value Object）的 Go 表达
- 聚合（Aggregate）与聚合根
- 仓储（Repository）模式
- 领域事件（Domain Event）与事件发布
- 应用服务（Application Service）vs 领域服务（Domain Service）

## 📦 本章产出（待 OpenSpec change 填充）

> ⚠️ 当前本章内容尚未实现。
>
> 请通过 `/opsx-propose chapter-13-ddd-patterns` 来落地本章内容。

## 🔗 前置依赖

- 第 12 章

## 📚 推荐扩展阅读

- 《领域驱动设计：软件核心复杂性应对之道》Eric Evans
- 《实现领域驱动设计》Vaughn Vernon
- [DDD with Go (Three Dots Labs)](https://threedots.tech/post/ddd-lite-in-go-introduction/)

## ✅ 自测清单（落地后填充）

- [ ] 能讲清实体 vs 值对象的判定标准
- [ ] 能用 Go 写一个不变性受保护的聚合根
- [ ] 能用仓储接口隔离领域层与持久化层
- [ ] 能用领域事件触发一次跨聚合的协作
- [ ] 能讲清应用服务和领域服务的边界
