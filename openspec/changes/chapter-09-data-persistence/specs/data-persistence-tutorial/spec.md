## ADDED Requirements

### Requirement: 第 09 章 SHALL 提供可运行的数据持久化入口程序

第 09 章目录 `stage-2-business/09-data-persistence/` SHALL 包含一个 `package main` 的入口文件 `main.go`，其 `main` 函数 MUST 通过调用本章子包导出函数生成持久化学习报告，并 MUST 能通过 `go run` 成功执行。

#### Scenario: go run 跑通入口程序
- **WHEN** 学习者在仓库根目录执行 `go run ./stage-2-business/09-data-persistence`
- **THEN** 程序 MUST 以退出码 0 结束，并在标准输出打印包含 `database/sql`、GORM、事务、迁移、预加载等关键词的报告文本

#### Scenario: 入口程序不承载全部业务逻辑
- **WHEN** 阅读者查看 `stage-2-business/09-data-persistence/main.go`
- **THEN** 该文件 MUST 通过 import 引入本章至少三个主题子包，并调用其导出函数，而非把全部逻辑写在 `main.go` 内

### Requirement: 第 09 章 SHALL 使用 database/sql 实现 CRUD

第 09 章 SHALL 包含 `sqlcrud` 子包，使用 `database/sql` 在 SQLite 内存数据库中实现文章的创建、查询、更新、删除和列表操作。查询 MUST 使用参数绑定，MUST NOT 通过字符串拼接注入用户输入。

#### Scenario: database/sql CRUD 全流程通过
- **WHEN** 测试代码依次创建、查询、更新、列表、删除文章
- **THEN** 每一步 MUST 返回符合预期的文章数据或删除结果

#### Scenario: 参数绑定阻止 SQL 注入
- **WHEN** 测试代码使用包含 SQL 片段的标题作为查询输入
- **THEN** 查询 MUST 只按普通字符串匹配，不得破坏表结构或返回非匹配记录

### Requirement: 第 09 章 SHALL 演示连接池与迁移

第 09 章 SHALL 包含 `dbx` 子包与 `migrations/` 目录。`dbx` MUST 提供打开 SQLite 内存数据库、配置连接池参数、执行迁移 SQL 的函数；迁移 SQL MUST 至少创建 articles 表。

#### Scenario: 迁移创建 articles 表
- **WHEN** 测试代码打开内存数据库并执行迁移 SQL
- **THEN** 后续插入 articles 表 MUST 成功

#### Scenario: 连接池配置可断言
- **WHEN** 测试代码调用连接池配置函数
- **THEN** 返回的配置摘要 MUST 包含最大打开连接数、最大空闲连接数和连接最大生命周期

### Requirement: 第 09 章 SHALL 使用 GORM 演示模型、CRUD 与关联

第 09 章 SHALL 包含 `gormdemo` 子包，使用 GORM 定义 User、Post、Tag 等模型，演示 AutoMigrate、CRUD、一对多、多对多关联和 `Preload`。

#### Scenario: GORM AutoMigrate 与 CRUD 通过
- **WHEN** 测试代码使用 GORM AutoMigrate 后创建用户和文章
- **THEN** 查询 MUST 能返回创建的数据

#### Scenario: GORM Preload 加载关联
- **WHEN** 测试代码创建用户及其文章并用 `Preload("Posts")` 查询用户
- **THEN** 返回用户 MUST 包含已加载的文章集合

#### Scenario: GORM 多对多标签关联可查询
- **WHEN** 测试代码创建文章并关联多个标签
- **THEN** 使用 `Preload("Tags")` 查询文章 MUST 返回对应标签集合

### Requirement: 第 09 章 SHALL 演示事务提交与回滚

第 09 章 SHALL 包含 `txdemo` 子包，演示成功事务提交与失败事务回滚。事务示例 MUST 可由单元测试断言。

#### Scenario: 成功事务提交数据
- **WHEN** 测试代码执行成功事务创建两篇文章
- **THEN** 事务结束后两篇文章 MUST 都可查询到

#### Scenario: 失败事务回滚数据
- **WHEN** 测试代码在事务中制造错误
- **THEN** 事务结束后该事务内写入的数据 MUST 不存在

### Requirement: 第 09 章 SHALL 提供单元测试

第 09 章 SHALL 为 `dbx`、`sqlcrud`、`gormdemo`、`txdemo` 或相关子包提供 `_test.go` 测试，且 MUST 使用真实 SQLite 内存数据库覆盖持久化行为。

#### Scenario: 章节测试全部通过
- **WHEN** 学习者在仓库根目录执行 `go test ./stage-2-business/09-data-persistence/...`
- **THEN** 所有测试 MUST 通过，命令 MUST 以退出码 0 返回

#### Scenario: 测试覆盖关键持久化语义
- **WHEN** 阅读者查看本章测试文件
- **THEN** 测试 MUST 覆盖 CRUD、事务回滚、迁移执行、GORM 关联预加载和防 SQL 注入场景

### Requirement: 第 09 章 SHALL 提供练习题与产出说明

第 09 章 SHALL 包含一个 `EXERCISES.md`，列出 3~5 道由浅入深的练习题（每题含验收标准）；同时章节 `README.md` 的“📦 本章产出”段落 MUST 被替换为实际内容（示例文件清单 + 运行命令），自测清单 MUST 与 ROADMAP 关键知识点对齐，且 MUST NOT 再包含“待 OpenSpec change 填充”或“尚未实现”等占位语。

#### Scenario: EXERCISES.md 含带验收标准的练习
- **WHEN** 阅读者打开 `stage-2-business/09-data-persistence/EXERCISES.md`
- **THEN** 该文件 MUST 包含 3 到 5 道练习题，且每道题 MUST 给出明确的验收标准

#### Scenario: README 产出说明不再是占位
- **WHEN** 阅读者打开 `stage-2-business/09-data-persistence/README.md` 的“📦 本章产出”段落
- **THEN** 该段落 MUST NOT 再包含“待 OpenSpec change 填充”或“尚未实现”占位语，且 MUST 列出本章 `.go` 文件清单与运行命令

#### Scenario: README 自测清单与知识点对齐
- **WHEN** 阅读者查看该 README 的“✅ 自测清单”
- **THEN** 该清单 MUST 包含与 ROADMAP 关键知识点对应的可勾选项，包括 `database/sql`、GORM、事务、迁移、防 SQL 注入、连接池和 N+1/Preload
