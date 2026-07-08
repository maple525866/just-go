# 11. 可观测性练习

## 练习 1：为 `/work` 增加错误路径指标

给 `/work` 增加一个可控错误路径，例如查询参数 `?fail=true` 返回 500，并记录 `http_errors_total`。

验收标准：

- `GET /work?fail=true` 返回 500 JSON 错误。
- `/metrics` 暴露 `http_errors_total 1`。
- 测试覆盖成功路径与失败路径，且成功路径不增加错误计数。

## 练习 2：为健康检查增加外部依赖模拟

新增一个 readiness check，模拟数据库连接不可用时返回失败。

验收标准：

- 依赖正常时 `/readyz` 返回 200 和 `"ok":true`。
- 依赖失败时 `/readyz` 返回 503，响应体包含失败 check 名称和原因。
- `/livez` 不受该依赖失败影响，仍返回 200。

## 练习 3：扩展 trace 字段

把 `tracex.Span` 扩展为包含 `StartedAt` 与 `EndedAt`，并提供 `End()` 方法计算耗时。

验收标准：

- 新增单元测试证明 child span 的 `ParentID` 仍等于 root span 的 `SpanID`。
- `End()` 后可以得到大于等于 0 的耗时。
- 结构化日志中增加 `duration_ms` 字段。

## 练习 4：设计 Histogram buckets

为 HTTP 延迟设计两组 buckets：一组适合低延迟内网接口，一组适合慢任务接口，并在 README 中说明理由。

验收标准：

- 代码中可以通过参数注入不同 bucket 配置。
- 测试证明同一个观测值在不同 bucket 下会落入不同边界统计。
- 文档解释 bucket 过细和过粗的代价。
