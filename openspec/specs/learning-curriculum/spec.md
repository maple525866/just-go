# learning-curriculum Specification

## Purpose
TBD - created by archiving change bootstrap-learning-repo. Update Purpose after archive.
## Requirements
### Requirement: 仓库 SHALL 提供顶层学习路线图

仓库根目录 SHALL 包含 `ROADMAP.md` 文件，作为整个学习路径的单一可信源。该文件 MUST 完整呈现三阶段、15 章、3 个 capstone 项目的全貌，并提供学习进度追踪机制。

#### Scenario: ROADMAP.md 存在且结构完整

- **WHEN** 任意阅读者在仓库根目录打开 `ROADMAP.md`
- **THEN** 该文件 MUST 至少包含以下六个段落：① 头部（定位 / 适用人群 / 学习方式）② 三段总览 ③ 15 章详表（按统一卡片模板）④ 3 个 capstone 说明 ⑤ 跟学方式（每章一个 OpenSpec change）⑥ 进度追踪表（含 15 个章节 + 3 个 capstone 的 checkbox 项）

#### Scenario: 每章按统一卡片模板呈现

- **WHEN** 阅读者查看 ROADMAP 中任意一章
- **THEN** 该章节卡片 MUST 包含以下六个字段：`🎯 学习目标`、`🧩 关键知识点`、`📦 章节产出`、`🔗 前置依赖`、`⏱️ 预计耗时`、`📚 推荐扩展`

#### Scenario: 进度追踪表覆盖所有学习单元

- **WHEN** 阅读者查看进度追踪表
- **THEN** 该表 MUST 按"阶段一 / 阶段二 / 阶段三"分组，且每组下 MUST 列出该阶段所有章节与 capstone 项目作为可勾选项（`- [ ]` 格式）

### Requirement: 仓库 SHALL 按三阶段物理分目录组织所有章节

仓库根目录 SHALL 包含三个顶层目录 `stage-1-syntax/`、`stage-2-business/`、`stage-3-architecture/`，分别承载语法精通、业务工程、架构进阶三个阶段的全部章节与 capstone 项目。

#### Scenario: 阶段一目录包含 7 章 + 1 capstone

- **WHEN** 阅读者列出 `stage-1-syntax/` 目录
- **THEN** 该目录下 MUST 存在如下 8 个子目录：`01-hello-go/`、`02-language-basics/`、`03-composite-types/`、`04-interface-error/`、`05-concurrency/`、`06-stdlib-essentials/`、`07-engineering/`、`capstone-1-cli-todo/`

#### Scenario: 阶段二目录包含 4 章 + 1 capstone

- **WHEN** 阅读者列出 `stage-2-business/` 目录
- **THEN** 该目录下 MUST 存在如下 5 个子目录：`08-web-foundations/`、`09-data-persistence/`、`10-cache-and-mq/`、`11-observability/`、`capstone-2-blog-api/`

#### Scenario: 阶段三目录包含 4 章 + 1 capstone

- **WHEN** 阅读者列出 `stage-3-architecture/` 目录
- **THEN** 该目录下 MUST 存在如下 5 个子目录：`12-clean-architecture/`、`13-ddd-patterns/`、`14-microservices/`、`15-resilience-perf/`、`capstone-3-blog-ms/`

### Requirement: 每个章节与 capstone 目录 SHALL 仅包含 README 与占位文件

章节目录与 capstone 目录的内容约束 SHALL 依其落地状态区分：

- **未落地章节**（尚无对应 `chapter-NN-*` / `capstone-N-*` change 实现）MUST 仅包含使用统一模板生成的 `README.md` 与用于 Git 跟踪空目录的 `.gitkeep`，且 MUST NOT 包含任何 `.go` 文件。
- **已落地章节**（对应 OpenSpec change 已落地）MAY 包含示例源码（`.go`）、测试文件（`_test.go`）以及 `EXERCISES.md` 等学习材料；其 `README.md` MUST 不再保留"待 OpenSpec change 填充"占位。

无论落地与否，仓库整体 MUST 始终保持 `go build ./...` 成功。

#### Scenario: 未落地章节目录内容受严格约束

- **WHEN** 阅读者列出一个尚未落地的章节目录（如 `stage-1-syntax/02-language-basics/`）
- **THEN** 该目录 MUST 仅包含 `README.md` 与 `.gitkeep` 两个文件，且 MUST NOT 存在任何后缀为 `.go` 的文件

#### Scenario: 已落地章节目录可包含源码与测试

- **WHEN** 阅读者列出一个已落地的章节目录（如 `stage-1-syntax/01-hello-go/`）
- **THEN** 该目录 MAY 包含 `.go` 源码、`_test.go` 测试与 `EXERCISES.md`，且其 `README.md` MUST NOT 再包含"待 OpenSpec change 填充"占位语

#### Scenario: 仓库整体可通过 go build

- **WHEN** 在仓库根目录执行 `go build ./...`
- **THEN** 命令 MUST 以退出码 0 成功返回

### Requirement: 章节 README SHALL 使用统一模板

每个章节目录下的 `README.md` MUST 遵循统一模板，且模板 MUST 显式声明本章内容"待 OpenSpec change 填充"，并指引学习者使用 `/opsx-propose chapter-NN-xxx` 来启动该章的实现。

#### Scenario: 章节 README 包含必备段落

- **WHEN** 阅读者打开任一章节目录下的 `README.md`
- **THEN** 该文件 MUST 包含以下段落：`# NN. 章节标题`、`🎯 学习目标`、`🧩 关键知识点`、`📦 本章产出（待 OpenSpec change 填充）`、`🔗 前置依赖`、`📚 推荐扩展阅读`、`✅ 自测清单（落地后填充）`

#### Scenario: 章节 README 引导后续工作流

- **WHEN** 阅读者查看任一章节 README 的"📦 本章产出"段落
- **THEN** 该段落 MUST 包含一句明确提示：调用 `/opsx-propose chapter-NN-xxx` 来落地本章的具体代码与练习

### Requirement: Capstone README SHALL 使用 capstone 专用模板

每个 capstone 目录下的 `README.md` MUST 使用 capstone 专用模板，明确列出本 capstone 综合应用的所有章节，并设置完成标准。

#### Scenario: Capstone README 列出综合应用的章节

- **WHEN** 阅读者打开任一 capstone 目录下的 `README.md`（如 `stage-1-syntax/capstone-1-cli-todo/README.md`）
- **THEN** 该文件 MUST 包含 `🧩 综合应用的章节` 段落，且 MUST 列出对应阶段内所有章节编号与知识点用途

#### Scenario: Capstone README 含完成标准

- **WHEN** 阅读者查看任一 capstone README
- **THEN** 该文件 MUST 包含 `✅ 完成标准（落地后填充）` 段落，且其中 MUST 包含"代码可运行 / 有测试 / 有 README 说明"三项底线要求

### Requirement: 仓库 SHALL 提供跨章节支持文档骨架

仓库根目录 SHALL 包含 `docs/` 目录，且其中 SHALL 至少存在三个文件：`glossary.md`（术语表）、`faq.md`（常见问题）、`references.md`（参考资料）。这三个文件在本次提案中只建骨架，内容由后续章节 change 增量追加。

#### Scenario: docs 目录结构完整

- **WHEN** 阅读者列出 `docs/` 目录
- **THEN** 该目录 MUST 至少包含 `glossary.md`、`faq.md`、`references.md` 三个文件

#### Scenario: glossary.md 提供表格骨架

- **WHEN** 阅读者打开 `docs/glossary.md`
- **THEN** 该文件 MUST 包含一个表头为 `| 术语 | 英文 | 解释 | 出现章节 |` 的 Markdown 表格（初始无数据行）

#### Scenario: references.md 预填三类资料

- **WHEN** 阅读者打开 `docs/references.md`
- **THEN** 该文件 MUST 至少分为三个段落：`官方资源`（含 go.dev、Tour of Go、Effective Go）、`中文书单`、`博客与社区`

### Requirement: 仓库 SHALL 配置基础工程文件

仓库根目录 SHALL 包含 `README.md`（仓库门面）、`LICENSE`（MIT）、`.gitignore`（Go 标准）、`.editorconfig` 四个工程基础文件，且 `go.mod` MUST 声明 `go 1.24`。

#### Scenario: 顶层文档与配置就位

- **WHEN** 阅读者列出仓库根目录
- **THEN** 该目录 MUST 包含 `README.md`、`ROADMAP.md`、`LICENSE`、`.gitignore`、`.editorconfig`、`go.mod` 六个文件

#### Scenario: go.mod 声明 Go 1.24

- **WHEN** 阅读者打开 `go.mod`
- **THEN** 该文件 MUST 包含 `module just-go` 与 `go 1.24` 两个语句

#### Scenario: 顶层 README 不复述路线图

- **WHEN** 阅读者打开根目录 `README.md`
- **THEN** 该文件 MUST 在显眼处提供指向 `ROADMAP.md` 的链接，且 MUST NOT 复制 15 章详表的内容（路线图作为单一可信源）

### Requirement: 仓库 SHALL 完成 Git 版本控制初始化

仓库 SHALL 在本提案落地时完成 `git init`，默认分支为 `main`，并将本提案产出的所有文件提交为首次 commit。

#### Scenario: Git 仓库就位

- **WHEN** 在仓库根目录运行 `git status`
- **THEN** 命令 MUST 不再报告"not a git repository"错误，且当前分支 MUST 为 `main`

#### Scenario: 首次 commit 包含本提案全部产出

- **WHEN** 在仓库根目录运行 `git log --oneline`
- **THEN** 输出 MUST 至少包含一条 commit，且其内容 MUST 涵盖 `ROADMAP.md`、`README.md`、所有阶段目录、`docs/` 以及 `openspec/changes/bootstrap-learning-repo/` 下的本提案文件

### Requirement: 学习者跟学工作流 SHALL 遵循统一命名规范

为保证后续 15+3 个独立 OpenSpec change 在命名与组织上一致，本提案 SHALL 在 `ROADMAP.md` 中明确以下命名规范，且未来所有相关 change MUST 遵守。

#### Scenario: 章节 change 命名规范

- **WHEN** 学习者准备启动任一章节内容的实现
- **THEN** 该章节对应的 OpenSpec change MUST 命名为 `chapter-NN-<kebab-name>`，其中 `NN` 为两位数章节编号（如 `chapter-01-hello-go`、`chapter-15-resilience-perf`）

#### Scenario: Capstone change 命名规范

- **WHEN** 学习者准备启动任一 capstone 项目的实现
- **THEN** 该 capstone 对应的 OpenSpec change MUST 命名为 `capstone-N-<kebab-name>`，其中 `N` 为单位数阶段编号（如 `capstone-1-cli-todo`、`capstone-3-blog-ms`）

#### Scenario: 跨章修订 change 命名规范

- **WHEN** 学习者需要对路线图或章节模板做修订
- **THEN** 对应的 OpenSpec change MUST 命名为 `revise-<topic>`（如 `revise-roadmap-add-security`、`revise-chapter-template-add-quiz`）

