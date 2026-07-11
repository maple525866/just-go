# 第 14 章微服务设计

## 背景

第 14 章目前只有占位 README。学习者已经在第 12 章掌握整洁架构，在第 13 章掌握聚合、仓储和领域事件，现在需要进一步理解进程边界、契约优先通信、服务发现、API Gateway 与动态配置。

本仓库是单一 Go module 的教材式项目。章节必须可独立运行和测试，默认不得要求 Docker、Consul、etcd 或其他外部服务，同时不能把微服务退化成仅在内存中互相调用的普通 Go 对象。

## 目标与非目标

### 目标

- 使用 protobuf 定义商品与库存服务，并提交可复现生成的 Go/gRPC 代码。
- 通过真实 gRPC 客户端与服务端演示 unary、server-streaming 和 bidirectional-streaming RPC。
- 提供可替换、并发安全的服务发现接口及内存适配器，注册和解析真实网络地址。
- 提供带版本、订阅和确定性灰度能力的动态配置中心接口及内存适配器。
- 提供 HTTP API Gateway，演示路由、Bearer 鉴权、限流、服务解析和并发聚合。
- 通过测试、可运行示例、中文说明、图示和带验收标准的练习形成完整教学闭环。

### 非目标

- 实现生产级 Consul、etcd 或 Nacos 客户端。
- 提供 Docker Compose、Kubernetes、TLS/mTLS 或跨主机部署方案。
- 实现生产级健康检查、持久化配置、全功能负载均衡或分布式限流。
- 重复实现第 10 章消息队列；本章只比较同步 gRPC 与异步 MQ 的适用边界。
- 为后续 Capstone 3 提前实现完整博客微服务。

## 方案选择

采用“真实 gRPC + 自包含基础设施适配器”的方案。商品、库存和 Gateway 之间使用真实网络地址及 gRPC 连接；服务发现和配置中心使用内存实现，以保证默认运行、测试和 CI 不依赖外部进程。

相比直接接入 Consul/etcd，该方案保留微服务的核心边界，同时把外部系统运维噪声移入进阶练习。相比只演示 gRPC，该方案覆盖路线图要求的服务发现、网关、配置中心与通信模式权衡。

## 示例领域与协议

章节使用商品与库存作为贯穿示例：

- `ProductService.GetProduct` 使用 unary RPC 返回商品名称与价格。
- `InventoryService.GetStock` 使用 unary RPC 返回指定 SKU 的库存。
- `InventoryService.WatchStock` 使用 server-streaming RPC 推送库存变化。
- `InventoryService.SyncStock` 使用 bidirectional-streaming RPC 接收库存调整并逐条返回最新库存。

protobuf 文件是协议的唯一手写来源。生成的 `.pb.go` 和 `_grpc.pb.go` 文件随仓库提交，并在 README 中记录固定工具版本和重新生成命令。生成文件不得手工修改。

## 组件与职责

### protobuf API

`api/product/v1` 和 `api/inventory/v1` 保存协议与生成代码。协议字段使用稳定编号，消息验证由服务端完成；教学文档解释兼容性原则，但不引入额外验证框架。

### 商品与库存服务

商品服务从并发安全的内存目录读取商品。库存服务维护库存并实现三种 RPC 形态。两个服务都把领域输入错误转换成标准 gRPC status，不依赖 Gateway、服务发现或配置中心的具体实现。

### 服务发现

服务发现定义注册、注销、解析和监听接口。内存适配器保存“服务名到实例元数据和网络地址”的映射，不保存服务对象或客户端对象。注册返回显式注销函数，监听受 `context.Context` 控制，快照顺序保持确定。

### 动态配置中心

配置中心保存带单调递增版本的 Gateway 配置，支持读取、更新和订阅。配置至少包含路由开关、请求超时、限流参数和灰度百分比。灰度选择基于稳定请求键计算确定性分桶，使同一请求键在配置不变时得到相同结果。

### API Gateway

Gateway 使用标准库 HTTP server，对外提供商品详情聚合端点。请求依次经过 Bearer 鉴权、限流和动态配置检查；随后通过服务发现解析商品与库存实例，建立或复用 gRPC 连接，并发请求两个下游，再组合成稳定 JSON 响应。

Gateway 只依赖服务发现和配置读取接口。限流器采用单进程实现，明确说明它不能替代生产环境的分布式限流。

### 组合根

`main.go` 在随机本地端口启动商品与库存 gRPC server，向服务发现注册地址，启动 HTTP Gateway，执行一次完整聚合调用并按 Gateway、连接、gRPC server、注册信息的所有权顺序关闭资源。

## 数据流

```text
HTTP Client
    |
    v
Bearer Auth -> Rate Limit -> Dynamic Config
    |
    v
API Gateway
    |
    +-> Service Discovery -> Product gRPC
    |
    +-> Service Discovery -> Inventory gRPC
    |
    v
Aggregated HTTP JSON
```

流式 RPC 由独立客户端示例或测试驱动，不通过 HTTP Gateway 伪装成长连接。这样可以分别讲清 gRPC 流控制与 Gateway 聚合职责。

## 错误处理

- 缺少或非法参数返回 `codes.InvalidArgument`。
- 商品或 SKU 不存在返回 `codes.NotFound`。
- 没有可用服务实例返回 `codes.Unavailable`。
- 调用取消与截止时间分别保留 `Canceled` 和 `DeadlineExceeded`。
- Gateway 将已知 gRPC 状态稳定映射为 HTTP `400`、`404`、`503` 和 `504`。
- 未知内部错误返回通用 HTTP `500`，不得向客户端暴露底层错误文本。
- Gateway 聚合采用统一超时；任一下游必要结果失败时，整体失败，不返回不完整聚合数据。
- 服务发现和配置订阅必须响应 context 取消，关闭后拒绝新操作，避免 goroutine 和资源泄漏。

## 并发与生命周期

内存商品目录、库存、服务发现和配置中心都明确由自身拥有共享状态并使用锁保护。订阅者接收不可变快照；生产者不得因慢订阅者无限阻塞。测试会覆盖取消、关闭、重复注册和并发读写，并使用 race detector 验证。

gRPC 连接由 Gateway 的连接管理组件拥有并关闭。server 的监听器、注册句柄和 goroutine 都由组合根显式管理，不使用隐式全局状态。

## 测试策略

实施严格遵循测试先行：每个新行为先编写失败测试，确认失败原因与缺失行为一致，再写最小实现使其通过。

- 协议服务测试覆盖 unary、server-streaming、bidirectional-streaming、验证和 status code。
- 服务发现测试覆盖注册、确定性解析、监听快照、注销、重复实例、取消、关闭和并发访问。
- 配置中心测试覆盖版本更新、非法配置、订阅、取消、慢订阅者和灰度边界。
- Gateway 测试覆盖鉴权、限流、动态路由、服务解析、并发聚合、超时与错误映射。
- 集成测试使用随机 TCP 端口或 `bufconn`，证明客户端通过真实 gRPC transport 与服务交互。
- `main` 测试证明完整示例能运行并有序退出。
- 文档和 OpenSpec 验证保证命令、路径、学习目标、任务与实现一致。

完成前依次运行章节测试、全仓测试、race、vet、build、golangci-lint 和 OpenSpec validation；缺失工具或环境限制必须如实记录。

## 教学材料

README 将包含：

- 契约优先与 protobuf 兼容性说明；
- 三种 RPC 形态及选择依据；
- 服务发现、动态配置和 Gateway 的职责边界；
- 同步 gRPC 与异步 MQ 的对比；
- 结构图、请求时序、运行、生成和验证命令；
- 自包含实现与生产方案之间的限制。

`EXERCISES.md` 将提供带明确验收标准的练习，包括替换 Consul/etcd、健康检查、客户端流、TLS、配置持久化和 MQ 解耦设计。

## OpenSpec 与课程同步

新增 `openspec/changes/chapter-14-microservices/`，包括 proposal、design、规格、任务清单与 change metadata。规格覆盖可运行示例、protobuf/gRPC、服务发现、配置中心、Gateway、学习材料和全仓可构建约束。

实现完成后更新第 14 章 README 的占位内容以及 `ROADMAP.md` 的章节产出与进度复选框。第 15 章和 Capstone 3 保持未实现状态。

## 风险与权衡

- 内存发现和配置可能被误解为生产方案：README 将明确接口的教学价值和外部适配器需要补齐的可靠性能力。
- 多组件可能分散学习重点：所有组件围绕一次商品详情聚合请求组织，并保持包职责单一。
- 提交生成代码增加 diff：这是为了让学习者无需本地安装生成器也能直接测试，同时保留可复现生成命令。
- 流式 RPC 容易产生泄漏：所有循环都监听 context，测试覆盖客户端取消和服务关闭。
- 单进程演示无法证明跨主机部署：真实监听地址和 gRPC transport 保留进程边界语义，部署编排明确列为后续扩展。

## 迁移顺序

1. 创建 OpenSpec change 与协议来源，生成并验证 Go/gRPC 代码。
2. 测试先行实现商品和库存 gRPC 服务。
3. 测试先行实现服务发现与动态配置中心。
4. 测试先行实现 Gateway 中间件、连接管理和聚合端点。
5. 添加组合根、集成测试、README、练习和路线图更新。
6. 完成质量门、subagent 代码审查、问题修复、GitHub issue 和 Pull Request。

## 已解决问题

- 默认运行边界：采用真实 gRPC 与内存基础设施适配器，不要求 Docker 或外部服务。
- 示例领域：采用商品与库存，避免提前实现完整博客微服务。
- MQ 范围：讲解同步与异步选择并提供练习，不重复实现消息中间件。
- Gateway 失败语义：必要下游失败时整体失败，不返回部分聚合结果。
