# 08. Web 基础

> 阶段：② 业务工程 ｜ 难度：⭐⭐⭐☆☆ ｜ 预计耗时：2 天

## 🎯 学习目标

用 `net/http` 起服务，掌握路由、中间件、参数校验、JSON 编解码。

## 🧩 关键知识点

- `net/http` Server / Handler / ServeMux
- 主流 router（chi / gin / echo）选型与使用
- 中间件链：日志、Recover、CORS、限流（简易版）
- 参数校验：`go-playground/validator`
- JSON 编解码与 binding
- RESTful 风格与状态码约定

## 📦 本章产出（待 OpenSpec change 填充）

> ⚠️ 当前本章内容尚未实现。
>
> 请通过 `/opsx-propose chapter-08-web-foundations` 来落地本章内容。

## 🔗 前置依赖

- 第 07 章（阶段一全部）

## 📚 推荐扩展阅读

- [chi](https://github.com/go-chi/chi)
- [gin](https://gin-gonic.com/)
- [Building Web Apps with Go](https://gopherguides.com/articles)

## ✅ 自测清单（落地后填充）

- [ ] 能用标准库 `net/http` 起一个不依赖框架的 HTTP 服务
- [ ] 能写一条标准的中间件链（含 Recover / 日志 / CORS）
- [ ] 能用 validator 做请求体参数校验
- [ ] 能正确处理 `context.Context` 在 HTTP 请求中的传播
- [ ] 能讲清"路由 → handler → middleware → response"的完整生命周期
