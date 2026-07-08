# 09. 数据持久化练习

## 练习 1：为 Article 增加分页列表

在 `sqlcrud.Repository` 中新增 `ListPage(limit, offset int)`，使用 SQL `LIMIT ? OFFSET ?` 返回分页结果。

**验收标准：**

- 使用参数绑定传入 limit 和 offset。
- 表驱动测试覆盖第一页、第二页、超出范围三种情况。
- `go test ./stage-2-business/09-data-persistence/sqlcrud` 通过。

## 练习 2：增加一次事务失败场景

在 `txdemo` 中新增一个事务函数：先插入一篇文章，再执行一条必然失败的 SQL，最后确认事务回滚。

**验收标准：**

- 函数返回非 nil error。
- 测试断言失败事务中的文章不存在。
- 不吞掉原始 SQL 错误，调用方能看到失败原因。

## 练习 3：观察 GORM Preload 前后的差异

为 `gormdemo` 添加一个测试或示例，分别用不带 `Preload` 和带 `Preload("Posts")` 的方式查询同一用户。

**验收标准：**

- 不带 `Preload` 时 `len(user.Posts) == 0`。
- 带 `Preload("Posts")` 时 `len(user.Posts) > 0`。
- README 中用一句话解释这与 N+1 查询治理的关系。

## 练习 4：新增第二个迁移文件

新增 `migrations/002_add_article_status.sql`，为 articles 表增加 `status` 字段，并更新迁移执行测试。

**验收标准：**

- 测试连续执行 001 和 002 两个迁移文件。
- 插入文章时能写入 `status` 字段。
- README 中说明迁移文件命名为什么要带递增序号。

## 练习 5：改造为 MySQL DSN 示例

不要求连接真实 MySQL，只在 README 中补充一段“如何把本章 SQLite 示例迁移到 MySQL”的步骤。

**验收标准：**

- 说明需要替换 driver import、GORM dialector、`database/sql` driver name 和 DSN。
- 给出 `parseTime=true` 的 DSN 示例。
- 明确测试仍默认使用 SQLite 内存库以保持 CI 稳定。
