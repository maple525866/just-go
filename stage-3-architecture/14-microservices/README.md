# 14. 微服务

> 阶段：③ 架构进阶 ｜ 难度：⭐⭐⭐⭐⭐ ｜ 预计耗时：4 天

## 🎯 学习目标

gRPC + protobuf + 服务发现 + 网关 + 配置中心，全套微服务基础设施。

## 🧩 关键知识点

- gRPC 与 protobuf：IDL、代码生成、流式 RPC
- 服务发现：Consul / etcd / Nacos 任选其一
- API 网关：基础职责（路由 / 鉴权 / 限流 / 聚合）
- 配置中心：动态配置 + 灰度
- 服务间通信模式：同步（gRPC）vs 异步（MQ）

## 📦 本章产出（待 OpenSpec change 填充）

> ⚠️ 当前本章内容尚未实现。
>
> 请通过 `/opsx-propose chapter-14-microservices` 来落地本章内容。

## 🔗 前置依赖

- 第 13 章

## 📚 推荐扩展阅读

- [gRPC-Go](https://github.com/grpc/grpc-go)
- [protobuf 官方教程](https://protobuf.dev/getting-started/gotutorial/)
- 《微服务架构设计模式》Chris Richardson

## ✅ 自测清单（落地后填充）

- [ ] 能用 protobuf 定义并生成一组 gRPC 服务
- [ ] 能用 gRPC 实现单向 / 服务端流 / 双向流
- [ ] 能用一种服务发现方案让服务互相找到彼此
- [ ] 能用网关把多个服务暴露为统一外部 API
- [ ] 能讲清楚什么时候用 gRPC 同步、什么时候用 MQ 异步
