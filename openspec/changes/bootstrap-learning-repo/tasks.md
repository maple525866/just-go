## 1. 基础设施与配置

- [x] 1.1 在仓库根目录运行 `git init`，并将默认分支设置为 `main`
- [x] 1.2 更新 `go.mod`，在 `module just-go` 之后追加 `go 1.24` 语句
- [x] 1.3 创建 `.gitignore`，覆盖 Go 标准忽略项（bin/、*.exe、*.test、*.out、vendor/、覆盖率报告）以及常见 IDE / 系统文件（.idea/、.vscode/、.DS_Store）
- [x] 1.4 创建 `.editorconfig`，统一缩进（Go 用 tab，其它用 2 空格）、换行符（LF）、字符集（utf-8）、末尾空格策略
- [x] 1.5 创建 `LICENSE` 文件，使用 MIT 协议模板

## 2. 顶层文档

- [x] 2.1 创建 `README.md`：含定位、适用人群（L0 / L2 / L3 三类）、学习方式、快速开始、目录速览、指向 ROADMAP 的链接；**不**复述 15 章详表
- [x] 2.2 创建 `ROADMAP.md`：按设计 4 段中确认的六段结构（① 头部 ② 三段总览 ③ 15 章详表 ④ 3 个 capstone ⑤ 跟学方式 ⑥ 进度追踪表）逐段填充，每章按统一卡片模板（🎯 / 🧩 / 📦 / 🔗 / ⏱️ / 📚 六字段）
- [x] 2.3 在 `ROADMAP.md` 中将"📦 章节产出"对全部 15 章 + 3 capstone 统一写为"待 OpenSpec change 填充"占位
- [x] 2.4 在 `ROADMAP.md` 末尾的"进度追踪"段中，按阶段一/二/三分组列出 15 章 + 3 capstone 的全部 `- [ ]` 项

## 3. 阶段一目录骨架（7 章 + 1 capstone）

- [x] 3.1 创建 `stage-1-syntax/01-hello-go/README.md`（章节模板，预填学习目标与关键知识点）与 `.gitkeep`
- [x] 3.2 创建 `stage-1-syntax/02-language-basics/README.md` 与 `.gitkeep`
- [x] 3.3 创建 `stage-1-syntax/03-composite-types/README.md` 与 `.gitkeep`
- [x] 3.4 创建 `stage-1-syntax/04-interface-error/README.md` 与 `.gitkeep`
- [x] 3.5 创建 `stage-1-syntax/05-concurrency/README.md` 与 `.gitkeep`
- [x] 3.6 创建 `stage-1-syntax/06-stdlib-essentials/README.md` 与 `.gitkeep`
- [x] 3.7 创建 `stage-1-syntax/07-engineering/README.md` 与 `.gitkeep`
- [x] 3.8 创建 `stage-1-syntax/capstone-1-cli-todo/README.md`（capstone 模板）与 `.gitkeep`

## 4. 阶段二目录骨架（4 章 + 1 capstone）

- [x] 4.1 创建 `stage-2-business/08-web-foundations/README.md` 与 `.gitkeep`
- [x] 4.2 创建 `stage-2-business/09-data-persistence/README.md` 与 `.gitkeep`
- [x] 4.3 创建 `stage-2-business/10-cache-and-mq/README.md` 与 `.gitkeep`
- [x] 4.4 创建 `stage-2-business/11-observability/README.md` 与 `.gitkeep`
- [x] 4.5 创建 `stage-2-business/capstone-2-blog-api/README.md`（capstone 模板）与 `.gitkeep`

## 5. 阶段三目录骨架（4 章 + 1 capstone）

- [x] 5.1 创建 `stage-3-architecture/12-clean-architecture/README.md` 与 `.gitkeep`
- [x] 5.2 创建 `stage-3-architecture/13-ddd-patterns/README.md` 与 `.gitkeep`
- [x] 5.3 创建 `stage-3-architecture/14-microservices/README.md` 与 `.gitkeep`
- [x] 5.4 创建 `stage-3-architecture/15-resilience-perf/README.md` 与 `.gitkeep`
- [x] 5.5 创建 `stage-3-architecture/capstone-3-blog-ms/README.md`（capstone 模板）与 `.gitkeep`

## 6. 跨章节支持文档

- [x] 6.1 创建 `docs/glossary.md`：表头 `| 术语 | 英文 | 解释 | 出现章节 |` 的空表 + 顶部一段使用说明
- [x] 6.2 创建 `docs/faq.md`：清单结构 + 一条示范条目（"为什么 `go run` 提示找不到包？"）
- [x] 6.3 创建 `docs/references.md`：预填三段（① 官方资源含 go.dev / Tour of Go / Effective Go；② 中文书单含《Go 程序设计语言》《Go 语言高级编程》等；③ 博客与社区）

## 7. 验证与提交

- [x] 7.1 列举仓库目录树（`tree -L 3` 或 PowerShell 等效命令），与 `specs/learning-curriculum/spec.md` 中要求的目录清单逐项核对
- [x] 7.2 在仓库根目录执行 `go build ./...`，确认退出码为 0（应无包可编译）
- [x] 7.3 抽查任意一个章节 README（建议 `01-hello-go` 与 `15-resilience-perf` 各一）与一个 capstone README，确认模板字段齐全
- [x] 7.4 抽查 `ROADMAP.md` 的进度追踪表，确认 15 章 + 3 capstone 全部以 `- [ ]` 列出
- [x] 7.5 运行 `openspec validate bootstrap-learning-repo --strict`（如可用），确保本 change 的 4 件产出物合规
- [ ] 7.6 将所有产出物加入 Git 暂存区（`git add -A`），并以 `chore: bootstrap learning curriculum and repo skeleton` 作为首次 commit 信息提交
- [ ] 7.7 执行 `git log --oneline` 与 `git status`，确认首次 commit 已就位且工作区干净
