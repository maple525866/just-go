# 09. 数据持久化

> 阶段：② 业务工程 ｜ 难度：⭐⭐⭐☆☆ ｜ 预计耗时：2 天

## 🎯 学习目标

用 `database/sql` + GORM 操作 MySQL，理解连接池、事务、迁移。

## 🧩 关键知识点

- `database/sql` 基础与连接池
- GORM 模型定义、CRUD、关联
- 事务（`db.Transaction`）、隔离级别
- 数据库迁移（`goose` / `golang-migrate`）
- 防 SQL 注入与 prepared statement
- N+1 与预加载（`Preload` / `Joins`）

## 📦 本章产出（待 OpenSpec change 填充）

> ⚠️ 当前本章内容尚未实现。
>
> 请通过 `/opsx-propose chapter-09-data-persistence` 来落地本章内容。

## 🔗 前置依赖

- 第 08 章

## 📚 推荐扩展阅读

- [GORM Docs](https://gorm.io/docs/)
- [database/sql tutorial](http://go-database-sql.org/)
- [golang-migrate](https://github.com/golang-migrate/migrate)

## ✅ 自测清单（落地后填充）

- [ ] 能用 `database/sql` 写一个完整的 CRUD（不借助 ORM）
- [ ] 能用 GORM 实现一对多 / 多对多关联
- [ ] 能正确处理事务回滚
- [ ] 能识别并消除一段 N+1 查询
- [ ] 能用迁移工具管理数据库版本
