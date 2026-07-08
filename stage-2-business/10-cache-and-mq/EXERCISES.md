# 10. 缓存与消息练习

## 练习 1：为缓存值增加 JSON 编码

在 `cachex` 旁新增一个小 helper，把结构体编码为 JSON 字符串后写入缓存，再读出并解码。

**验收标准：**

- 使用标准库 `encoding/json`。
- 测试覆盖编码、解码和非法 JSON 错误。
- `go test ./stage-2-business/10-cache-and-mq/cachex` 通过。

## 练习 2：实现锁续约

为 `LockManager` 增加 `Renew(resource, token string, ttl time.Duration) bool`。

**验收标准：**

- 只有当前持有者 token 可以续约。
- 续约后，原 TTL 到期点之后锁仍有效。
- 非持有者续约返回 false。

## 练习 3：观察缓存击穿

写一个测试：不用 `SingleFlightCache` 时 8 个 goroutine 同时 miss 会调用 loader 多次；使用 `SingleFlightCache` 后只调用一次。

**验收标准：**

- 测试能稳定复现 loader 调用次数差异。
- README 中用一句话解释 singleflight 适合热点 key 失效场景。

## 练习 4：实现最多一次消费

在 `mqdemo` 中新增 `FetchAtMostOnce`：消息一被取出就从队列删除，不需要 Ack。

**验收标准：**

- Fetch 后即使消费者不 Ack，消息也不会重投。
- README 中对比最多一次与至少一次的丢失/重复风险。

## 练习 5：补充真实 Redis/NATS 替换说明

在 README 中新增一节，说明如何把 `cachex.Store` 替换为 go-redis，把 `mqdemo.Broker` 替换为 NATS。

**验收标准：**

- 给出 go-redis `SetNX`、`Get`、`Set` 的伪代码。
- 给出 NATS publish/subscribe/ack 的伪代码或命令。
- 明确本章测试仍使用内存组件以保持 CI 稳定。
