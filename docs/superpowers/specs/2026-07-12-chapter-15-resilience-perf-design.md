# 第 15 章韧性与性能设计

## 背景

第 15 章是阶段三架构进阶的最后一个章节，承接第 14 章已经完成的商品/库存微服务切片。第 14 章重点解决“跨进程通信、服务发现、动态配置、Gateway 聚合”等微服务基础设施问题；第 15 章进一步回答“服务能跑之后，如何在慢响应、失败、突发流量和性能瓶颈下保持可控”。

本仓库是教学型 Go 学习路线，章节必须默认可本地运行、可测试、可复盘。第 15 章应保持 `go test ./...` 无外部服务依赖，不要求 Docker、Redis、Envoy、服务网格或压测工具预装。同时，文档必须明确区分“教学实现”和“工业生产方案”：本章手写部分主要用于锻炼机制理解、测试能力和策略编排能力，实际生产系统通常优先复用成熟库、网关、服务网格或平台能力。

## 目标与非目标

### 目标

- 在 `stage-3-architecture/15-resilience-perf/` 下提供一个自包含的商品详情韧性 Gateway 教学切片。
- 延续第 14 章商品/库存领域概念，但不直接修改第 14 章代码，避免破坏已完成章节的稳定性。
- 演示 token bucket 限流、bulkhead 隔离、指数退避重试、抖动、熔断、降级和超时控制。
- 使用 `github.com/sony/gobreaker` 展示生产常见熔断库的封装和集成。
- 使用标准库 `net/http/pprof` 展示 CPU/heap profile 的入口、风险和基本分析流程。
- 提供稳定测试、可运行示例、中文 README、练习和路线图/OpenSpec 同步。
- 在文档中列出工业级生产方案，包括 Go 库、网关/代理、服务网格、分布式限流、可观测性和压测工具。

### 非目标

- 不把本章手写 token bucket、retry 或 bulkhead 包装成生产级通用框架。
- 不引入 Redis、Envoy、Kong、APISIX、Istio、Linkerd、Sentinel 或真实服务网格作为默认运行依赖。
- 不要求 `vegeta`、`wrk`、`hey` 在默认测试中可用；它们只作为文档命令或进阶练习。
- 不重复实现第 14 章完整 gRPC 微服务基础设施，也不提前实现 Capstone 3 博客微服务。
- 不做脆弱的性能阈值测试，例如依赖固定耗时、固定 QPS 或机器相关 benchmark 结果。

## 方案选择

采用“标准库为主 + 少量生产常见库”的教学方案：

- token bucket、retry/backoff、jitter、bulkhead、fallback 和可控 fake upstream 自己实现，以便学习者读懂机制和测试边界；
- circuit breaker 使用路线图点名的 `github.com/sony/gobreaker`，但通过本章自己的 `breaker` 包封装，避免第三方类型散落到业务编排代码；
- pprof 使用标准库 `net/http/pprof`；
- 压测工具只进入 README/EXERCISES，不进入默认依赖。

该方案比“全部自己实现”更贴近真实项目中复用成熟库的习惯，也比“全部第三方库组合”更适合教学，因为核心策略仍能通过小包和测试清楚呈现。

## 目录结构

```text
15-resilience-perf/
├── internal/
│   ├── upstream/      # 可控假上游：成功、慢响应、失败、热点分配
│   ├── gateway/       # HTTP handler、策略编排、响应映射
│   ├── limiter/       # 教学版 token bucket，时间可注入
│   ├── retry/         # 指数退避 + jitter，时间/随机源可注入
│   ├── breaker/       # gobreaker 封装，隔离第三方库细节
│   ├── bulkhead/      # 并发隔离 semaphore
│   └── profiler/      # pprof 注册、热点 workload、说明入口
├── main.go
├── main_test.go
├── README.md
└── EXERCISES.md
```

## 组件职责

### upstream

`upstream` 提供可控的假商品/库存上游。它按场景返回成功、慢响应、间歇失败、持续失败、客户端错误和热点内存分配结果。测试可以通过确定性脚本驱动它，而不是依赖随机网络故障。

### limiter

`limiter` 实现教学版 token bucket。它支持容量、补充速率、突发令牌、可注入时钟和下一次可用时间计算。Gateway 在进入下游调用前执行限流，超额直接返回 `429 Too Many Requests`。

### bulkhead

`bulkhead` 使用 semaphore 表达并发隔离。请求必须先拿到槽位才能访问下游；槽满时快速失败，避免无限排队拖垮进程。实现必须保证成功、失败和 context 取消都释放槽位。

### retry

`retry` 实现指数退避、最大退避、最大尝试次数和 jitter。它只重试临时错误、上游 `5xx` 或明确可重试错误；不重试 `4xx`、参数错误和已取消 context。时间等待和随机源可注入，测试不得依赖 `time.Sleep`。

### breaker

`breaker` 封装 `github.com/sony/gobreaker`，提供本章内部使用的窄接口。它负责把连续失败转换为 closed/open/half-open 状态切换，并让 Gateway 能观察“熔断打开”这一稳定错误。业务编排代码不直接依赖 `gobreaker` 的具体类型。

### gateway

`gateway` 是策略编排中心，对外提供 `GET /api/v1/products/{sku}`。它负责限流、bulkhead、统一超时、retry、breaker、fallback 和稳定 JSON/HTTP 错误映射。策略顺序保持清晰：先拒绝本地过载，再进入隔离槽，再在统一 deadline 内重试下游调用，每次真实调用都经过熔断器，最终决定成功、降级或失败。

### profiler

`profiler` 注册 pprof handler，并提供可触发 CPU/heap 热点的演示 workload。README 说明如何用 `go tool pprof`、`go test -bench` 或压测请求观察热点，同时强调 pprof 不应默认暴露到公网。

## 请求数据流

```text
HTTP Client
    |
    v
Token Bucket Limiter
    |
    v
Bulkhead Semaphore
    |
    v
Unified Timeout Context
    |
    v
Retry Loop with Backoff + Jitter
    |
    v
Circuit Breaker
    |
    v
Fake Product/Stock Upstream
    |
    v
Fallback or Stable JSON Response
```

成功响应包含商品和库存信息。降级响应必须显式标记 `degraded: true` 和 `degrade_reason`，避免调用方把降级数据误认为完整真实数据。

示例降级响应：

```json
{
  "sku": "book-1",
  "name": "Go Microservices",
  "price_cents": 9900,
  "quantity": 0,
  "degraded": true,
  "degrade_reason": "stock_unavailable"
}
```

## 错误处理

Gateway 对外返回稳定语义：

- token bucket 超额返回 `429 Too Many Requests`，响应包含 `retry_after_millis`；
- bulkhead 并发槽满返回 `503 Service Unavailable`，响应标记 `reason: "bulkhead_full"`；
- context deadline 返回 `504 Gateway Timeout`，不继续重试；
- 熔断打开时，如果 fallback 开启，返回 `200` 且 `degraded: true`；如果 fallback 关闭，返回 `503`；
- 上游 `4xx` 不重试，映射为稳定客户端错误；
- 上游 `5xx` 或临时网络错误按指数退避 + jitter 重试，耗尽后进入 fallback 或 `503`；
- 未知内部错误返回通用 `500`，不得泄露内部错误文本。

## 并发与生命周期

limiter、bulkhead、breaker 和 fake upstream 都必须明确拥有自己的共享状态，并通过锁、channel 或原子操作保护。请求级 context 控制超时与取消。示例 server 启动后必须能有序关闭 HTTP server、后台 workload 和 pprof 入口，不留下 goroutine 泄漏。

## 测试策略

实施严格遵循测试先行。核心测试包括：

- `limiter`：令牌补充、突发、耗尽、等待时间和可注入时钟；
- `retry`：只重试临时错误/5xx，指数退避上限、jitter 范围、context 取消；
- `bulkhead`：并发槽满快速失败，释放槽位，不泄漏；
- `breaker`：closed → open → half-open → closed 的封装行为；
- `gateway`：成功、限流、隔离、重试后成功、熔断后降级、超时和错误映射；
- `profiler`：pprof handler 注册和热点 workload 可触发，不做机器相关性能断言；
- `main`：示例启动、发请求、关闭；
- 文档同步：README、EXERCISES、ROADMAP 和 OpenSpec change 一致。

完成前至少运行：

```bash
gofmt -w stage-3-architecture/15-resilience-perf
go test ./stage-3-architecture/15-resilience-perf/... -count=1
go test -race -count=1 ./stage-3-architecture/15-resilience-perf/...
go test ./... -count=1
go vet ./...
go build ./...
golangci-lint run
openspec validate chapter-15-resilience-perf --strict
```

缺失工具或环境限制必须如实记录。

## 文档与工业方案对照

README 必须明确说明：本章手写实现是为了锻炼和理解，不代表推荐生产中从零手写所有治理组件。生产中常见选择包括：

| 主题 | 教学实现 | 生产常见选择 |
|---|---|---|
| 本地限流 | 手写 token bucket | `golang.org/x/time/rate`、Envoy local rate limit、Nginx/OpenResty、APISIX/Kong 插件 |
| 分布式限流 | 练习设计 | Redis + Lua/GCRA、Envoy global rate limit、Sentinel、API Gateway 平台能力 |
| 熔断 | `gobreaker` 封装 | `sony/gobreaker`、failsafe-go、Sentinel、Envoy outlier detection、Service Mesh |
| 重试/退避 | 手写 retry/backoff | `hashicorp/go-retryablehttp`、`cenkalti/backoff`、gRPC retry policy、Envoy retry policy |
| 隔离 | semaphore bulkhead | worker pool、连接池隔离、服务网格/网关级并发限制、舱壁化部署 |
| 降级 | 静态 fallback | 缓存、只读副本、功能开关、灰度配置、SLO/error-budget 驱动策略 |
| 性能分析 | pprof 示例 | pprof、go tool trace、Prometheus/Grafana、Pyroscope、Parca、Datadog/New Relic |
| 压测 | 文档命令 | vegeta、wrk、hey、k6、JMeter、生产影子流量/回放平台 |
| 服务治理 | 代码内编排 | Envoy、Kong、APISIX、Istio、Linkerd、Consul Connect、平台网关 |

文档还要补充生产化关注点：多租户配额、身份维度、配置热更新、指标告警、熔断策略治理、灰度发布、容量评估、安全暴露控制、故障演练和回滚策略。

## OpenSpec 与课程同步

新增或完善 `openspec/changes/chapter-15-resilience-perf/`，包含 proposal、design、tasks 和 change metadata。实现完成后更新第 15 章 README 的占位内容以及 `ROADMAP.md` 的章节产出与进度复选框。Capstone 3 保持未实现状态。

## 风险与权衡

- 手写组件可能被误解为生产推荐：README 用独立章节说明教学目的，并列出工业级替代方案。
- 韧性策略过多可能分散重点：所有策略围绕一条商品详情请求组织，保持包职责单一。
- retry、timeout、breaker 交互容易产生歧义：测试固定策略顺序，并在 README 中解释为什么不是所有错误都可重试。
- 性能测试容易受机器影响：默认测试只验证行为，不断言固定性能数字；压测和 profile 进入文档和练习。
- pprof 暴露存在安全风险：示例默认只本地使用，文档强调生产环境需鉴权、内网隔离或按需开启。

## 已解决问题

- 第 15 章采用第 14 章商品/库存概念，但代码在第 15 章自包含实现。
- 核心实现采用标准库为主，熔断使用 `gobreaker`。
- 文档必须说明教学实现与生产方案的差异，并列出成熟库、网关、服务网格和压测/可观测性工具。
- 默认测试不得依赖外部服务或机器相关性能阈值。
