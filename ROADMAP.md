# ROADMAP · 从 Go 小白到 Go 架构师

> 这是 `just-go` 仓库的**单一可信源**。所有章节的实现 / 修订都必须以本文件为对齐基线。
>
> 修订本文件须通过专门的 OpenSpec change（命名为 `revise-roadmap-*`），禁止直接编辑。

## 一、定位与使用方式

### 这份路线图是什么

一条从 **L0 完全零基础** 走到 **L4 能用 Go 搭建业务架构** 的系统化学习路径。

- 横跨三阶段：**语法精通 → 业务工程 → 架构进阶**。
- 拆解为 **15 章 + 3 个 capstone 综合项目**。
- 每一章独立可运行，每一段都有阶段答辩级别的实战项目作为锚点。

### 适用人群

| 你目前的水平 | 推荐起点 |
|---|---|
| L0 装都没装过 Go | 第 01 章 |
| L1 写过 Hello World 但语法不熟 | 第 02 章 |
| L2 基础语法会写但并发 / 泛型不熟 | 第 05 章 |
| L3 语法熟练但没做过业务工程 | 第 08 章 |
| L4 有业务经验但缺架构视角 | 第 12 章 |

### 如何使用

1. **按编号顺序**学习每一章（除非你的水平允许跳学）。
2. **每开学一章**，运行 `/opsx-propose chapter-NN-<kebab-name>` 来启动该章的 OpenSpec change。
3. **每章完成**后勾选本文末 [进度追踪](#进度追踪) 中的 checkbox。
4. **每段结束**做对应的 capstone 项目作为阶段答辩。

---

## 二、三段总览

```text
═══ 阶段一：Go 语法精通（L0 → L2） ═══
目标：能独立写出符合 Go 风格的程序；能驾驭并发；熟悉标准库；具备工程素养。
产出：7 章 + capstone-1（CLI Todo）

═══ 阶段二：Go 业务工程（L2 → L3） ═══
目标：能用 Go 独立交付一个生产级单体后端，含数据持久化、缓存、消息、可观测性。
产出：4 章 + capstone-2（博客后端 API）

═══ 阶段三：Go 架构进阶（L3 → L4） ═══
目标：从"会写业务"跃迁到"能设计架构"——分层、DDD、微服务、韧性与性能。
产出：4 章 + capstone-3（微服务版博客）
```

---

## 三、章节详表

### 阶段一：Go 语法精通

#### 01. Hello, Go ─ `stage-1-syntax/01-hello-go/`

- 🎯 **学习目标**：装好 Go 工具链，跑通第一个程序，理解 `go run / build / mod` 三件套与项目结构。
- 🧩 **关键知识点**：
  - Go 安装与环境变量（GOPATH / GOROOT / GOPROXY）
  - `go mod init / tidy / run / build / test` 命令
  - 包（package）、import 路径、`main` 函数
  - VS Code / GoLand 基础配置与调试器
- 📦 **本章产出**：_待 OpenSpec change `chapter-01-hello-go` 填充_
- 🔗 **前置依赖**：无
- ⏱️ **预计耗时**：3 小时
- 📚 **推荐扩展**：[Tour of Go](https://go.dev/tour/) §Welcome、[Effective Go](https://go.dev/doc/effective_go) §Introduction

#### 02. 语法基础 ─ `stage-1-syntax/02-language-basics/`

- 🎯 **学习目标**：掌握变量、常量、基本类型、控制流、函数、包的全部基础语法。
- 🧩 **关键知识点**：
  - 变量声明（`var` / `:=`）、零值、常量与 `iota`
  - 基本类型（数值 / 布尔 / 字符串 / rune / byte）与类型转换
  - 控制流（`if` / `for` / `switch` / `defer`）
  - 函数：多返回值、命名返回、可变参数、闭包
  - 包的组织、可见性（首字母大小写）
- 📦 **本章产出**：_待 OpenSpec change `chapter-02-language-basics` 填充_
- 🔗 **前置依赖**：01
- ⏱️ **预计耗时**：1 天
- 📚 **推荐扩展**：《Go 程序设计语言》第 1~3 章

#### 03. 复合类型 ─ `stage-1-syntax/03-composite-types/`

- 🎯 **学习目标**：玩转 array / slice / map / struct / 指针，理解值语义与引用语义。
- 🧩 **关键知识点**：
  - 数组与切片：底层数组、len/cap、扩容机制、共享底层数组的坑
  - map：声明、零值陷阱、并发不安全
  - struct：字段标签（tag）、组合（嵌入）
  - 指针：取址、解引用、何时用指针接收者
  - 值类型 vs 引用类型的传递成本
- 📦 **本章产出**：_待 OpenSpec change `chapter-03-composite-types` 填充_
- 🔗 **前置依赖**：02
- ⏱️ **预计耗时**：1 天
- 📚 **推荐扩展**：《Go 程序设计语言》第 4 章

#### 04. 接口与错误 ─ `stage-1-syntax/04-interface-error/`

- 🎯 **学习目标**：理解 interface 隐式实现，掌握标准错误处理，初识泛型。
- 🧩 **关键知识点**：
  - interface 隐式实现、空接口 `any`、类型断言、类型 switch
  - 接口设计原则（小接口、接受接口返回结构体）
  - `error` 接口、`errors.New` / `fmt.Errorf` / `%w` 包装
  - `errors.Is` / `errors.As` 错误判别
  - 泛型基础：类型参数、约束（constraints）、基本用法
- 📦 **本章产出**：_待 OpenSpec change `chapter-04-interface-error` 填充_
- 🔗 **前置依赖**：03
- ⏱️ **预计耗时**：1 天
- 📚 **推荐扩展**：《Go 程序设计语言》第 7 章 + Go 1.18 泛型官方博客

#### 05. 并发编程 ─ `stage-1-syntax/05-concurrency/`

- 🎯 **学习目标**：用 goroutine + channel + sync + context 写出正确的并发程序。
- 🧩 **关键知识点**：
  - goroutine 启动与生命周期、调度概览（GMP 浅层）
  - channel：有缓冲 / 无缓冲、关闭语义、`select`、超时模式
  - `sync` 包：`Mutex` / `RWMutex` / `WaitGroup` / `Once`
  - `context.Context`：取消、超时、值传递
  - 常见并发坑：data race、goroutine 泄漏、channel 死锁
- 📦 **本章产出**：_待 OpenSpec change `chapter-05-concurrency` 填充_
- 🔗 **前置依赖**：04
- ⏱️ **预计耗时**：2 天
- 📚 **推荐扩展**：《Go 程序设计语言》第 8~9 章、`go test -race`

#### 06. 标准库精要 ─ `stage-1-syntax/06-stdlib-essentials/`

- 🎯 **学习目标**：熟练使用 Go 标准库中最常用的 7 个包。
- 🧩 **关键知识点**：
  - `fmt`：格式化输入输出
  - `io` / `bufio`：读写抽象与缓冲
  - `os` / `os/exec`：文件与进程
  - `net` / `net/http`：网络与 HTTP（仅基础，深入留给阶段二）
  - `encoding/json` / `encoding/xml`：序列化
  - `time`：时间、定时器、ticker
  - `reflect`：反射（只学读，不滥用）
- 📦 **本章产出**：_待 OpenSpec change `chapter-06-stdlib-essentials` 填充_
- 🔗 **前置依赖**：05
- ⏱️ **预计耗时**：2 天
- 📚 **推荐扩展**：[pkg.go.dev/std](https://pkg.go.dev/std)

#### 07. 工程化基础 ─ `stage-1-syntax/07-engineering/`

- 🎯 **学习目标**：把"写代码"升级为"做工程"——模块、测试、性能、调试样样通。
- 🧩 **关键知识点**：
  - module 与 `go.work`、版本语义
  - 单元测试：`testing` 包、表驱动测试、子测试
  - benchmark 与 `go test -bench`
  - lint：`go vet` / `golangci-lint`
  - 调试：`dlv` / IDE 断点 / `log` & `slog`
  - 性能 profile：`pprof`、CPU / 内存 / 阻塞
- 📦 **本章产出**：_待 OpenSpec change `chapter-07-engineering` 填充_
- 🔗 **前置依赖**：06
- ⏱️ **预计耗时**：1.5 天
- 📚 **推荐扩展**：《Go 语言高级编程》第 2~3 章

#### 🚩 Capstone 1: CLI Todo ─ `stage-1-syntax/capstone-1-cli-todo/`

- 🎯 **项目目标**：综合阶段一所学，做一个**带文件持久化的命令行 Todo**（含子命令、单元测试、CI）。
- 🧩 **综合应用的章节**：
  - 01 / 02 / 03（基础语法 + struct + slice）
  - 04（自定义 error）
  - 05（goroutine 异步持久化）
  - 06（`os` / `encoding/json` / `time`）
  - 07（表驱动测试 + benchmark + lint）
- 📦 **本项目产出**：_待 OpenSpec change `capstone-1-cli-todo` 填充_
- 🔗 **前置依赖**：01–07
- ⏱️ **预计耗时**：3 天

---

### 阶段二：Go 业务工程

#### 08. Web 基础 ─ `stage-2-business/08-web-foundations/`

- 🎯 **学习目标**：用 `net/http` 起服务，掌握路由、中间件、参数校验、JSON 编解码。
- 🧩 **关键知识点**：
  - `net/http` Server / Handler / ServeMux
  - 主流 router（chi / gin / echo）选型与使用
  - 中间件链：日志、Recover、CORS、限流（简易版）
  - 参数校验：`go-playground/validator`
  - JSON 编解码与 binding
  - RESTful 风格与状态码约定
- 📦 **本章产出**：_待 OpenSpec change `chapter-08-web-foundations` 填充_
- 🔗 **前置依赖**：07
- ⏱️ **预计耗时**：2 天
- 📚 **推荐扩展**：[chi](https://github.com/go-chi/chi)、[gin](https://gin-gonic.com/)

#### 09. 数据持久化 ─ `stage-2-business/09-data-persistence/`

- 🎯 **学习目标**：用 `database/sql` + GORM 操作 MySQL，理解连接池、事务、迁移。
- 🧩 **关键知识点**：
  - `database/sql` 基础与连接池
  - GORM 模型定义、CRUD、关联
  - 事务（`db.Transaction`）、隔离级别
  - 数据库迁移（`goose` / `golang-migrate`）
  - 防 SQL 注入与 prepared statement
  - N+1 与预加载（`Preload` / `Joins`）
- 📦 **本章产出**：_待 OpenSpec change `chapter-09-data-persistence` 填充_
- 🔗 **前置依赖**：08
- ⏱️ **预计耗时**：2 天
- 📚 **推荐扩展**：[GORM Docs](https://gorm.io/docs/)

#### 10. 缓存与消息 ─ `stage-2-business/10-cache-and-mq/`

- 🎯 **学习目标**：Redis 缓存的常见模式 + NATS/Kafka 消息系统的入门用法。
- 🧩 **关键知识点**：
  - Redis 客户端（`go-redis`）基础与连接池
  - 缓存模式：Cache-Aside、Read-Through、Write-Through
  - 缓存三大问题：穿透 / 雪崩 / 击穿 的工程对策
  - 分布式锁（Redis 实现）
  - 消息队列：NATS（轻量）或 Kafka（重量）的 Go 客户端
  - 消息的至少一次 / 最多一次 / 精确一次语义
- 📦 **本章产出**：_待 OpenSpec change `chapter-10-cache-and-mq` 填充_
- 🔗 **前置依赖**：09
- ⏱️ **预计耗时**：2 天
- 📚 **推荐扩展**：[go-redis](https://github.com/redis/go-redis)

#### 11. 可观测性 ─ `stage-2-business/11-observability/`

- 🎯 **学习目标**：让线上 Go 服务"可见"——结构化日志、metrics、trace 全套。
- 🧩 **关键知识点**：
  - 结构化日志：`slog`（标准库）
  - Prometheus metrics：Counter / Gauge / Histogram
  - OpenTelemetry trace：span / context 传播
  - 健康检查：liveness / readiness
  - 日志规范：traceID、level、字段约定
- 📦 **本章产出**：_待 OpenSpec change `chapter-11-observability` 填充_
- 🔗 **前置依赖**：10
- ⏱️ **预计耗时**：2 天
- 📚 **推荐扩展**：[OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)

#### 🚩 Capstone 2: 博客后端 API（单体）─ `stage-2-business/capstone-2-blog-api/`

- 🎯 **项目目标**：做一个**生产级单体**博客后端：用户注册登录 + 文章 CRUD + 评论 + 列表分页 + 完整可观测性。
- 🧩 **综合应用的章节**：08 / 09 / 10 / 11 全套
- 📦 **本项目产出**：_待 OpenSpec change `capstone-2-blog-api` 填充_
- 🔗 **前置依赖**：08–11
- ⏱️ **预计耗时**：5 天

---

### 阶段三：Go 架构进阶

#### 12. 整洁架构 ─ `stage-3-architecture/12-clean-architecture/`

- 🎯 **学习目标**：分层（domain / usecase / interface / infra）+ 依赖注入。
- 🧩 **关键知识点**：
  - 整洁架构 / 六边形架构核心思想
  - 四层分层：domain / usecase / interface / infrastructure
  - 依赖倒置：业务核心不依赖外部
  - `wire` 编译期依赖注入
  - 单元测试在分层架构下的玩法（mock 接口）
- 📦 **本章产出**：_待 OpenSpec change `chapter-12-clean-architecture` 填充_
- 🔗 **前置依赖**：capstone-2
- ⏱️ **预计耗时**：2 天
- 📚 **推荐扩展**：[wire](https://github.com/google/wire)、《架构整洁之道》

#### 13. DDD 战术模式 ─ `stage-3-architecture/13-ddd-patterns/`

- 🎯 **学习目标**：实体 / 值对象 / 聚合根 / 领域事件 / 仓储模式的 Go 落地。
- 🧩 **关键知识点**：
  - 实体（Entity）与值对象（Value Object）的 Go 表达
  - 聚合（Aggregate）与聚合根
  - 仓储（Repository）模式
  - 领域事件（Domain Event）与事件发布
  - 应用服务 vs 领域服务
- 📦 **本章产出**：_待 OpenSpec change `chapter-13-ddd-patterns` 填充_
- 🔗 **前置依赖**：12
- ⏱️ **预计耗时**：3 天
- 📚 **推荐扩展**：《领域驱动设计》Eric Evans

#### 14. 微服务 ─ `stage-3-architecture/14-microservices/`

- 🎯 **学习目标**：gRPC + protobuf + 服务发现 + 网关 + 配置中心，全套微服务基础设施。
- 🧩 **关键知识点**：
  - gRPC 与 protobuf：IDL、代码生成、流式
  - 服务发现：Consul / etcd / Nacos
  - API 网关：基础职责（路由、鉴权、限流、聚合）
  - 配置中心：动态配置 + 灰度
  - 服务间通信模式：同步（gRPC）vs 异步（MQ）
- 📦 **本章产出**：_待 OpenSpec change `chapter-14-microservices` 填充_
- 🔗 **前置依赖**：13
- ⏱️ **预计耗时**：4 天
- 📚 **推荐扩展**：[gRPC-Go](https://github.com/grpc/grpc-go)

#### 15. 韧性与性能 ─ `stage-3-architecture/15-resilience-perf/`

- 🎯 **学习目标**：从"能跑"到"扛得住"——限流 / 熔断 / 降级 / 重试，外加性能调优。
- 🧩 **关键知识点**：
  - 限流：令牌桶 / 漏桶 / 滑动窗口
  - 熔断器（circuit breaker）原理与实现（`gobreaker`）
  - 降级与服务隔离
  - 超时与重试策略（指数退避 / 抖动）
  - 性能调优：`pprof` 实战、火焰图、GC 调优
  - 压测：`vegeta` / `wrk`
- 📦 **本章产出**：_待 OpenSpec change `chapter-15-resilience-perf` 填充_
- 🔗 **前置依赖**：14
- ⏱️ **预计耗时**：3 天
- 📚 **推荐扩展**：Netflix Hystrix 论文、《数据密集型应用系统设计》

#### 🚩 Capstone 3: 微服务版博客 ─ `stage-3-architecture/capstone-3-blog-ms/`

- 🎯 **项目目标**：把 capstone-2 的**单体博客拆成 3 个微服务**（user-svc / post-svc / comment-svc），加网关、链路追踪、容器化部署。
- 🧩 **综合应用的章节**：12 / 13 / 14 / 15 全套
- 📦 **本项目产出**：_待 OpenSpec change `capstone-3-blog-ms` 填充_
- 🔗 **前置依赖**：12–15
- ⏱️ **预计耗时**：7 天

---

## 四、跟学工作流

### 命名规范

| 类型 | 命名 | 示例 |
|---|---|---|
| 章节内容落地 | `chapter-NN-<kebab-name>` | `chapter-01-hello-go` |
| Capstone 项目落地 | `capstone-N-<kebab-name>` | `capstone-1-cli-todo` |
| 路线图/模板修订 | `revise-<topic>` | `revise-roadmap-add-security` |

### 标准闭环（一章一 change）

```
1. 看 ROADMAP 找下一章
2. /opsx-propose chapter-NN-<name>   ← 提案（why + design + spec + tasks）
3. /opsx-apply chapter-NN-<name>     ← 实施（按 tasks 推进）
4. 自测清单全打勾
5. /opsx-archive chapter-NN-<name>   ← 归档
6. 回到 ROADMAP 勾选进度
```

---

## 进度追踪

### 阶段一：Go 语法精通

- [x] 01-hello-go
- [ ] 02-language-basics
- [ ] 03-composite-types
- [ ] 04-interface-error
- [ ] 05-concurrency
- [ ] 06-stdlib-essentials
- [ ] 07-engineering
- [ ] 🚩 capstone-1-cli-todo

### 阶段二：Go 业务工程

- [ ] 08-web-foundations
- [ ] 09-data-persistence
- [ ] 10-cache-and-mq
- [ ] 11-observability
- [ ] 🚩 capstone-2-blog-api

### 阶段三：Go 架构进阶

- [ ] 12-clean-architecture
- [ ] 13-ddd-patterns
- [ ] 14-microservices
- [ ] 15-resilience-perf
- [ ] 🚩 capstone-3-blog-ms
