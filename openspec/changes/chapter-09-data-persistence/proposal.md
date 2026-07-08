## Why

第 09 章目前仍是占位章节，尚未提供数据持久化示例、数据库测试和练习材料。完成本章可以承接第 08 章 Web API，帮助学习者掌握 `database/sql`、GORM、连接池、事务、迁移、防 SQL 注入和 N+1 查询治理。

## What Changes

- 在 `stage-2-business/09-data-persistence/` 下新增可运行入口程序，输出数据持久化学习报告。
- 新增主题子包，分别演示 `database/sql` CRUD、GORM 模型与关联、事务回滚、迁移脚本说明、prepared statement 与 N+1/预加载示例。
- 使用 SQLite 内存数据库作为教学与测试后端，避免本地必须启动 MySQL，同时在文档中说明 MySQL 生产语境与差异。
- 补充 `_test.go` 覆盖 CRUD、事务回滚、GORM 关联预加载、防注入查询和迁移应用。
- 将本章 `README.md` 从占位内容更新为实际产出说明、自测清单和运行命令，并新增 `EXERCISES.md`。
- 引入教学依赖：`gorm.io/gorm`、`gorm.io/driver/sqlite`、`github.com/mattn/go-sqlite3`。

## Capabilities

### New Capabilities
- `data-persistence-tutorial`: 覆盖第 09 章数据持久化学习单元的可运行代码、数据库测试、文档和练习。

### Modified Capabilities
- `learning-curriculum`: 将第 09 章从未落地占位章节更新为已落地章节，允许其目录包含源码、测试与练习材料。

## Impact

- 主要影响目录：`stage-2-business/09-data-persistence/`。
- 新增 OpenSpec 规格：`openspec/changes/chapter-09-data-persistence/specs/data-persistence-tutorial/spec.md`。
- 修改现有规格：`openspec/changes/chapter-09-data-persistence/specs/learning-curriculum/spec.md`。
- 更新依赖：`go.mod` / `go.sum` 增加 GORM 与 SQLite driver 相关依赖。
- 验证命令包括 `go test ./stage-2-business/09-data-persistence/...`、`go run ./stage-2-business/09-data-persistence`、`go test ./...` 和 `go build ./...`。
