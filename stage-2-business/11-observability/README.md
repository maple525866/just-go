# 11. 可观测性

> 阶段：② 业务工程 ｜ 难度：⭐⭐⭐⭐☆ ｜ 预计耗时：2 天

## 🎯 学习目标

让线上 Go 服务"可见"——结构化日志、metrics、trace 全套，并理解健康检查在发布与运维中的作用。

## 🧩 关键知识点

- 结构化日志：标准库 `log/slog`、JSON handler、稳定字段命名
- Prometheus metrics：Counter / Gauge / Histogram 与文本暴露格式
- Trace：trace ID / span ID、父子 span、`context.Context` 传播
- 健康检查：liveness / readiness 的职责差异
- 日志规范：trace_id、level、业务字段与请求字段约定

## 📦 本章产出

本章使用轻量内存实现模拟生产可观测性组件，避免必须启动 Prometheus 或 OpenTelemetry Collector；概念和接口命名与真实系统保持一致，便于后续替换为官方 SDK。

```text
stage-2-business/11-observability/
├── main.go       # 输出可观测性学习报告
├── tracex/       # trace ID、span ID、父子 span 与 context 传播
├── loggingx/     # 基于 slog 的 trace-aware JSON 结构化日志
├── metricsx/     # Counter / Gauge / Histogram 与 Prometheus 文本格式
├── healthx/      # liveness / readiness 检查器与聚合报告
└── server/       # /livez、/readyz、/metrics、/work HTTP 示例
```

运行测试：

```bash
go test ./stage-2-business/11-observability/...
```

运行示例：

```bash
go run ./stage-2-business/11-observability
```

在 HTTP 侧使用（核心片段，省略 import 与错误处理）：

```go
checker := healthx.NewChecker()
checker.AddLiveness("process", healthx.OK)
checker.AddReadiness("database", healthx.OK)
handler := server.NewRouter(checker, metricsx.NewRegistry())
http.ListenAndServe(":8080", handler)
```

示例端点：

```bash
curl http://localhost:8080/livez
curl http://localhost:8080/readyz
curl http://localhost:8080/metrics
curl -H 'X-Trace-ID: trace-demo' http://localhost:8080/work
```

真实组件替换方向：

| 本章教学组件 | 生产替换 | 关注点 |
|---|---|---|
| `loggingx` | 标准库 `slog` + 日志平台采集器 | JSON 格式、trace_id 字段、敏感信息脱敏 |
| `metricsx` | `github.com/prometheus/client_golang/prometheus` | 指标命名、label 基数、Histogram bucket 设计 |
| `tracex` | `go.opentelemetry.io/otel` | propagator、sampler、exporter、span 生命周期 |
| `healthx` | Kubernetes probes / 负载均衡健康检查 | liveness 不误杀、readiness 保护流量入口 |

## 🔗 前置依赖

- 第 10 章

## 📚 推荐扩展阅读

- [OpenTelemetry Go](https://opentelemetry.io/docs/languages/go/)
- [Prometheus Go Client](https://github.com/prometheus/client_golang)
- [Go `slog` 标准库文档](https://pkg.go.dev/log/slog)
- [Kubernetes Probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)

## ✅ 自测清单

- [ ] 能用 `slog` 输出含 `trace_id` 的结构化 JSON 日志。
- [ ] 能区分 Counter、Gauge、Histogram，并为 HTTP 请求选择合适指标类型。
- [ ] 能用 `context.Context` 把一次请求的 trace 跨多个函数串起来。
- [ ] 能解释 liveness 与 readiness 的语义差别，以及 readiness 失败时为什么应返回 503。
- [ ] 能讲清"日志 / metric / trace"三者各自的最适场景：日志回答发生了什么，指标回答整体是否异常，trace 回答一次请求慢在哪里。
