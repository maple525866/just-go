# 🚩 Capstone 3: 微服务版博客

> 阶段：③ 架构进阶 · 阶段综合项目 ｜ 难度：⭐⭐⭐⭐⭐ ｜ 预计耗时：7 天

## 🎯 项目目标

把 capstone-2 的**单体博客拆成 3 个微服务**：

- `user-svc`：用户注册 / 登录 / 鉴权（JWT）
- `post-svc`：文章 CRUD / 标签 / 列表
- `comment-svc`：评论 / 嵌套层级 / 软删除

外加：
- 一个 API 网关聚合三个服务
- 用 OpenTelemetry 实现端到端链路追踪
- 容器化（Dockerfile + docker-compose）
- 韧性策略（限流 / 熔断 / 重试）就位
- 提供一份 `docker-compose up` 即可跑起来的全套环境

## 🧩 综合应用的章节

- [✓ 12-clean-architecture] 每个服务都按整洁架构分层
- [✓ 13-ddd-patterns] 三个服务各对应一个限界上下文
- [✓ 14-microservices] gRPC + protobuf + 服务发现 + 网关
- [✓ 15-resilience-perf] 限流 / 熔断 / 重试 / pprof

## 📋 功能清单（待 OpenSpec change 填充）

> ⚠️ 当前项目尚未实现。
>
> 请通过 `/opsx-propose capstone-3-blog-ms` 来启动本项目。

## ✅ 完成标准（落地后填充）

- [ ] 三个服务可独立启动、独立测试
- [ ] 网关聚合后端到端冒烟测试通过
- [ ] 端到端链路追踪可视化（Jaeger / Tempo）
- [ ] 韧性策略可被压测验证有效（人工注入故障）
- [ ] 阶段答辩：能讲清楚拆服务的依据、限界上下文边界、CAP 取舍
