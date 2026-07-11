# 第 13 章练习：DDD 战术模式

建议先运行本章全部测试，再按顺序完成练习。每个练习都要求先写失败测试，再实现行为。

## 练习 1：修改收货地址（基础）

为草稿订单增加 `ChangeShippingAddress` 行为。

要求：

- 只接受由 `NewAddress` 创建的有效值对象。
- 草稿状态修改成功并推进聚合版本。
- 已确认订单返回 `domain.ErrInvalidState`。
- `Address()` 返回值不能让调用方绕过聚合行为修改内部状态。

验收命令：

```bash
go test ./stage-3-architecture/13-ddd-patterns/domain/order -run ChangeShippingAddress
```

## 练习 2：阶梯折扣策略（基础）

实现 `ThresholdDiscount`：小计达到指定门槛后减去固定金额，否则不打折。

要求：

- 门槛、折扣额和小计必须使用相同货币。
- 折扣额不能超过小计。
- 策略保持无状态，不导入 `application` 或 `infrastructure` 包。
- 覆盖未达门槛、恰好达到、超过门槛和币种不一致四组测试。

## 练习 3：修改订单行数量（进阶）

在不暴露可变 `Line` 的前提下，为 `Order` 增加 `ChangeLineQuantity`。

要求：

- 只能通过 `LineID` 定位实体。
- 数量必须大于 0；不存在的行返回 `domain.ErrLineNotFound`。
- 已确认订单不可修改。
- 修改成功只推进一次版本号。

思考：为什么不直接给 `Line` 增加一个公开的 `SetQuantity`？

## 练习 4：持久化事件幂等性（进阶）

当前库存投影只在进程内按 `OrderID` 去重。请把幂等记录抽象为可替换的 `ProcessedEventStore`，并为领域事件增加独立 `EventID`，使重启后的消费者仍能识别重复投递。

要求：

- `OrderConfirmed` 创建时获得稳定且非空的 `EventID`，重试时 ID 不变。
- `InventoryProjection` 通过端口查询和记录已处理 ID，不直接依赖数据库实现。
- 模拟消费者重启后重复投递同一事件，库存数量不能重复增加。
- 两个不同事件确认相同商品时，数量正确累加。
- 使用 `go test -race` 验证并发投递，并说明“更新投影 + 记录 ID”需要原子性的原因。

验收命令：

```bash
go test -race ./stage-3-architecture/13-ddd-patterns/infrastructure/events
```

## 练习 5：设计 transactional outbox（挑战）

不用引入真实数据库，先写一份 `OUTBOX_DESIGN.md`，再为端口设计测试替身。

文档至少回答：

- 聚合数据与 outbox 记录如何在同一事务中提交？
- 后台发布器如何重试并避免重复消费？
- 事件需要哪些 ID、类型、版本和时间字段？
- 发布成功后如何清理或归档？
- 进程在“消息已发送、状态未更新”时崩溃会发生什么？

验收标准：应用服务不再直接依赖即时发布成功；仓储失败仍然不能产生可发布事件；测试覆盖提交、重试和重复投递。

## 完成检查

```bash
gofmt -w stage-3-architecture/13-ddd-patterns
go test ./stage-3-architecture/13-ddd-patterns/...
go test -race ./stage-3-architecture/13-ddd-patterns/...
go vet ./stage-3-architecture/13-ddd-patterns/...
```
