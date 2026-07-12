## MODIFIED Requirements

### Requirement: 每个章节与 capstone 目录 SHALL 仅包含 README 与占位文件

章节目录与 capstone 目录的内容约束 SHALL 依其落地状态区分：

- **未落地章节**（尚无对应 `chapter-NN-*` / `capstone-N-*` change 实现）MUST 仅包含使用统一模板生成的 `README.md` 与用于 Git 跟踪空目录的 `.gitkeep`，且 MUST NOT 包含任何 `.go` 文件。
- **已落地章节**（对应 OpenSpec change 已落地）MAY 包含示例源码（`.go`）、测试文件（`_test.go`）以及 `EXERCISES.md` 等学习材料；其 `README.md` MUST 不再保留"待 OpenSpec change 填充"占位。

无论落地与否，仓库整体 MUST 始终保持 `go build ./...` 成功。

#### Scenario: 未落地章节目录内容受严格约束

- **WHEN** 阅读者列出一个尚未落地的章节目录（如 `stage-3-architecture/capstone-3-blog-ms/`）
- **THEN** 该目录 MUST 仅包含 `README.md` 与 `.gitkeep` 两个文件，且 MUST NOT 存在任何后缀为 `.go` 的文件

#### Scenario: 已落地章节目录可包含源码与测试

- **WHEN** 阅读者列出一个已落地的章节目录（如 `stage-3-architecture/15-resilience-perf/`）
- **THEN** 该目录 MAY 包含 `.go` 源码、`_test.go` 测试与 `EXERCISES.md`，且其 `README.md` MUST NOT 再包含"待 OpenSpec change 填充"占位语

#### Scenario: 仓库整体可通过 go build

- **WHEN** 在仓库根目录执行 `go build ./...`
- **THEN** 命令 MUST 以退出码 0 成功返回

### Requirement: 仓库 SHALL 提供顶层学习路线图

仓库根目录 SHALL 包含 `ROADMAP.md` 文件，作为整个学习路径的单一可信源。该文件 MUST 完整呈现三阶段、15 章、3 个 capstone 项目的全貌，并提供学习进度追踪机制。

#### Scenario: ROADMAP.md 存在且结构完整

- **WHEN** 任意阅读者在仓库根目录打开 `ROADMAP.md`
- **THEN** 该文件 MUST 至少包含以下六个段落：① 头部（定位 / 适用人群 / 学习方式）② 三段总览 ③ 15 章详表（按统一卡片模板）④ 3 个 capstone 说明 ⑤ 跟学方式（每章一个 OpenSpec change）⑥ 进度追踪表（含 15 个章节 + 3 个 capstone 的 checkbox 项）

#### Scenario: Chapter 15 路线图展示已落地产出

- **WHEN** 阅读者查看 ROADMAP 中第 15 章 `15-resilience-perf`
- **THEN** 该章节产出 MUST 描述可运行的商品详情韧性 Gateway、fake upstream、token bucket、bulkhead、retry、`gobreaker`、fallback、pprof、生产级方案对照、测试与练习
- **AND** 进度追踪表 MUST 将 `15-resilience-perf` 标记为已完成
- **AND** `capstone-3-blog-ms` MUST 保持未完成

#### Scenario: 每章按统一卡片模板呈现

- **WHEN** 阅读者查看 ROADMAP 中任意一章
- **THEN** 该章节卡片 MUST 包含以下六个字段：`🎯 学习目标`、`🧩 关键知识点`、`📦 章节产出`、`🔗 前置依赖`、`⏱️ 预计耗时`、`📚 推荐扩展`

#### Scenario: 进度追踪表覆盖所有学习单元

- **WHEN** 阅读者查看进度追踪表
- **THEN** 该表 MUST 按"阶段一 / 阶段二 / 阶段三"分组，且每组下 MUST 列出该阶段所有章节与 capstone 项目作为可勾选项（`- [ ]` 或 `- [x]` 格式）
