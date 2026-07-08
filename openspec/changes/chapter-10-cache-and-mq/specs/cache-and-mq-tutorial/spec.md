## ADDED Requirements

### Requirement: 第 10 章 SHALL 提供可运行的缓存与消息入口程序

第 10 章目录 `stage-2-business/10-cache-and-mq/` SHALL 包含一个 `package main` 的入口文件 `main.go`，其 `main` 函数 MUST 通过调用本章子包导出函数生成缓存与消息学习报告，并 MUST 能通过 `go run` 成功执行。

#### Scenario: go run 跑通入口程序
- **WHEN** 学习者在仓库根目录执行 `go run ./stage-2-business/10-cache-and-mq`
- **THEN** 程序 MUST 以退出码 0 结束，并在标准输出打印包含 Cache-Aside、Read-Through、Write-Through、分布式锁、消息 ack 等关键词的报告文本

#### Scenario: 入口程序不承载全部业务逻辑
- **WHEN** 阅读者查看 `stage-2-business/10-cache-and-mq/main.go`
- **THEN** 该文件 MUST 通过 import 引入本章至少两个主题子包，并调用其导出函数，而非把全部逻辑写在 `main.go` 内

### Requirement: 第 10 章 SHALL 演示缓存模式

第 10 章 SHALL 包含 `cachex` 子包，提供带 TTL 的内存 key-value store，并实现 Cache-Aside、Read-Through、Write-Through 三种缓存模式。缓存示例 MUST 可由单元测试断言命中、未命中、回源和写穿行为。

#### Scenario: Cache-Aside miss 后写入缓存
- **WHEN** 测试代码第一次读取不存在于缓存但存在于源数据的 key
- **THEN** loader MUST 被调用一次，返回值 MUST 写入缓存，第二次读取同一 key MUST 命中缓存且不再调用 loader

#### Scenario: Read-Through 封装加载逻辑
- **WHEN** 测试代码通过 read-through cache 读取缺失 key
- **THEN** cache MUST 自动调用内部 loader 并缓存结果

#### Scenario: Write-Through 同步写源与缓存
- **WHEN** 测试代码通过 write-through cache 写入 key/value
- **THEN** 源数据与缓存中 MUST 都能读到新值

### Requirement: 第 10 章 SHALL 演示缓存穿透、雪崩、击穿对策

第 10 章 SHALL 在 `cachex` 子包中提供负缓存、TTL jitter 计算和互斥加载机制，分别对应缓存穿透、雪崩和击穿治理。

#### Scenario: 负缓存避免重复穿透
- **WHEN** 测试代码连续读取不存在于源数据的 key
- **THEN** 第一次读取 MUST 写入负缓存，第二次读取 MUST 不再调用 loader

#### Scenario: TTL jitter 产生不同过期时间
- **WHEN** 测试代码对多个 key 计算带 jitter 的 TTL
- **THEN** 返回 TTL MUST 在允许范围内且至少存在两个不同值

#### Scenario: 互斥加载避免同 key 并发击穿
- **WHEN** 多个 goroutine 同时读取同一个缓存 miss key
- **THEN** loader MUST 只执行一次，所有 goroutine MUST 获得相同结果

### Requirement: 第 10 章 SHALL 演示 Redis 风格分布式锁

第 10 章 SHALL 在 `cachex` 子包中实现教学版分布式锁。锁 MUST 使用 token 标识持有者，MUST 支持 TTL，MUST 只允许持有者释放。

#### Scenario: 同一资源同一时间只能被一个 token 持有
- **WHEN** token A 已成功获取某资源锁且未过期
- **THEN** token B 获取同一资源锁 MUST 失败

#### Scenario: 只有持有者可以释放锁
- **WHEN** token B 尝试释放 token A 持有的锁
- **THEN** 释放 MUST 失败且锁仍然存在

#### Scenario: 锁过期后可重新获取
- **WHEN** 锁 TTL 到期
- **THEN** 新 token MUST 能重新获取该资源锁

### Requirement: 第 10 章 SHALL 演示消息生产消费与 ack 语义

第 10 章 SHALL 包含 `mqdemo` 子包，提供内存 broker，支持 Publish、Fetch、Ack、RequeueExpired。消息示例 MUST 演示至少一次投递语义。

#### Scenario: 发布后可消费消息
- **WHEN** 测试代码发布一条消息并调用 Fetch
- **THEN** Fetch MUST 返回该消息及其投递 ID

#### Scenario: ack 后消息不再投递
- **WHEN** 消费者 Fetch 一条消息后调用 Ack
- **THEN** 后续 Fetch MUST 不再返回该消息

#### Scenario: 未 ack 消息可重投
- **WHEN** 消费者 Fetch 一条消息但不 Ack，且可见性超时到期
- **THEN** RequeueExpired 后再次 Fetch MUST 返回同一消息

### Requirement: 第 10 章 SHALL 提供单元测试

第 10 章 SHALL 为 `cachex`、`mqdemo` 或相关子包提供 `_test.go` 测试，覆盖缓存模式、缓存问题对策、分布式锁和消息 ack/retry 行为。

#### Scenario: 章节测试全部通过
- **WHEN** 学习者在仓库根目录执行 `go test ./stage-2-business/10-cache-and-mq/...`
- **THEN** 所有测试 MUST 通过，命令 MUST 以退出码 0 返回

#### Scenario: 测试覆盖关键缓存与消息语义
- **WHEN** 阅读者查看本章测试文件
- **THEN** 测试 MUST 覆盖 Cache-Aside、Read-Through、Write-Through、负缓存、TTL jitter、互斥加载、锁 token 校验、消息 ack 和未 ack 重投

### Requirement: 第 10 章 SHALL 提供练习题与产出说明

第 10 章 SHALL 包含一个 `EXERCISES.md`，列出 3~5 道由浅入深的练习题（每题含验收标准）；同时章节 `README.md` 的“📦 本章产出”段落 MUST 被替换为实际内容（示例文件清单 + 运行命令），自测清单 MUST 与 ROADMAP 关键知识点对齐，且 MUST NOT 再包含“待 OpenSpec change 填充”或“尚未实现”等占位语。

#### Scenario: EXERCISES.md 含带验收标准的练习
- **WHEN** 阅读者打开 `stage-2-business/10-cache-and-mq/EXERCISES.md`
- **THEN** 该文件 MUST 包含 3 到 5 道练习题，且每道题 MUST 给出明确的验收标准

#### Scenario: README 产出说明不再是占位
- **WHEN** 阅读者打开 `stage-2-business/10-cache-and-mq/README.md` 的“📦 本章产出”段落
- **THEN** 该段落 MUST NOT 再包含“待 OpenSpec change 填充”或“尚未实现”占位语，且 MUST 列出本章 `.go` 文件清单与运行命令

#### Scenario: README 自测清单与知识点对齐
- **WHEN** 阅读者查看该 README 的“✅ 自测清单”
- **THEN** 该清单 MUST 包含与 ROADMAP 关键知识点对应的可勾选项，包括 Cache-Aside、Read-Through、Write-Through、穿透/雪崩/击穿、分布式锁、消息生产消费、至少一次语义
