# 第 14 章练习：微服务基础设施

所有练习都应保持默认测试无需外部服务。需要 Consul、etcd 或证书的集成测试应使用 build tag 或独立命令，不能破坏 `go test ./...`。

## 练习 1：增加 client-streaming RPC（基础）

### 目标

为库存服务增加批量调整 RPC：客户端连续发送多条调整，服务端在客户端结束发送后返回一份汇总结果。

### 约束

- 在 `.proto` 中新增独立 request/response 类型；
- 不改变现有字段号和三个 RPC 的行为；
- 任一调整非法时返回 `InvalidArgument`，并在文档中说明是否允许部分应用；
- 重新运行固定版本生成命令，不手改 `.pb.go`。

### 验收标准

- transport 测试通过真实 gRPC client stream 发送至少 3 条请求；
- 测试覆盖正常 EOF、客户端取消和非法 delta；
- `buf lint` 与 `go test ./stage-3-architecture/14-microservices/...` 通过。

### 边界

思考流中第 3 条失败时，前两条是否已提交。若要求原子性，需要在 Store 边界增加怎样的批量操作？

## 练习 2：健康感知的服务发现（进阶）

### 目标

让 Resolver 只返回健康实例，并在健康状态变化时更新 watcher 快照。

### 约束

- 健康状态属于实例元数据，不把 gRPC client 放入 registry；
- 探测循环必须接受 context 并有明确 owner；
- 不允许持锁执行网络探测；
- 健康切换要保持快照排序确定。

### 验收标准

- 测试覆盖 healthy → unhealthy → healthy；
- 所有实例不健康时 `Resolve` 返回 `ErrUnavailable`；
- 慢 watcher 不阻塞健康更新；
- `go test -race -count=1 ./stage-3-architecture/14-microservices/internal/discovery` 无 race。

### 边界

区分“进程存活”“端口可连”“服务可处理请求”三种健康语义。

## 练习 3：替换为 Consul 或 etcd 适配器（挑战）

### 目标

在不修改 Gateway 的前提下，让一个外部适配器实现现有 discovery 合同。

### 约束

- 二选一：Consul 或 etcd；
- 注册必须使用租约/TTL，并在 context 取消时撤销；
- 处理 watch 重连、压缩/游标失效或长轮询失败；
- 外部依赖测试不得进入默认测试路径。

### 验收标准

- 编译期断言适配器满足 `discovery.Registry`；
- 独立集成命令能启动依赖、注册两个实例、观察注销并清理数据；
- 模拟短暂断连后 watcher 能恢复且不重复泄漏 goroutine；
- README 记录启动、验证与清理命令。

### 边界

解释底层系统的一致性模型如何影响 `Resolve` 与 `Watch` 的可见性。

## 练习 4：为 gRPC 增加 TLS/mTLS（进阶）

### 目标

用 TLS 替换 `insecure.NewCredentials`；挑战部分要求服务端验证客户端身份。

### 约束

- 测试证书只用于测试，不提交私钥；
- server name 验证不能通过 `InsecureSkipVerify` 绕过；
- 证书加载错误必须显式返回；
- Connections 仍然拥有并关闭 ClientConn。

### 验收标准

- 合法 CA/server name 调用成功；
- 未知 CA、错误 server name 和无客户端证书分别失败；
- 错误信息不输出私钥内容；
- Gateway 集成测试使用 TLS transport 通过。

### 边界

说明加密、服务身份和业务用户身份不是同一个概念。

## 练习 5：持久化并审计动态配置（挑战）

### 目标

重启后恢复最新配置，并保留版本、操作者、时间和变更原因。

### 约束

- 更新必须先持久化成功再发布新快照；
- 版本在重启后继续单调递增；
- 写失败不能改变内存当前版本；
- secret 不得以明文进入审计日志。

### 验收标准

- 重启测试能恢复最后一次有效配置；
- 注入写失败后 `Current` 和 watcher 都保持旧版本；
- 两个并发更新不会得到相同版本；
- race test 与故障恢复测试通过。

### 边界

选择 SQLite、etcd 或其他存储，并解释事务边界与 watcher 发布顺序。

## 练习 6：设计分布式限流（挑战）

### 目标

把当前单进程 fixed-window limiter 替换为多 Gateway 实例共享的限流方案。

### 约束

- 先定义语义：全局配额、用户配额或实例配额；
- 明确算法：token bucket、sliding window 或 GCRA；
- 定义存储不可用时 fail-open 或 fail-closed；
- Gateway 仍需在调用下游前拒绝超额请求。

### 验收标准

- 两个 limiter 实例共享同一配额的测试通过；
- 时间边界、突发、时钟偏差和存储超时均有测试；
- 文档说明精确性、可用性和延迟权衡；
- 不使用 `time.Sleep` 构造脆弱测试，时间必须可注入。

### 边界

比较“强一致计数”与“本地令牌 + 周期协调”的吞吐和误差。

## 练习 7：把同步协作改为 MQ（挑战）

### 目标

选取一个不需要立即结果的流程，例如“库存变化后通知搜索索引”，用 MQ 事件替代同步调用。

### 约束

- Gateway 的商品详情查询仍使用同步 gRPC；
- 事件包含稳定 ID、发生时间、schema version 与幂等键；
- consumer 必须幂等，并定义重试与死信策略；
- 说明事务提交与事件发布之间的一致性处理。

### 验收标准

- 重复投递不会重复更新投影；
- consumer 暂时失败后可重试；
- 不兼容 schema 被明确拒绝或迁移；
- 设计文档用时序图比较修改前后的耦合与失败路径。

### 边界

解释为什么“改成 MQ”不等于自动获得 exactly-once，以及 transactional outbox 解决了哪一段一致性问题。

## 完成检查

完成任一练习后至少执行：

```bash
gofmt -w stage-3-architecture/14-microservices
go test ./stage-3-architecture/14-microservices/... -count=1
go test -race -count=1 ./stage-3-architecture/14-microservices/...
go vet ./stage-3-architecture/14-microservices/...
```

若修改 protobuf，再执行 README 中的 Buf lint 与 generate 命令，并确认生成结果与协议来源一致。
