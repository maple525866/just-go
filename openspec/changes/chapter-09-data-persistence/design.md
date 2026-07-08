## Context

`stage-2-business/09-data-persistence/` 当前只有占位 README 与 `.gitkeep`。第 09 章目标是从第 08 章的内存 Web API 过渡到持久化层，覆盖 `database/sql`、GORM、连接池、事务、迁移、防 SQL 注入与 N+1 查询治理。

本仓库每章需要独立可运行且本地验证稳定。ROADMAP 目标写的是 MySQL，但 CI 与学习者本地不应强制依赖外部 MySQL 服务，因此本章使用 SQLite 内存数据库执行真实 SQL/GORM 测试，并在 README 中说明 API 迁移到 MySQL 时连接字符串和 driver 的差异。

## Goals / Non-Goals

**Goals:**

- 提供可运行的第 09 章入口程序，输出持久化学习报告。
- 用 `database/sql` 实现完整 CRUD，并演示连接池配置与 prepared statement 参数绑定。
- 用 GORM 实现用户、文章、评论/标签等模型，以及一对多、多对多关联和 `Preload`。
- 用测试证明事务提交与回滚行为。
- 提供迁移 SQL 与迁移执行函数，演示数据库版本管理的基本形态。
- README 与练习题替换占位内容，形成完整学习单元。

**Non-Goals:**

- 不要求学习者本地启动 MySQL、Docker 或外部数据库。
- 不引入真实线上迁移工具命令作为测试前置；迁移工具以 SQL 文件与可复制命令说明呈现。
- 不实现复杂业务仓储或 Web handler；本章聚焦持久化层。
- 不处理数据库性能压测或分库分表。

## Decisions

### 1. 测试后端使用 SQLite 内存数据库

本章的 `sqlcrud` 与 `gormdemo` 子包都使用 SQLite `:memory:` 数据库运行测试。代码保持 SQL 参数绑定、事务、迁移、连接池等通用概念，README 额外说明 MySQL driver 与 DSN 的替换方式。

**Rationale:** SQLite 内存数据库能让测试无需外部服务且运行快速，仍可验证真实 SQL 行为。

**Alternative considered:** 使用 MySQL testcontainer。该方案更贴近 ROADMAP，但会显著增加本地和 CI 依赖复杂度。

### 2. 子包按持久化主题拆分

- `dbx/`：打开 SQLite 内存 DB、连接池配置、迁移执行 helper。
- `sqlcrud/`：使用 `database/sql` 实现 Article CRUD 与 prepared statement 查询。
- `gormdemo/`：GORM 模型、AutoMigrate、CRUD、关联、Preload。
- `txdemo/`：事务提交与回滚示例。
- `migrations/`：SQL 迁移文件。

**Rationale:** 持久化概念较多，按主题拆分能让测试和文档分别聚焦。

**Alternative considered:** 只用 GORM。该方案会跳过 `database/sql` 基础，不符合 ROADMAP。

### 3. N+1 治理以可断言的预加载结果呈现

GORM 示例中创建用户和文章，通过 `Preload("Posts")` 一次性加载关联，并在测试中断言用户携带文章集合。

**Rationale:** 初学者可以先理解“关联未加载 vs 预加载”的行为差异，再深入查询计数和性能分析。

**Alternative considered:** 在测试中统计 SQL 查询条数。该方案需要 logger hook，复杂度较高，容易分散章节重点。

### 4. 迁移用 SQL 文件 + helper 演示

仓库保存 `migrations/001_create_articles.sql`，`dbx.ApplyMigration` 读取并执行 SQL。README 同时给出 goose/golang-migrate 的真实命令形态。

**Rationale:** 测试可稳定验证迁移效果，同时不把外部迁移工具安装作为硬依赖。

**Alternative considered:** 直接使用 goose 库。该方案会引入额外工具概念，本章先讲迁移思想和 SQL 版本文件即可。

## Risks / Trade-offs

- [Risk] SQLite 与 MySQL 语法/行为存在差异 → Mitigation：README 明确 SQLite 是教学测试后端，并列出迁移到 MySQL 的 driver/DSN 差异。
- [Risk] GORM 关联示例可能过度简化 N+1 → Mitigation：练习要求学习者增加查询日志或计数来观察 Preload 前后差异。
- [Risk] 不使用真实迁移工具可能不够生产化 → Mitigation：保留 SQL 迁移文件和工具命令说明，capstone 再落地完整工具链。
