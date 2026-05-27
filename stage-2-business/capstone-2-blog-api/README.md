# 🚩 Capstone 2: 博客后端 API（单体）

> 阶段：② 业务工程 · 阶段综合项目 ｜ 难度：⭐⭐⭐⭐☆ ｜ 预计耗时：5 天

## 🎯 项目目标

做一个**生产级单体博客后端**：

- 用户注册 / 登录（JWT）
- 文章 CRUD + 列表分页 + 标签
- 评论：嵌套层级 / 软删除
- 全套可观测性（结构化日志 + Prometheus + OpenTelemetry）
- 完整的 README、Dockerfile、docker-compose、CI 模板

## 🧩 综合应用的章节

- [✓ 08-web-foundations] 路由 / 中间件 / 参数校验
- [✓ 09-data-persistence] MySQL + GORM + 迁移 + 事务
- [✓ 10-cache-and-mq] Redis 缓存 + 消息队列异步通知
- [✓ 11-observability] slog + Prometheus + OpenTelemetry

## 📋 功能清单（待 OpenSpec change 填充）

> ⚠️ 当前项目尚未实现。
>
> 请通过 `/opsx-propose capstone-2-blog-api` 来启动本项目。

## ✅ 完成标准（落地后填充）

- [ ] 所有列出的章节知识点至少综合使用过一次
- [ ] 代码可运行（`docker-compose up` 一键启）、有测试、有 README
- [ ] 接口文档（OpenAPI / Swagger）齐全
- [ ] 端到端冒烟测试通过
- [ ] 阶段答辩：能讲清楚分层（即使还没正式学整洁架构，也能初步表达自己的分层逻辑）
