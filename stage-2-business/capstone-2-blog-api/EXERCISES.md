# Capstone 2 练习

## 练习 1：为文章更新与删除增加所有权校验

当前 `PUT /api/articles/{id}` 与 `DELETE /api/articles/{id}` 只要求登录；请扩展为只有作者本人可以更新或删除。

验收标准：

- 作者本人更新/删除返回 200/204。
- 其他登录用户更新/删除返回 403。
- 更新后再次 `GET /api/articles/{id}` 返回新内容。
- 删除后详情返回 404，列表不再出现该文章。
- 更新和删除都必须让文章详情缓存失效。

## 练习 2：把内存仓库替换为 GORM 仓库

基于第 09 章，把 `store.MemoryStore` 的接口形态迁移到 SQLite/GORM 实现。

验收标准：

- 保留现有 HTTP 测试，并新增 GORM repository 测试。
- 用户名唯一约束由数据库保证。
- 评论软删除在数据库中保留记录但隐藏正文。
- README 说明 SQLite 与 MySQL DSN 替换方式。

## 练习 3：加入异步通知事件

创建文章后发布一条 `ArticleCreated` 事件，使用第 10 章的内存 MQ demo 或等价接口消费。

验收标准：

- 创建文章成功后 broker 中能消费到事件。
- 消费者 ack 后事件不再重投。
- 未 ack 超时后事件会被重新投递。
- metrics 增加事件发布与消费计数。

## 练习 4：接入真实 OpenTelemetry SDK

用 `go.opentelemetry.io/otel` 替换教学版 `tracex`。

验收标准：

- HTTP 请求会生成 server span。
- store 和 cache 操作生成 child span。
- trace ID 继续写入响应头和结构化日志。
- 测试使用内存 exporter 验证 span 关系。
