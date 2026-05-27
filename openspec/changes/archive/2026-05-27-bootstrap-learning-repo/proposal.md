## Why

仓库当前是一张白纸（无 .go 源码、无 git、无路线规划），但目标是承载一条「Go 零基础 → 语法熟练 → 业务工程 → 业务架构师」的系统化学习路径。在写任何章节代码之前，必须先把"路线图 + 仓库骨架 + 跟学工作流"立起来，否则后续每一章都会失去对齐参照，最终演变成无方向的代码堆。

本次提案就是这块"奠基石"——不写章节内容，只交付一份可被后续 N 个章节提案稳定引用的「课程契约」。

## What Changes

- **新增**学习路线图 `ROADMAP.md`：三阶段 15 章 + 3 个 capstone 项目，每章按统一卡片模板呈现（学习目标 / 关键知识点 / 章节产出 / 前置依赖 / 预计耗时 / 推荐扩展）。
- **新增**仓库门面 `README.md`：定位、适用人群、学习方式、目录速览。
- **新增**仓库物理骨架：
  - `stage-1-syntax/`（7 章 + 1 capstone）
  - `stage-2-business/`（4 章 + 1 capstone）
  - `stage-3-architecture/`（4 章 + 1 capstone）
  - 每章/每 capstone 目录内仅含 `README.md`（统一模板）+ `.gitkeep`，**不含任何 .go 源码**。
- **新增**跨章节文档 `docs/`：`glossary.md`（术语表骨架）/ `faq.md`（含一条示范）/ `references.md`（预填三类参考资料）。
- **新增**工程基础设施：`LICENSE`（MIT）、`.gitignore`（Go 标准）、`.editorconfig`。
- **更新** `go.mod`：声明 `go 1.24`。
- **初始化** Git 仓库：`git init`，默认分支 `main`，首次 commit。
- **建立**跟学工作流约定：每一章 = 一个独立 OpenSpec change，命名规范 `chapter-NN-<kebab-name>` 与 `capstone-N-<kebab-name>`。

## Capabilities

### New Capabilities

- `learning-curriculum`：定义本仓库作为"Go 学习教科书"的完整契约，包括三阶段划分、15 章 + 3 capstone 的学习路线、每个章节目录的统一结构（README 模板 + .gitkeep）、跨章节文档的组织方式，以及每章一个 OpenSpec change 的跟学工作流。后续所有 `chapter-NN-*` 和 `capstone-N-*` 提案都将以本 spec 为合规基线。

### Modified Capabilities

（无——仓库为全新状态，没有既有 capability 需要修改。）

## Impact

- **代码**：本次零业务代码改动，仅增加文档与目录占位。`go build ./...` 与 `go vet ./...` 必须保持通过（因不存在 .go 文件，应直接成功返回）。
- **目录结构**：仓库根目录从 1 个 `go.mod` 扩张为约 48 个新增文件（详见 `tasks.md`）。
- **版本控制**：仓库从「非 Git 项目」转变为 Git 项目，并产生第一条提交记录。
- **依赖**：不引入任何 Go 第三方依赖。`go.mod` 仅声明语言版本。
- **未来提案的契约**：本提案确立的章节命名、目录结构、README 模板成为约束，后续修改需走 `revise-roadmap-*` 或 `revise-chapter-template-*` 这类显式修订提案。
- **学习者体验**：从"不知道从哪开始"变为"打开 `ROADMAP.md` 即可看到完整路径，按编号顺序运行 `/opsx-propose chapter-NN-xxx` 推进"。
