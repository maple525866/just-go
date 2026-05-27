## Context

`just-go` 仓库目前处于"白纸"状态：

- 仅有 `go.mod`（一行 `module just-go`，无 `go` 版本声明）
- 无 `.go` 源码、无 README、无路线规划
- 尚未 `git init`
- 已配置好 OpenSpec 工作流（`openspec/config.yaml`）以及多套 AI assistant skills（`.cursor/`、`.claude/`、`.codex/`、`.opencode/`、`.pi/`）

仓库所有者（学习者）当前 Go 水平为 **L0（完全零基础）**，目标终点为"能用 Go 搭建业务架构的架构师"。学习者已经选定：

- **仓库形态**：章节教科书式（按主题分目录，非渐进项目式）
- **路线图方案**：三段式 15 章 + 3 个综合 capstone 项目（详见 proposal.md）
- **本次提案范围**：仅交付路线图 + 仓库骨架，不写章节内容

本次变更必须为后续 15 + 3 = 18 个独立 OpenSpec change 提供稳定的"对接契约"（章节命名、目录结构、README 模板、跟学工作流）。

## Goals / Non-Goals

**Goals:**

- 输出一份**单一可信源（Single Source of Truth）**式的 `ROADMAP.md`，让学习者随时知道"我在哪 / 下一步是什么 / 终点是什么"。
- 把三段式（语法 / 业务 / 架构）的叙事**物理化**到目录结构（`stage-1-syntax/` / `stage-2-business/` / `stage-3-architecture/`），让"阶段感"在文件系统层面可见。
- 为所有未来章节提案确立**统一的对接契约**：章节目录布局、README 模板、跟学工作流命名规范。
- 让 `go build ./...` 在零业务代码状态下**仍能成功通过**（不引入任何编译错误的占位代码）。
- 让仓库自带 Git 版本控制，从第一天起就具备可回溯性。

**Non-Goals:**

- ❌ **不**写任何章节的实际 `.go` 教学代码（这是后续每个 `chapter-NN-*` change 的职责）。
- ❌ **不**实现任何 capstone 项目（这是后续 `capstone-N-*` change 的职责）。
- ❌ **不**引入任何第三方依赖（`go.mod` 仅声明语言版本）。
- ❌ **不**配置 CI/CD、Docker、k8s 等任何运行时基础设施（按需在对应章节引入）。
- ❌ **不**做术语表 / FAQ / references 的内容填充——只建骨架，内容由后续章节落地时增量追加。
- ❌ **不**预判后续每个章节的实现细节（仅给出 ROADMAP 中的"卡片元信息"）。

## Decisions

### D1：仓库形态 = 章节教科书式（A 方案）

**选择**：按主题分目录，每章独立可运行的 demo。
**替代方案**：
- B. 渐进单项目式（一条主线不断重构升级）——叙事性强但跨知识点跳跃不便。
- C. 多迷你项目式——知识点独立但缺乏路线感。
- D. 混合式——最全面但最复杂，对 L0 学习者认知负担过重。

**理由**：学习者明确选择 A，且 L0 阶段需要"知识点对应清晰"的可检索结构，章节式最优。

### D2：路线图深度 = 三段式 15 章 + 3 capstone（P2 方案）

**选择**：阶段一 7 章 + capstone-1；阶段二 4 章 + capstone-2；阶段三 4 章 + capstone-3。
**替代方案**：
- P1（10 章无 capstone）：路线最清爽但 09→10 跨度过大。
- P3（22 章多 capstone）：最完整但马拉松感强，半途放弃风险高。

**理由**：15 章既保留心理可承受目标感，又能在三段末端用 capstone 锁定阶段成果。每章对接一个独立 change，颗粒度恰好（一个 change 一周内可完成）。

### D3：单一 Go module（`module just-go`），统一 `go 1.24`

**选择**：整个仓库使用一个 module，所有章节代码共用 `go.mod`。
**替代方案**：每章独立 module（per-chapter `go.mod`）。

**理由**：
- L0 阶段引入 multi-module / `go.work` 心智成本过高。
- 单 module 下章节间互相 import 简单（用 `just-go/stage-1-syntax/03-composite-types` 这种路径即可）。
- 待 capstone 项目需要独立依赖管理时再用 `go.work` 局部升级，不在本提案预设。

### D4：章节目录仅含 README + .gitkeep（不放任何 .go）

**选择**：本次提案对每个章节目录只生成 `README.md`（章节模板，含"📦 本章产出（待 OpenSpec change 填充）"占位）和 `.gitkeep`。
**替代方案**：放一个 `package xxx` 空文件占位。

**理由**：
- 空 `.go` 文件会让 `go build ./...` 在没有任何代码的情况下成功，但会让后续章节落地时不得不修改既有文件，违反"每章新建"的清爽语义。
- `.gitkeep` 是社区通行做法，零理解成本。
- 由于完全不放 `.go`，`go build ./...` 直接成功（没有任何包要编译）。

### D5：三段在文件系统层面物理分目录

**选择**：`stage-1-syntax/` / `stage-2-business/` / `stage-3-architecture/` 三个顶层目录承载所有章节。
**替代方案**：所有章节平铺在根目录（`01-hello-go/` ... `15-resilience-perf/`）。

**理由**：
- 物理分目录让"三段式叙事"在 `tree` 输出中一眼可见，强化"小白→大神→架构师"的心理路标。
- 平铺方式 15 个章节目录散在根目录，与 `docs/` `openspec/` 等管理目录混在一起，视觉噪声大。

### D6：跟学工作流 = 一章一 OpenSpec change

**选择**：未来每完成一章学习内容，对应一个独立的 OpenSpec change：
- 章节：`chapter-NN-<kebab-name>`（如 `chapter-01-hello-go`）
- Capstone：`capstone-N-<kebab-name>`（如 `capstone-1-cli-todo`）
- 跨章修订：`revise-<topic>`（如 `revise-roadmap-add-security`）

**替代方案**：以"阶段"为粒度，每阶段一个大 change。

**理由**：
- 章节粒度恰好对应"一周可完成"的工作量，与 OpenSpec 工作流的最佳实践吻合。
- 阶段粒度的 change 容易膨胀，违反 OpenSpec"一提案一职责"的原则。
- 章节粒度便于回顾 / 复盘 / 单独重做（学完发现某章学得不扎实可单独重提 revise change）。

### D7：文档语言 = 中文为主；代码标识符 = 英文

**选择**：所有 README / ROADMAP / glossary / faq / references 用中文；未来章节代码的注释用中文，但变量 / 函数 / 包名遵守 Go 社区惯例用英文。
**替代方案**：中英双语 / 全英文。

**理由**：
- 学习者母语中文，技术文档全英文会增加无谓的认知阻力。
- 代码标识符英文是 Go 社区强约束（`gofmt` / `golint` / 第三方库），不可妥协。
- 双语文档会让文档维护成本翻倍，性价比低。

### D8：License = MIT

**选择**：MIT 协议。
**替代方案**：无 License / Apache 2.0 / 私有。

**理由**：MIT 是学习仓库的默认选择，允许任何复用且对学习者本人毫无负担。

### D9：本次提案抽象出 1 个 capability（`learning-curriculum`）

**选择**：把"课程结构 + 章节脚手架 + 跟学工作流"合并为一个 capability。
**替代方案**：拆为 `learning-roadmap` / `chapter-scaffold` / `learning-workflow` 三个 capability。

**理由**：
- 这三块**强耦合**：动一个就要联动改其它两个（如果重构章节模板，roadmap 中的卡片字段也要变）。
- OpenSpec 的 capability 颗粒度建议是"能独立演进"的功能区，强耦合的三个不应拆分。
- 若未来某一块演化出独立性，可在那时拆分（用 `rename` 类型 delta）。

## Risks / Trade-offs

| 风险 | 缓解 |
|---|---|
| **后续 15+3 个 change 偏离本路线图** | 在每个章节 change 的 `proposal.md` 中显式引用 `ROADMAP.md` 对应卡片作为合规基线；spec 自审阶段校对一致性。 |
| **路线图本身需要修订（如 Go 新版本带来新主题）** | 通过 `revise-roadmap-*` 命名的独立 change 修订，禁止直接编辑 `ROADMAP.md`。变更可追溯。 |
| **章节难度估算失真（预计 3h，实际 1 天）** | 章节落地后由 chapter change 用 `MODIFIED Requirements` 回填实测耗时；预计耗时本就是"指南"非"承诺"。 |
| **章节 README 模板流于形式（学习者不填自测清单）** | 模板中"✅ 自测清单"作为章节验收强制项，章节 change 的 tasks.md 必须包含"补全自测清单"任务。 |
| **学习中断后失去上下文** | 三重保障：① `ROADMAP.md` 的进度追踪表；② Git 历史（每章 change 一个 PR/commit 序列）；③ Capstone 项目作为阶段锚点。 |
| **同时打开多份章节 change 互相阻塞** | 显式约定一次只活跃一个 chapter change；OpenSpec 状态机本身也提示当前 in-progress change。 |
| **`.gitkeep` 在 Windows / Linux 大小写不一致** | 全部使用小写 `.gitkeep`，与 GitHub 等平台习惯一致。 |
| **Go 1.24 在学习者机器上未安装** | 在 `01-hello-go` 章节的 README 中明确指引（章节落地时处理）；本提案只在 `go.mod` 中声明，不强行校验。 |

## Migration Plan

本提案**无需迁移**——仓库从无到有的奠基。

回退策略：如本提案产生的所有文件需要撤回，直接 `git reset --hard` 到 `git init` 之前的状态即可（但因为本提案就是首次 commit 的内容，回退即等同于回到"白纸"状态）。

## Open Questions

无。所有关键决策已在 brainstorming 阶段与学习者对齐确认。
