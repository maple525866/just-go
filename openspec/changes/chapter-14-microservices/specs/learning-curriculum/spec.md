## MODIFIED Requirements

### Requirement: 每个章节与 capstone 目录 SHALL 仅包含 README 与占位文件

章节目录与 capstone 目录的内容约束 SHALL 依其落地状态区分：

- **未落地章节**（尚无对应 `chapter-NN-*` / `capstone-N-*` change 实现）MUST 仅包含使用统一模板生成的 `README.md` 与用于 Git 跟踪空目录的 `.gitkeep`，且 MUST NOT 包含任何 `.go` 文件。
- **已落地章节**（对应 OpenSpec change 已落地）MAY 包含示例源码（`.go`）、测试文件（`_test.go`）以及 `EXERCISES.md` 等学习材料；其 `README.md` MUST 不再保留"待 OpenSpec change 填充"占位。

无论落地与否，仓库整体 MUST 始终保持 `go build ./...` 成功。

#### Scenario: 未落地章节目录内容受严格约束

- **WHEN** 阅读者列出一个尚未落地的章节目录（如 `stage-3-architecture/15-resilience-perf/`）
- **THEN** 该目录 MUST 仅包含 `README.md` 与 `.gitkeep` 两个文件，且 MUST NOT 存在任何后缀为 `.go` 的文件

#### Scenario: 已落地章节目录可包含源码与测试

- **WHEN** 阅读者列出一个已落地的章节目录（如 `stage-3-architecture/14-microservices/`）
- **THEN** 该目录 MAY 包含 `.go` 源码、`_test.go` 测试与 `EXERCISES.md`，且其 `README.md` MUST NOT 再包含"待 OpenSpec change 填充"占位语

#### Scenario: 仓库整体可通过 go build

- **WHEN** 在仓库根目录执行 `go build ./...`
- **THEN** 命令 MUST 以退出码 0 成功返回
