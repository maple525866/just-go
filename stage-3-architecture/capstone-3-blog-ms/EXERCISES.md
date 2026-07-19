# Capstone 3 练习

## 练习 1：替换持久化层

为三个服务分别接入独立 SQLite/PostgreSQL schema，不得共享表或跨库 join。为 repository 增加契约测试，并提供数据迁移命令。

## 练习 2：拆分 protobuf module

把单一 `blog.proto` 按 bounded context 拆为 user/post/comment 三组合同，使用 Buf lint 和 breaking check 保证兼容。

## 练习 3：接入 OpenTelemetry

用 OTel Go SDK 替换 `tracekit`，经 Collector 导出到 Jaeger。验收时必须在 UI 中看到 Gateway 到三个服务的父子 span。

## 练习 4：故障注入与压测

为 comment-svc 增加可控延迟和 `Unavailable` 注入，用 vegeta/k6 比较正常、超时、重试、breaker open 四种场景的 p95/p99、错误率和降级率。

## 练习 5：事件驱动一致性

用 outbox + 消息队列发布 `PostDeleted` 事件，让 comment-svc 异步清理或归档评论。设计重复消费、乱序、失败重放和死信处理。

## 练习 6：服务发现与滚动发布

把静态地址替换为 Consul/Kubernetes DNS，并演示一个服务双实例滚动更新。证明连接复用、健康检查和优雅退出不会造成明显请求中断。

## 练习 7：安全加固

替换教学 token 为短期 access token + refresh token；补充密钥轮换、TLS/mTLS、速率维度、审计日志和 pprof 管理端口隔离。

## 通用验收

```bash
gofmt -w stage-3-architecture/capstone-3-blog-ms
go test ./stage-3-architecture/capstone-3-blog-ms/... -count=1
go test -race ./stage-3-architecture/capstone-3-blog-ms/... -count=1
go vet ./stage-3-architecture/capstone-3-blog-ms/...
```
