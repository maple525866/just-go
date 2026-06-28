# 09. 数据持久化

> 阶段：② 业务工程 ｜ 难度：⭐⭐⭐☆☆ ｜ 预计耗时：2 天

## 🎯 学习目标

用 `database/sql` + GORM 操作关系型数据库，理解连接池、事务、迁移、防 SQL 注入与 N+1 查询治理。

## 🧩 关键知识点

- `database/sql` 基础、CRUD 与连接池
- GORM 模型定义、CRUD、一对多与多对多关联
- 事务提交与回滚
- 数据库迁移 SQL 与迁移工具形态
- 防 SQL 注入与 prepared statement / 参数绑定
- N+1 查询与 `Preload` 预加载

## 📦 本章产出

本章使用 SQLite 内存数据库运行真实 SQL/GORM 示例，避免本地必须启动 MySQL；生产或 capstone 中可把 driver 和 DSN 换成 MySQL。

```text
stage-2-business/09-data-persistence/
├── main.go                         # 输出持久化学习报告
├── dbx/                            # SQLite 内存 DB、迁移执行、连接池配置
├── sqlcrud/                        # database/sql Article CRUD + 参数绑定
├── gormdemo/                       # GORM 模型、AutoMigrate、关联与 Preload
├── txdemo/                         # 事务提交与回滚示例
└── migrations/001_create_articles.sql
```

运行测试：

```bash
go test ./stage-2-business/09-data-persistence/...
```

运行示例：

```bash
go run ./stage-2-business/09-data-persistence
```

SQLite 与 MySQL 对照：

| 主题 | 本章 SQLite | MySQL 生产语境 |
|---|---|---|
| driver | `github.com/mattn/go-sqlite3` | `github.com/go-sql-driver/mysql` |
| DSN | `:memory:` | `user:pass@tcp(localhost:3306)/db?parseTime=true` |
| 迁移 | 读取 SQL 文件并执行 | goose / golang-migrate 管理版本 |
| 测试 | 无外部服务依赖 | 可用 Docker/Testcontainers 做集成测试 |

迁移工具命令示例：

```bash
# goose 示例
goose -dir ./migrations mysql "$MYSQL_DSN" up

# golang-migrate 示例
migrate -path ./migrations -database "$MYSQL_DSN" up
```

## 🔗 前置依赖

- 第 08 章

## 📚 推荐扩展阅读

- [GORM Docs](https://gorm.io/docs/)
- [database/sql tutorial](http://go-database-sql.org/)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [goose](https://github.com/pressly/goose)

## ✅ 自测清单

- [ ] 能用 `database/sql` 写一个完整的 CRUD（不借助 ORM）。
- [ ] 能解释连接池中的最大打开连接、最大空闲连接和连接最大生命周期。
- [ ] 能用 GORM 实现一对多与多对多关联。
- [ ] 能正确处理事务提交与回滚。
- [ ] 能识别并消除一段 N+1 查询，知道何时使用 `Preload`。
- [ ] 能说明参数绑定如何避免 SQL 注入。
- [ ] 能用迁移 SQL 或迁移工具管理数据库版本。
