# 11. 可观测性

> 阶段：② 业务工程 ｜ 难度：⭐⭐⭐⭐☆ ｜ 预计耗时：2 天

## 🎯 学习目标

让线上 Go 服务"可见"——结构化日志、metrics、trace 全套。

## 🧩 关键知识点

- 结构化日志：`slog`（标准库）
- Prometheus metrics：Counter / Gauge / Histogram
- OpenTelemetry trace：span / context 传播
- 健康检查：liveness / readiness
- 日志规范：traceID、level、字段约定

## 📦 本章产出（待 OpenSpec change 填充）

> ⚠️ 当前本章内容尚未实现。
>
> 请通过 `/opsx-propose chapter-11-observability` 来落地本章内容。

## 🔗 前置依赖

- 第 10 章

## 📚 推荐扩展阅读

- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)
- [Prometheus Go Client](https://github.com/prometheus/client_golang)
- [Go `slog` 标准库文档](https://pkg.go.dev/log/slog)

## ✅ 自测清单（落地后填充）

- [ ] 能用 `slog` 输出含 traceID 的结构化日志
- [ ] 能埋一组业务 metrics 并被 Prometheus 抓取
- [ ] 能用 OpenTelemetry 把一次请求的 trace 跨多个函数串起来
- [ ] 能解释 liveness 与 readiness 的语义差别
- [ ] 能讲清"日志 / metric / trace"三者各自的最适场景
